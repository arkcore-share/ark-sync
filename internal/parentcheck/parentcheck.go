// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// Package parentcheck exits the process unless the immediate parent executable
// is the expected ArkSync client (arksync_client.exe on Windows, arksync_client elsewhere),
// or the same resolved path as this binary (monitor respawning the inner process).
//
// Set ARKSYNC_SKIP_PARENT_CHECK=1 to disable (tests, CI, or service wrappers).
// Set ARKSYNC_DEBUG_PARENT=1 to print parent PID, inferred parent image base name,
// and self path to stderr (useful on Windows when the console is not obvious).
package parentcheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	psprocess "github.com/shirou/gopsutil/v4/process"

	"github.com/syncthing/syncthing/lib/build"
)

const skipEnv = "ARKSYNC_SKIP_PARENT_CHECK"

// ErrWrongParent indicates the parent process did not match the required client
// (or the same executable for the monitor/inner-process chain).
var ErrWrongParent = errors.New("parent process must be arksync_client (arksync_client.exe on Windows), or the same arksync executable respawned by the monitor")

// CheckParentProcess returns nil if the parent is the approved launcher, or if
// checks are skipped via ARKSYNC_SKIP_PARENT_CHECK=1.
func CheckParentProcess() error {
	if os.Getenv(skipEnv) == "1" {
		return nil
	}

	self, err := psprocess.NewProcess(int32(os.Getpid()))
	if err != nil {
		return fmt.Errorf("parent check: %w", err)
	}
	ppid, err := self.Ppid()
	if err != nil {
		return fmt.Errorf("parent check: %w", err)
	}
	if ppid <= 1 {
		return fmt.Errorf("%w (no valid parent)", ErrWrongParent)
	}

	parent, err := psprocess.NewProcess(ppid)
	if err != nil {
		return fmt.Errorf("parent check: %w", err)
	}

	base, err := parentExecutableBase(parent)
	if err != nil {
		return fmt.Errorf("parent check: %w", err)
	}

	if isApprovedParentBase(base) {
		debugParent(base, "allowed: parent matches arksync_client")
		return nil
	}
	// Monitor respawns the same binary for the inner process only (child has
	// STMONITORED set). Without this guard, any same-path arksync parent could
	// spawn a second unrestricted main.
	if os.Getenv("STMONITORED") != "" {
		if same, err := sameResolvedExecutable(parent); err == nil && same {
			debugParent(base, "allowed: same executable as parent (STMONITORED inner process)")
			return nil
		}
	}
	debugParent(base, "rejected")
	return fmt.Errorf("%w (got %q)", ErrWrongParent, base)
}

// ExitIfWrongParent prints a message to stderr and exits with code 1 if
// CheckParentProcess returns an error.
func ExitIfWrongParent() {
	if err := CheckParentProcess(); err != nil {
		fmt.Fprintf(os.Stderr, "arksync: %v\n", err)
		fmt.Fprintf(os.Stderr, "arksync: set %s=1 only for tests or controlled service environments.\n", skipEnv)
		os.Exit(1)
	}
}

func parentExecutableBase(parent *psprocess.Process) (string, error) {
	exe, err := parent.Exe()
	if err == nil && exe != "" {
		return filepath.Base(exe), nil
	}
	name, err := parent.Name()
	if err != nil {
		return "", err
	}
	return filepath.Base(name), nil
}

func isApprovedParentBase(base string) bool {
	b := strings.ToLower(filepath.Base(strings.TrimSpace(base)))
	if build.IsWindows {
		// Some APIs report the image name without ".exe".
		return b == "arksync_client.exe" || b == "arksync_client"
	}
	return b == "arksync_client"
}

func debugParent(parentBase, outcome string) {
	if os.Getenv("ARKSYNC_DEBUG_PARENT") != "1" {
		return
	}
	selfExe, _ := os.Executable()
	ppid := int32(-1)
	if self, err := psprocess.NewProcess(int32(os.Getpid())); err == nil {
		if p, err := self.Ppid(); err == nil {
			ppid = p
		}
	}
	fmt.Fprintf(os.Stderr, "arksync: parent-check debug: ppid=%d parentBase=%q selfExe=%q STMONITORED=%q outcome=%s\n",
		ppid, parentBase, selfExe, os.Getenv("STMONITORED"), outcome)
}

func sameResolvedExecutable(parent *psprocess.Process) (bool, error) {
	selfRaw, err := os.Executable()
	if err != nil {
		return false, err
	}
	pRaw, err := parent.Exe()
	if err != nil || pRaw == "" {
		return false, err
	}
	selfExe := selfRaw
	if r, err := filepath.EvalSymlinks(selfRaw); err == nil {
		selfExe = r
	}
	pExe := pRaw
	if r, err := filepath.EvalSymlinks(pRaw); err == nil {
		pExe = r
	}
	if build.IsWindows {
		return strings.EqualFold(filepath.Clean(selfExe), filepath.Clean(pExe)), nil
	}
	return filepath.Clean(selfExe) == filepath.Clean(pExe), nil
}
