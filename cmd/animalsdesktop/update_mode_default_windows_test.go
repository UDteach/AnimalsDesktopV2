//go:build windows && !animalsdesktop_nonetwork

package main

import "testing"

func TestDefaultBuildEnablesNetworkUpdates(t *testing.T) {
	if !networkUpdatesEnabled {
		t.Fatalf("default Windows build should keep update network access enabled")
	}
}

func TestVer2UpdatesUseSeparateRepository(t *testing.T) {
	if updateAPIURL != "https://api.github.com/repos/UDteach/AnimalsDesktopV2/releases/latest" {
		t.Fatalf("Ver2 must fetch its own releases, got %q", updateAPIURL)
	}
	if !isNewerVersion("v2.0.0", "v0.2.16") {
		t.Fatal("Ver2 must be newer than the last V1 release")
	}
}
