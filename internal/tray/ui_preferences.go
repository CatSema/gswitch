package tray

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

// UIPreferences belongs to the graphical user, not the daemon configuration.
type UIPreferences struct {
	IconMode TrayIconMode
	Language string
}

var (
	trayPreferencesPath = defaultTrayPreferencesPath
	preferencesMu       sync.Mutex
)

func defaultTrayPreferencesPath() (string, error) {
	directory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(directory, "gswitch", "ui.conf"), nil
}

// LoadUIPreferences migrates an existing tray.conf only when ui.conf is absent.
// After successful creation or reading of ui.conf, the obsolete file is removed.
// An existing ui.conf always wins; failed migration leaves tray.conf intact.
func LoadUIPreferences() (UIPreferences, error) {
	preferencesMu.Lock()
	defer preferencesMu.Unlock()
	_, data, err := readUIPreferences()
	return parseUIPreferences(data), err
}

func parseUIPreferences(data []byte) UIPreferences {
	prefs := UIPreferences{IconMode: DefaultTrayIconMode, Language: "auto"}
	for line := range strings.SplitSeq(string(data), "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "tray-icon-mode":
			prefs.IconMode = normalizeTrayIconMode(TrayIconMode(value))
		case "ui-language":
			if value != "" {
				prefs.Language = value
			}
		}
	}
	return prefs
}

func readUIPreferences() (string, []byte, error) {
	path, err := trayPreferencesPath()
	if err != nil {
		return "", nil, err
	}
	// #nosec G304 -- user config directory or test-controlled path.
	data, err := os.ReadFile(path)
	if err == nil {
		return path, data, removeLegacyUIPreferences(path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return path, nil, err
	}
	legacy := filepath.Join(filepath.Dir(path), "tray.conf")
	if legacy == path {
		return path, nil, nil
	} // Compatibility with path-injected tests.
	// #nosec G304 -- fixed legacy filename in the user config directory.
	data, err = os.ReadFile(legacy)
	if errors.Is(err, os.ErrNotExist) {
		return path, nil, nil
	}
	if err != nil {
		return path, nil, err
	}
	if err := writeUIPreferences(path, data, true); err != nil {
		if errors.Is(err, os.ErrExist) {
			// Another instance created ui.conf during migration; never overwrite it.
			// #nosec G304 -- user config directory or test-controlled path.
			current, readErr := os.ReadFile(path)
			if readErr != nil {
				return path, current, readErr
			}
			return path, current, removeLegacyUIPreferences(path)
		}
		return path, data, fmt.Errorf("migrate interface preferences: %w", err)
	}
	return path, data, removeLegacyUIPreferences(path)
}

// removeLegacyUIPreferences runs only after the destination is available.
func removeLegacyUIPreferences(path string) error {
	legacy := filepath.Join(filepath.Dir(path), "tray.conf")
	if legacy == path {
		return nil // Compatibility with path-injected tests.
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		if _, err := os.Lstat(legacy); errors.Is(err, os.ErrNotExist) {
			return nil
		}
		// A manually linked ui.conf may still depend on the legacy pathname.
		return errors.New("cannot remove legacy preferences while ui.conf is a symbolic link")
	}
	// Unlink removes files and symlinks, never a directory at the legacy path.
	if err := unix.Unlink(legacy); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove obsolete interface preferences: %w", err)
	}
	return nil
}

// SaveUIPreferences atomically updates both preferences, retaining other keys.
func SaveUIPreferences(prefs UIPreferences) error {
	preferencesMu.Lock()
	defer preferencesMu.Unlock()
	return saveUIPreferences(prefs)
}

func saveUIPreferences(prefs UIPreferences) error {
	path, data, err := readUIPreferences()
	if err != nil {
		return err
	}
	if strings.ContainsAny(prefs.Language, "\r\n=") {
		return errors.New("invalid interface language")
	}
	if prefs.Language == "" {
		prefs.Language = "auto"
	}
	var lines []string
	for line := range strings.SplitSeq(strings.TrimRight(string(data), "\n"), "\n") {
		key, _, _ := strings.Cut(line, "=")
		if key = strings.TrimSpace(key); key != "tray-icon-mode" && key != "ui-language" && line != "" {
			lines = append(lines, line)
		}
	}
	lines = append(lines, "tray-icon-mode="+string(normalizeTrayIconMode(prefs.IconMode)), "ui-language="+prefs.Language)
	return writeUIPreferences(path, []byte(strings.Join(lines, "\n")+"\n"), false)
}

func writeUIPreferences(path string, data []byte, exclusive bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create interface preferences directory: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".ui.conf-*")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if exclusive {
		return os.Link(temporary, path)
	}
	return os.Rename(temporary, path)
}

// LoadTrayIconMode retains the existing default-on-read-error caller contract.
func LoadTrayIconMode() TrayIconMode { prefs, _ := LoadUIPreferences(); return prefs.IconMode }

// SaveTrayIconMode updates just the icon without discarding the language.
func SaveTrayIconMode(mode TrayIconMode) error {
	preferencesMu.Lock()
	defer preferencesMu.Unlock()
	_, data, err := readUIPreferences()
	if err != nil {
		return err
	}
	prefs := parseUIPreferences(data)
	prefs.IconMode = mode
	return saveUIPreferences(prefs)
}
