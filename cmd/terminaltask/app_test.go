package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/config"
)

// --- Test helpers ---

type bufferPrinter struct {
	buf *bytes.Buffer
}

func (p bufferPrinter) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(p.buf, format, a...)
}

type fakeProgramRunner struct {
	runs int
	err  error
	last tea.Model
}

func (f *fakeProgramRunner) Run(m tea.Model) error {
	f.runs++
	f.last = m
	return f.err
}

// --- Tests ---

func TestVersionFlagPrintsAndExits(t *testing.T) {
	var out bytes.Buffer
	printer := bufferPrinter{buf: &out}

	app := NewApp(AppEnv{
		Printer: printer,
		// LoadConfig and ProgramRunner can be left nil; they shouldn't be used on --version
	})

	err := app.Run([]string{"-version"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "terminaltask v") {
		t.Fatalf("expected version output, got %q", got)
	}
}

func TestConfigLoadErrorIsWrapped(t *testing.T) {
	loadErr := errors.New("boom")

	a := NewApp(AppEnv{
		LoadConfig: func() (config.Config, error) {
			return config.Config{}, loadErr
		},
		ProgramRunner: &fakeProgramRunner{}, // won't be called
	})

	err := a.Run([]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, loadErr) {
		t.Fatalf("expected to wrap load error, got %v", err)
	}
}

func TestProgramRunnerIsInvoked(t *testing.T) {
	fakeRunner := &fakeProgramRunner{}

	a := NewApp(AppEnv{
		LoadConfig: func() (config.Config, error) {
			return config.Config{
				TasksFile: "/tmp/tasks.json",
			}, nil
		},
		ProgramRunner: fakeRunner,
	})

	err := a.Run([]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fakeRunner.runs != 1 {
		t.Fatalf("expected program runner to run once, ran %d times", fakeRunner.runs)
	}
}
