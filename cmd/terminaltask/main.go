package main

import (
	"flag"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jacobdanielrose/terminaltask/internal/app"
	"github.com/jacobdanielrose/terminaltask/internal/config"
	"github.com/jacobdanielrose/terminaltask/internal/store"
)

var (
	version   = "dev"
	commit    = "dev"
	buildDate = "dev"
)

func main() {
	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		fmt.Printf("terminaltask v%s (commit=%s, built=%s)\n",
			version, commit, buildDate)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", "err", err)
	}

	store := store.NewFileTaskStore(cfg.TasksFile)
	model := app.NewModel(cfg, store)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
}
