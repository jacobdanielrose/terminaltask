package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/storage"
)

var version = "dev"

func main() {

	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		fmt.Println("terminaltask v0.1.0-alpha")
		return
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
