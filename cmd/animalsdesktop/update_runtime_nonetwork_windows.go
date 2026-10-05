//go:build windows && animalsdesktop_nonetwork

package main

import "fmt"

func downloadAndStartUpdater(githubReleaseAsset) error {
	return fmt.Errorf("network updates are disabled in this build")
}

func runUpdaterUtility([]string) bool {
	return false
}

func updateCleanupDir([]string) string {
	return ""
}

func cleanupUpdateTempDirLater(string) {}
