package config

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to restore environment variables after a test
func withEnv(key, value string, fn func()) {
	orig, had := os.LookupEnv(key)
	if value == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, value)
	}
	defer func() {
		if had {
			_ = os.Setenv(key, orig)
		} else {
			_ = os.Unsetenv(key)
		}
	}()
	fn()
}

func TestLoad_UsesEnvConfigDir(t *testing.T) {
	// Use a unique temp directory under the system temp dir.
	tmpBase := t.TempDir()
	customDir := filepath.Join(tmpBase, "my-terminaltask-config")

	withEnv("TERMINALTASK_CONFIG_DIR", customDir, func() {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		if cfg.ConfigDir != customDir {
			t.Fatalf("ConfigDir = %q, want %q", cfg.ConfigDir, customDir)
		}

		// TasksFile should be ConfigDir/tasks.json
		wantTasksFile := filepath.Join(customDir, "tasks.json")
		if cfg.TasksFile != wantTasksFile {
			t.Fatalf("TasksFile = %q, want %q", cfg.TasksFile, wantTasksFile)
		}

		// The directory should have been created.
		info, err := os.Stat(customDir)
		if err != nil {
			t.Fatalf("expected config dir %q to exist, stat error: %v", customDir, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %q to be a directory", customDir)
		}
	})
}

func TestLoad_UsesUserConfigDirWhenEnvUnset(t *testing.T) {
	// Ensure the env var is unset for this test.
	withEnv("TERMINALTASK_CONFIG_DIR", "", func() {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		// When env is unset, ConfigDir should live under os.UserConfigDir().
		userCfgDir, err := os.UserConfigDir()
		if err != nil {
			t.Fatalf("os.UserConfigDir() returned error: %v", err)
		}

		expectedPrefix := filepath.Join(userCfgDir, "terminaltask")
		if cfg.ConfigDir != expectedPrefix {
			t.Fatalf("ConfigDir = %q, want %q", cfg.ConfigDir, expectedPrefix)
		}

		// TasksFile should be ConfigDir/tasks.json
		wantTasksFile := filepath.Join(cfg.ConfigDir, "tasks.json")
		if cfg.TasksFile != wantTasksFile {
			t.Fatalf("TasksFile = %q, want %q", cfg.TasksFile, wantTasksFile)
		}

		// The directory should have been created.
		info, err := os.Stat(cfg.ConfigDir)
		if err != nil {
			t.Fatalf("expected config dir %q to exist, stat error: %v", cfg.ConfigDir, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %q to be a directory", cfg.ConfigDir)
		}
	})
}
