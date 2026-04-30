// Copyright (c) 2023-2026, Nubificus LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unikontainers

import (
	"strings"
	"syscall"
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestSignalRejectsInvalidPIDForNonKill(t *testing.T) {
	t.Parallel()

	u := &Unikontainer{
		State: &specs.State{
			Pid: 0,
		},
	}

	err := u.Signal(syscall.SIGTERM)
	if err == nil {
		t.Fatalf("expected error for invalid pid, got nil")
	}

	if !strings.Contains(err.Error(), "invalid monitor pid") {
		t.Fatalf("expected invalid monitor pid error, got: %v", err)
	}
}

func TestSignalForwardsNonKillSignals(t *testing.T) {
	t.Parallel()

	u := &Unikontainer{
		State: &specs.State{
			Pid: syscall.Getpid(),
		},
	}

	// Signal 0 checks process existence without delivering a signal.
	if err := u.Signal(syscall.Signal(0)); err != nil {
		t.Fatalf("expected non-kill signal forwarding to succeed, got: %v", err)
	}
}
