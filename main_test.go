package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name    string
		proj    string
		out     string
		expErr  error
		mockCmd func(ctx context.Context, name string, args ...string) *exec.Cmd
	}{
		{
			name: "success", proj: "./testdata/tool/",
			out:     "Go Build: SUCCESS\n" + "Go Test: SUCCESS\n" + "Gofmt: SUCCESS\n" + "Git Push: SUCCESS\n",
			expErr:  nil,
			mockCmd: nil,
		},
		{
			name: "successMock", proj: "./testdata/tool/",
			out:     "Go Build: SUCCESS\n" + "Go Test: SUCCESS\n" + "Gofmt: SUCCESS\n" + "Git Push: SUCCESS\n",
			expErr:  nil,
			mockCmd: mockCmdContext,
		},
		{
			name: "fail", proj: "./testdata/toolErr",
			out:     "",
			expErr:  &stepErr{step: "go build"},
			mockCmd: nil,
		},
		{
			name: "failFormat", proj: "./testdata/toolFmtErr",
			out:     "",
			expErr:  &stepErr{step: "go fmt"},
			mockCmd: nil,
		},
		{
			name: "failTimeout", proj: "./testdata/tool/",
			out:     "",
			expErr:  context.DeadlineExceeded,
			mockCmd: mockCmdTimeout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			var out bytes.Buffer
			err := run(tc.proj, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got 'nil' instead.", tc.expErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tc.out {
				t.Errorf("Expected output: %q. Got %q instead.", tc.out, out.String())
			}
		})
	}
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)

	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}

	os.Exit(1)
}
