//go:build windows && animalsdesktop_nonetwork

package main

import (
	"strings"
	"testing"
)

func TestNoNetworkBuildDisablesUpdateMenus(t *testing.T) {
	if networkUpdatesEnabled {
		t.Fatalf("no-network Windows build should disable update network access")
	}
	a := &petApp{lang: langEnglish}
	if got := a.updateCheckMenuLabel(); !strings.Contains(got, "Network disabled") {
		t.Fatalf("updateCheckMenuLabel() = %q, want network disabled label", got)
	}
	if a.hasInstallableUpdate() {
		t.Fatalf("no-network build should not report installable updates")
	}
}
