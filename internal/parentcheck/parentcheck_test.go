// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package parentcheck

import (
	"os"
	"testing"
)

func TestCheckParentProcess_SkipEnv(t *testing.T) {
	t.Setenv(skipEnv, "1")
	if err := CheckParentProcess(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckParentProcess_WithoutSkip(t *testing.T) {
	_ = os.Unsetenv(skipEnv)
	// Under `go test`, the parent is the test binary, not arksync_client — expect an error.
	if err := CheckParentProcess(); err == nil {
		t.Fatal("expected error when parent is not arksync_client")
	}
}
