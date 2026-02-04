package main

import (
	"flag"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/app"
	"github.com/jacobdanielrose/terminaltask/internal/config"
	taskservice "github.com/jacobdanielrose/terminaltask/internal/service"
	"github.com/jacobdanielrose/terminaltask/internal/store"
)

type CLIOptions struct {
	ShowVersion bool
}

func parseArgs(args []string) (CLIOptions, error) {
	fs := flag.NewFlagSet("terminaltask", flag.ContinueOnError)
	fs.SetOutput(nil)

	var opts CLIOptions
	fs.BoolVar(&opts.ShowVersion, "version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		return CLIOptions{}, err
	}

	return opts, nil
}

type Printer interface {
	Printf(format string, a ...any)
}

type ConfigLoader func() (config.Config, error)

type ProgramRunner interface {
	Run(model tea.Model) error
}

// StdoutPrinter uses fmt.Printf; in tests you can provide a different Printer.
type StdoutPrinter struct{}

func (StdoutPrinter) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

// TeaProgramRunner runs the real Bubble Tea program.
type TeaProgramRunner struct{}

func (TeaProgramRunner) Run(model tea.Model) error {
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

type AppEnv struct {
	Printer       Printer
	LoadConfig    ConfigLoader
	ProgramRunner ProgramRunner
}

type App struct {
	env AppEnv
}

func NewApp(env AppEnv) *App {
	// Provide reasonable defaults for production when fields are nil
	if env.Printer == nil {
		env.Printer = StdoutPrinter{}
	}
	if env.LoadConfig == nil {
		env.LoadConfig = config.Load
	}
	if env.ProgramRunner == nil {
		env.ProgramRunner = TeaProgramRunner{}
	}
	return &App{env: env}
}

func (a *App) Run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	if opts.ShowVersion {
		a.env.Printer.Printf(
			"terminaltask v%s (commit=%s, built=%s)\n",
			version, commit, buildDate,
		)
		return nil
	}

	cfg, err := a.env.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	taskStore := store.NewFileTaskStore(cfg.TasksFile)
	taskService := taskservice.NewFileTaskService(taskStore)
	model := app.NewModel(cfg, taskService)

	if err := a.env.ProgramRunner.Run(model); err != nil {
		return fmt.Errorf("run program: %w", err)
	}

	return nil
}
