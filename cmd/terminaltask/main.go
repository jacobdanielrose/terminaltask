package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/app"
	"github.com/jacobdanielrose/terminaltask/internal/storage"
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

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(cfgDir, "terminaltask", "tasks.json")
	store := storage.NewFileTaskStore(path)

	m := app.NewModel(store)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
