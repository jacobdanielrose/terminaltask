package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jacobdanielrose/terminaltask/internal/storage"
)

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
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
