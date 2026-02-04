package main

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	version   = "dev"
	commit    = "dev"
	buildDate = "dev"
)

func run() error {
	app := NewApp(AppEnv{})
	return app.Run(os.Args[1:])
}

func main() {
	if err := run(); err != nil {
		log.Fatal("terminaltask exited with error", "err", err)
	}
}
