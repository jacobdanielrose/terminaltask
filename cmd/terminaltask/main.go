package main

import (
	"flag"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jacobdanielrose/terminaltask/internal/app"
	"github.com/jacobdanielrose/terminaltask/internal/config"
	taskservice "github.com/jacobdanielrose/terminaltask/internal/service"
	"github.com/jacobdanielrose/terminaltask/internal/store"
)

var (
	version   = "dev"
	commit    = "dev"
	buildDate = "dev"
)

func run() error {
	ver := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("terminaltask v%s (commit=%s, built=%s)\n",
			version, commit, buildDate)
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	taskStore := store.NewFileTaskStore(cfg.TasksFile)
	taskService := taskservice.NewFileTaskService(taskStore)
	model := app.NewModel(cfg, taskService)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run program: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal("terminaltask exited with error", "err", err)
	}
}
