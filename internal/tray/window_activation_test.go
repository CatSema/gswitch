package tray

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

// Run only through tools/scripts/test-settings-focus.sh, on its private display.
func TestSettingsRequestedFromTrayGetsFocus(t *testing.T) {
	if os.Getenv("GSWITCH_FOCUS_TEST") != "1" {
		t.Skip("requires isolated KWin display")
	}
	runtime.LockOSThread()
	gtk.Init(nil)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	bin := t.TempDir()
	for _, name := range []string{"systemctl", "gswitch", "pkexec"} {
		// #nosec G306 -- executable stub in test-owned storage.
		if err := os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\nexit 1\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)
	w, err := NewSettingsWindow(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.window.Destroy()
	w.deviceManager.sysDir = t.TempDir()
	w.window.ShowAll()
	w.window.PresentWithTime(settingsRequestTime(w.window))
	awaitSettingsFocus(t, "initial settings activation", w.window.IsActive)

	app := New()
	app.settingsWindow = w
	for _, state := range []string{"covered", "hidden", "minimized"} {
		t.Run(state, func(t *testing.T) {
			ready := filepath.Join(t.TempDir(), "peer-ready")
			// #nosec G204 -- fixed test helper; the only variable argument is a test-owned readiness path.
			peer := exec.Command("/usr/bin/python3", "testdata/focus/peer.py", ready)
			if err := peer.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = peer.Process.Kill()
				_ = peer.Wait()
			}()
			awaitSettingsFocus(t, "independent peer activation", func() bool {
				_, err := os.Stat(ready)
				return err == nil && !w.window.IsActive()
			})
			switch state {
			case "hidden":
				w.window.Hide()
			case "minimized":
				w.window.Iconify()
			}
			drainGTK()
			// The real tray callback arrives without a GTK key/button event.
			app.OnSettingsClicked()
			awaitSettingsFocus(t, "settings activation from tray", w.window.IsActive)
		})
	}
}

func awaitSettingsFocus(t *testing.T, description string, ready func() bool) {
	t.Helper()
	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		drainGTK()
		if ready() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", description)
}
