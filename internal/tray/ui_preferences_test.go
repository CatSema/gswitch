package tray

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUIPreferencesMigration(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	old := filepath.Join(dir, "gswitch", "tray.conf")
	if err := os.MkdirAll(filepath.Dir(old), 0o700); err != nil {
		t.Fatal(err)
	}
	original := "# preserved\ntray-icon-mode=app-with-flag\nfuture-setting=yes\n"
	if err := os.WriteFile(old, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	prefs, err := LoadUIPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.IconMode != TrayIconModeAppWithFlag || prefs.Language != "auto" {
		t.Fatalf("migration: %+v", prefs)
	}
	path := filepath.Join(dir, "gswitch", "ui.conf")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("migration did not create ui.conf:", err)
	}
	if _, err := os.Lstat(old); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("legacy file retained after migration: %v", err)
	}
	prefs.Language = "ru"
	if err := SaveUIPreferences(prefs); err != nil {
		t.Fatal(err)
	}
	if err := SaveTrayIconMode(TrayIconModeFlag); err != nil {
		t.Fatal(err)
	}
	prefs, err = LoadUIPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.Language != "ru" || prefs.IconMode != TrayIconModeFlag {
		t.Fatalf("lost preference: %+v", prefs)
	}
	// #nosec G304 -- path is inside t.TempDir.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "future-setting=yes") {
		t.Fatal("lost unknown setting")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatal("unsafe permissions")
	}
	// A legacy file must never overwrite an existing ui.conf.
	if err := os.WriteFile(old, []byte("tray-icon-mode=app\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	prefs, err = LoadUIPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.Language != "ru" || prefs.IconMode != TrayIconModeFlag {
		t.Fatal("legacy file overrode ui.conf")
	}
	if _, err := os.Lstat(old); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("legacy file retained beside existing ui.conf: %v", err)
	}
}

func TestUIPreferencesReadFailureDoesNotOverwrite(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, err := defaultTrayPreferencesPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadUIPreferences(); err == nil {
		t.Fatal("directory accepted as preferences")
	}
	if err := SaveUIPreferences(UIPreferences{Language: "ru", IconMode: TrayIconModeFlag}); err == nil {
		t.Fatal("read failure overwritten")
	}
}

func TestUIPreferencesMigrationFailureKeepsLegacy(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, err := defaultTrayPreferencesPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(filepath.Dir(path), "tray.conf")
	original := "tray-icon-mode=app-with-flag\nfuture-setting=yes\n"
	if err := os.WriteFile(legacy, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	// A dangling destination cannot be read or exclusively replaced.
	if err := os.Symlink(filepath.Join(t.TempDir(), "missing"), path); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadUIPreferences(); err == nil {
		t.Fatal("migration failure not reported")
	}
	// #nosec G304 -- temporary test configuration.
	data, err := os.ReadFile(legacy)
	if err != nil || string(data) != original {
		t.Fatalf("legacy changed on failure: %q, %v", data, err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	prefs, err := LoadUIPreferences()
	if err != nil || prefs.IconMode != TrayIconModeAppWithFlag {
		t.Fatalf("retry: %+v, %v", prefs, err)
	}
	if _, err := os.Lstat(legacy); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("retry retained legacy: %v", err)
	}
}

func TestUIPreferencesCleanupFailureKeepsCurrentSettings(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, err := defaultTrayPreferencesPath()
	if err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(filepath.Dir(path), "tray.conf")
	// Never remove a directory in place of the obsolete config file.
	if err := os.MkdirAll(legacy, 0o700); err != nil {
		t.Fatal(err)
	}
	original := "tray-icon-mode=app-with-flag\nui-language=ru\n"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	prefs, err := LoadUIPreferences()
	if err == nil {
		t.Fatal("cleanup failure not reported")
	}
	if prefs.IconMode != TrayIconModeAppWithFlag || prefs.Language != "ru" {
		t.Fatalf("cleanup lost settings: %+v", prefs)
	}
	// #nosec G304 -- temporary test configuration.
	data, err := os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatalf("current settings changed: %q, %v", data, err)
	}
	if err := os.Remove(legacy); err != nil {
		t.Fatal("legacy directory was removed:", err)
	}
	if _, err := LoadUIPreferences(); err != nil {
		t.Fatal("retry failed:", err)
	}
}

func TestUIPreferencesCleanupDoesNotBreakDestinationSymlink(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, err := defaultTrayPreferencesPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(filepath.Dir(path), "tray.conf")
	original := "tray-icon-mode=app-with-flag\nui-language=ru\n"
	if err := os.WriteFile(legacy, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("tray.conf", path); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadUIPreferences(); err == nil {
		t.Fatal("dependent destination accepted for cleanup")
	}
	// #nosec G304 -- temporary test configuration.
	data, err := os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatalf("cleanup broke current config: %q, %v", data, err)
	}
}
