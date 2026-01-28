package config

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// Config holds configuration for terminaltask.
type Config struct {
	// ConfigDir is the directory where terminaltask stores its config/data files.
	// Default: $TERMINALTASK_CONFIG_DIR or, if unset, UserConfigDir()/terminaltask.
	ConfigDir string

	// TasksFile is the full path to the tasks JSON file.
	// Default: ConfigDir/tasks.json.
	TasksFile string
}

// Load builds a Config from environment variables and sensible defaults.
func Load() (Config, error) {
	var cfg Config

	// Base config directory
	if envDir := os.Getenv("TERMINALTASK_CONFIG_DIR"); envDir != "" {
		cfg.ConfigDir = envDir
	} else {
		userCfgDir, err := os.UserConfigDir()
		if err != nil {
			log.Error("getting user config dir", "err", err)
			return Config{}, err
		}
		cfg.ConfigDir = filepath.Join(userCfgDir, "terminaltask")
	}

	// Ensure directory exists
	if err := os.MkdirAll(cfg.ConfigDir, 0o755); err != nil {
		log.Error("creating config dir", "dir", cfg.ConfigDir, "err", err)
		return Config{}, err
	}

	// Tasks file path (can be overridden later with another env var if desired)
	cfg.TasksFile = filepath.Join(cfg.ConfigDir, "tasks.json")

	return cfg, nil
}
