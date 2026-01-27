package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/storage"
)

var (
	version   = "dev" // ldflags overwrites this
	buildTime = "dev" // ldflags overwrites this
	commit    = "dev" // optional
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("TerminalTask %s (built %s) (commit %s)\n", version, buildTime, commit)
		os.Exit(0)
	}

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(cfgDir, "terminaltask", "tasks.json")
	store := storage.NewFileTaskStore(path)

	m := initialModel(store)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
