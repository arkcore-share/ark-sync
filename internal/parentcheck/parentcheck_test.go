// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package parentcheck

import (
	"os"
	"runtime"
	"testing"
)

func runtimeIsWindows() bool { return runtime.GOOS == "windows" }

func TestCheckParentProcess_SkipEnv(t *testing.T) {
	t.Setenv(skipEnv, "1")
	if err := CheckParentProcess(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckParentProcess_WithoutSkip(t *testing.T) {
	_ = os.Unsetenv(skipEnv)
	// Under `go test`, the parent is the test binary, not arksync_client/electron — expect an error.
	if err := CheckParentProcess(); err == nil {
		t.Fatal("expected error when parent is not an approved launcher")
	}
}

func TestIsApprovedParentBase(t *testing.T) {
	approvedAll := []string{"arksync_client", "electron"}
	approvedWindows := []string{"arksync_client.exe", "electron.exe", "ArkSync_Client.EXE", "Electron.Exe"}
	rejected := []string{"", "explorer.exe", "cmd.exe", "powershell.exe", "arksync.exe", "electronwrapper.exe"}

	for _, name := range approvedAll {
		if !isApprovedParentBase(name) {
			t.Errorf("expected %q to be approved", name)
		}
	}
	for _, name := range rejected {
		if isApprovedParentBase(name) {
			t.Errorf("expected %q to be rejected", name)
		}
	}

	// Windows-specific names must be approved on Windows. On other platforms
	// they should be rejected (e.g. ".exe" suffix isn't expected there).
	for _, name := range approvedWindows {
		got := isApprovedParentBase(name)
		want := runtimeIsWindows()
		if got != want {
			t.Errorf("isApprovedParentBase(%q) = %v, want %v (windows=%v)", name, got, want, runtimeIsWindows())
		}
	}
}
