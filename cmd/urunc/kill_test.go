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

package main

import (
	"strconv"
	"syscall"
	"testing"
)

func TestParseSignal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    syscall.Signal
		wantErr bool
	}{
		{
			name:  "defaults to SIGTERM when empty",
			input: "",
			want:  syscall.SIGTERM,
		},
		{
			name:  "accepts numeric signal",
			input: strconv.Itoa(int(syscall.SIGWINCH)),
			want:  syscall.SIGWINCH,
		},
		{
			name:  "accepts short signal name",
			input: "WINCH",
			want:  syscall.SIGWINCH,
		},
		{
			name:  "accepts prefixed signal name",
			input: "SIGWINCH",
			want:  syscall.SIGWINCH,
		},
		{
			name:    "rejects unknown signal",
			input:   "SIG_NOT_A_SIGNAL",
			wantErr: true,
		},
		{
			name:    "rejects non-positive numeric signal",
			input:   "0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseSignal(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseSignal(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseSignal(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parseSignal(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
