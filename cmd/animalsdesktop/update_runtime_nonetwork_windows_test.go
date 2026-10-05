//go:build windows && animalsdesktop_nonetwork

package main

import (
	"fmt"
	"os/exec"
)

type updateApplyOptions struct {
	Source     string
	Target     string
	ParentPID  int
	CleanupDir string
}

func verifyDownloadedAsset(string, githubReleaseAsset) error {
	return fmt.Errorf("network updates are disabled in this build")
}

func extractUpdateExe(string, string) (string, error) {
	return "", fmt.Errorf("network updates are disabled in this build")
}

func parseUpdateApplyArgs([]string) (updateApplyOptions, error) {
	return updateApplyOptions{}, fmt.Errorf("network updates are disabled in this build")
}

func newUpdaterHelperCommand(string, string, string, string, int) *exec.Cmd {
	return nil
}

func newUpdaterCleanupCommand(string, string) *exec.Cmd {
	return nil
}
