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
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var killCommand = &cli.Command{
	Name:  "kill",
	Usage: "kill sends the specified signal (default: SIGTERM) to the container's init process",
	ArgsUsage: `<container-id> [signal]

Where "<container-id>" is the name for the instance of the container and
"[signal]" is the signal to be sent to the init process.

EXAMPLE:
For example, if the container id is "ubuntu01" the following will send a "KILL"
signal to the init process of the "ubuntu01" container:

	# urunc kill ubuntu01 KILL`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "send the specified signal to all processes inside the container",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		logrus.WithField("command", "KILL").WithField("args", os.Args).Debug("urunc INVOKED")
		if err := checkArgs(cmd, 1, minArgs); err != nil {
			return err
		}
		if err := checkArgs(cmd, 2, maxArgs); err != nil {
			return err
		}

		// get Unikontainer data from state.json
		unikontainer, err := getUnikontainer(cmd)
		if err != nil {
			return err
		}
		sig, err := parseSignal(cmd.Args().Get(1))
		if err != nil {
			return err
		}
		return unikontainer.Signal(sig)
	},
}

var signalMap = map[string]syscall.Signal{
	"HUP":    syscall.SIGHUP,
	"INT":    syscall.SIGINT,
	"QUIT":   syscall.SIGQUIT,
	"ILL":    syscall.SIGILL,
	"TRAP":   syscall.SIGTRAP,
	"ABRT":   syscall.SIGABRT,
	"BUS":    syscall.SIGBUS,
	"FPE":    syscall.SIGFPE,
	"KILL":   syscall.SIGKILL,
	"USR1":   syscall.SIGUSR1,
	"SEGV":   syscall.SIGSEGV,
	"USR2":   syscall.SIGUSR2,
	"PIPE":   syscall.SIGPIPE,
	"ALRM":   syscall.SIGALRM,
	"TERM":   syscall.SIGTERM,
	"STKFLT": syscall.SIGSTKFLT,
	"CHLD":   syscall.SIGCHLD,
	"CONT":   syscall.SIGCONT,
	"STOP":   syscall.SIGSTOP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
	"VTALRM": syscall.SIGVTALRM,
	"PROF":   syscall.SIGPROF,
	"WINCH":  syscall.SIGWINCH,
	"IO":     syscall.SIGIO,
	"PWR":    syscall.SIGPWR,
	"SYS":    syscall.SIGSYS,
}

func parseSignal(sigStr string) (syscall.Signal, error) {
	if sigStr == "" {
		return syscall.SIGTERM, nil
	}

	if sigNum, err := strconv.Atoi(sigStr); err == nil {
		if sigNum <= 0 {
			return 0, fmt.Errorf("invalid signal %q", sigStr)
		}
		return syscall.Signal(sigNum), nil
	}

	normalized := strings.ToUpper(sigStr)
	normalized = strings.TrimPrefix(normalized, "SIG")
	if sig, ok := signalMap[normalized]; ok {
		return sig, nil
	}

	return 0, fmt.Errorf("unknown signal %q", sigStr)
}
