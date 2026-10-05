//go:build windows && !animalsdesktop_nonetwork

package main

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func downloadAndStartUpdater(asset githubReleaseAsset) error {
	if asset.BrowserDownloadURL == "" {
		return fmt.Errorf("update asset has no download URL")
	}
	tmpDir, err := os.MkdirTemp("", updateTempPrefix+"*")
	if err != nil {
		return err
	}
	zipPath := filepath.Join(tmpDir, "update.zip")
	if err := downloadFile(asset.BrowserDownloadURL, zipPath); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	if err := verifyDownloadedAsset(zipPath, asset); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	exePath, err := extractUpdateExe(zipPath, tmpDir)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	currentExe, err := os.Executable()
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	helperPath := filepath.Join(tmpDir, "helper", "AnimalsDesktop.exe")
	if err := os.MkdirAll(filepath.Dir(helperPath), 0o755); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	if err := copyFile(currentExe, helperPath); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}
	return startUpdaterHelper(helperPath, tmpDir, exePath, currentExe, os.Getpid())
}

func verifyDownloadedAsset(path string, asset githubReleaseAsset) error {
	if asset.Size > 0 {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.Size() != asset.Size {
			return fmt.Errorf("downloaded update size mismatch: got %d bytes, want %d", info.Size(), asset.Size)
		}
	}
	if asset.Digest == "" {
		return nil
	}
	algorithm, want, ok := strings.Cut(asset.Digest, ":")
	if !ok || !strings.EqualFold(algorithm, "sha256") || len(want) != 64 {
		return fmt.Errorf("unsupported update digest %q", asset.Digest)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	got := fmt.Sprintf("%x", sum[:])
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("downloaded update digest mismatch")
	}
	return nil
}

func extractUpdateExe(zipPath, tmpDir string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if !strings.EqualFold(filepath.Base(file.Name), "AnimalsDesktop.exe") {
			continue
		}
		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()
		exePath := filepath.Join(tmpDir, "payload", "AnimalsDesktop.exe")
		if err := os.MkdirAll(filepath.Dir(exePath), 0o755); err != nil {
			return "", err
		}
		dst, err := os.Create(exePath)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dst, src); err != nil {
			_ = dst.Close()
			return "", err
		}
		if err := dst.Close(); err != nil {
			return "", err
		}
		return exePath, os.Chmod(exePath, 0o755)
	}
	return "", fmt.Errorf("AnimalsDesktop.exe was not found in update zip")
}

func startUpdaterHelper(helperPath, tmpDir, sourceExe, targetExe string, pid int) error {
	return newUpdaterHelperCommand(helperPath, tmpDir, sourceExe, targetExe, pid).Start()
}

func newUpdaterHelperCommand(helperPath, tmpDir, sourceExe, targetExe string, pid int) *exec.Cmd {
	cmd := exec.Command(
		helperPath,
		updaterApplyArg,
		"--source", sourceExe,
		"--target", targetExe,
		"--parent-pid", strconv.Itoa(pid),
		"--cleanup-dir", tmpDir,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

type updateApplyOptions struct {
	Source     string
	Target     string
	ParentPID  int
	CleanupDir string
}

func runUpdaterUtility(args []string) bool {
	if len(args) == 0 || args[0] != updaterApplyArg {
		return false
	}
	opts, err := parseUpdateApplyArgs(args[1:])
	if err != nil {
		os.Exit(2)
	}
	if err := applyUpdate(opts); err != nil {
		os.Exit(1)
	}
	return true
}

func parseUpdateApplyArgs(args []string) (updateApplyOptions, error) {
	var opts updateApplyOptions
	for i := 0; i < len(args); i++ {
		if i+1 >= len(args) {
			return opts, fmt.Errorf("%s is missing a value", args[i])
		}
		value := args[i+1]
		switch args[i] {
		case "--source":
			opts.Source = value
		case "--target":
			opts.Target = value
		case "--parent-pid":
			pid, err := strconv.Atoi(value)
			if err != nil || pid < 0 {
				return opts, fmt.Errorf("invalid parent pid %q", value)
			}
			opts.ParentPID = pid
		case "--cleanup-dir":
			opts.CleanupDir = value
		default:
			return opts, fmt.Errorf("unknown updater argument %q", args[i])
		}
		i++
	}
	if opts.Source == "" || opts.Target == "" || opts.CleanupDir == "" {
		return opts, fmt.Errorf("updater source, target, and cleanup-dir are required")
	}
	if !isUpdateTempDir(opts.CleanupDir) {
		return opts, fmt.Errorf("refusing cleanup outside update temp dir")
	}
	if !isPathInsideDir(opts.Source, opts.CleanupDir) || !strings.EqualFold(filepath.Base(opts.Source), "AnimalsDesktop.exe") {
		return opts, fmt.Errorf("updater source must be AnimalsDesktop.exe inside the update temp dir")
	}
	if !strings.EqualFold(filepath.Base(opts.Target), "AnimalsDesktop.exe") || isPathInsideDir(opts.Target, opts.CleanupDir) {
		return opts, fmt.Errorf("updater target must be an installed AnimalsDesktop.exe outside the update temp dir")
	}
	return opts, nil
}

func applyUpdate(opts updateApplyOptions) error {
	if opts.ParentPID > 0 {
		if err := waitForProcessExit(opts.ParentPID, 120*time.Second); err != nil {
			return err
		}
		time.Sleep(300 * time.Millisecond)
	}
	if err := copyFile(opts.Source, opts.Target); err != nil {
		return err
	}
	return newUpdaterCleanupCommand(opts.Target, opts.CleanupDir).Start()
}

func newUpdaterCleanupCommand(targetExe, cleanupDir string) *exec.Cmd {
	cmd := exec.Command(targetExe, updaterCleanupArg, cleanupDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func waitForProcessExit(pid int, timeout time.Duration) error {
	handle, err := syscall.OpenProcess(syscall.SYNCHRONIZE, false, uint32(pid))
	if err != nil {
		return nil
	}
	defer syscall.CloseHandle(handle)
	waitMS := uint32(timeout / time.Millisecond)
	if timeout <= 0 {
		waitMS = syscall.INFINITE
	}
	result, err := syscall.WaitForSingleObject(handle, waitMS)
	if err != nil {
		return err
	}
	if result == syscall.WAIT_TIMEOUT {
		return fmt.Errorf("timed out waiting for process %d to exit", pid)
	}
	if result == syscall.WAIT_FAILED {
		return fmt.Errorf("failed waiting for process %d", pid)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dst, 0o755)
}

func updateCleanupDir(args []string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == updaterCleanupArg && isUpdateTempDir(args[i+1]) {
			return args[i+1]
		}
	}
	return ""
}

func cleanupUpdateTempDirLater(dir string) {
	for i := 0; i < 20; i++ {
		time.Sleep(500 * time.Millisecond)
		if err := os.RemoveAll(dir); err == nil {
			return
		}
	}
}

func isUpdateTempDir(path string) bool {
	if path == "" {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absTemp, err := filepath.Abs(os.TempDir())
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absTemp, absPath)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return false
	}
	return strings.HasPrefix(filepath.Base(absPath), updateTempPrefix)
}

func isPathInsideDir(path, dir string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absDir, absPath)
	if err != nil || rel == "." || filepath.IsAbs(rel) {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}
