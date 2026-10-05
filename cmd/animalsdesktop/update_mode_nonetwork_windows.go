//go:build windows && animalsdesktop_nonetwork

package main

import "fmt"

const (
	networkUpdatesEnabled = false
	updateAPIURL          = ""
)

func fetchLatestRelease() (*githubRelease, error) {
	return nil, fmt.Errorf("network updates are disabled in this build")
}

func downloadFile(string, string) error {
	return fmt.Errorf("network downloads are disabled in this build")
}
