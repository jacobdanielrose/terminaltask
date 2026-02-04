package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestRun_VersionFlag(t *testing.T) {
	// Save and restore original args and stdout.
	origArgs := os.Args
	origStdout := os.Stdout
	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
	}()

	// Reset the default flag set so we can parse flags cleanly for this test.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Override version variables for predictable output.
	version = "1.2.3"
	commit = "abcd1234"
	buildDate = "2026-02-03"

	os.Args = []string{"terminaltask", "-version"}

	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() failed: %v", err)
	}
	os.Stdout = pw

	if err := run(); err != nil {
		t.Fatalf("run() returned error with -version: %v", err)
	}

	_ = pw.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, pr); err != nil {
		t.Fatalf("io.Copy() failed: %v", err)
	}
	_ = pr.Close()

	got := buf.String()
	want := fmt.Sprintf("terminaltask v%s (commit=%s, built=%s)\n", version, commit, buildDate)
	if got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}
