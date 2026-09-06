package tray

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/arumata/gswitch/internal/i18n"
)

type translatedMenuItem struct{ title, tooltip string }

func (m *translatedMenuItem) Clicks() <-chan struct{} { return nil }
func (m *translatedMenuItem) SetTitle(s string)       { m.title = s }
func (m *translatedMenuItem) SetTooltip(s string)     { m.tooltip = s }
func (m *translatedMenuItem) Enable()                 {}
func (m *translatedMenuItem) Disable()                {}

func TestLanguageRefreshesUnchangedTrayStatus(t *testing.T) {
	old := interfaceLanguage.Load()
	defer interfaceLanguage.Store(old)
	menu := &translatedMenuItem{}
	status := &translatedMenuItem{}
	action := &translatedMenuItem{}
	backend := &iconModeRecordingTrayBackend{}
	tray := &Tray{backend: backend, mSettings: menu, mServiceStatus: status, mServiceAction: action, currentServiceStatus: StatusRunning,
		detectionInfo: DetectionInfo{Status: TrayStatusServiceError, Error: "service is not running", ErrorMessageID: strTooltipServiceStopped}}
	selectInterfaceLanguage("ru")
	tray.applyLanguage()
	if menu.title != "Настройки..." || status.title != "● Служба работает" || action.title != "Остановить" {
		t.Fatalf("untranslated menu: %+v %+v %+v", menu, status, action)
	}
	if backend.tooltip != "gswitch: служба не работает" {
		t.Fatalf("cached error remained untranslated: %s", backend.tooltip)
	}
	selectInterfaceLanguage("en")
	tray.applyLanguage()
	if menu.title != "Settings..." || action.title != "Stop" {
		t.Fatal("reverse language change failed")
	}
}

func TestInterfaceApplyAndCancel(t *testing.T) {
	if os.Getenv("GSWITCH_GTK_TEST") != "1" {
		t.Skip("requires isolated GTK display")
	}
	runtime.LockOSThread()
	gtk.Init(nil)
	old := interfaceLanguage.Load()
	defer interfaceLanguage.Store(old)
	for _, tt := range []struct{ language, title, conversion string }{
		{"ru", "gswitch - Настройки", "Двойное нажатие Shift"},
		{"de", "gswitch - Einstellungen", "Shift zweimal drücken"},
		{"fr", "gswitch - Paramètres", "Deux pressions sur Shift"},
		{"es", "gswitch - Configuración", "Doble pulsación de Shift"},
		{"uk", "gswitch - Налаштування", "Подвійне натискання Shift"},
		{"pl", "gswitch - Ustawienia", "Podwójny Shift"},
		{"pt", "gswitch - Definições", "Shift duas vezes"},
		{"it", "gswitch - Impostazioni", "Doppio Shift"},
		{"be", "gswitch - Налады", "Падвойнае націсканне Shift"},
		{"kk", "gswitch - Баптаулар", "Shift пернесін екі рет басу"},
	} {
		checkInterfaceApplyAndCancel(t, tt.language, tt.title, tt.conversion)
	}
}

func checkInterfaceApplyAndCancel(t *testing.T, language, wantTitle, wantConvert string) {
	t.Helper()
	selectInterfaceLanguage("en")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	// Any accidental systemctl/pkexec invocation is recorded without reaching the host.
	bin := t.TempDir()
	log := filepath.Join(bin, "called")
	for _, name := range []string{"systemctl", "pkexec"} {
		// #nosec G306 -- executable stub in a test-owned temporary directory.
		if err := os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\nprintf '%s\\n' called >> '"+log+"'\nexit 1\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)
	w, err := NewSettingsWindow(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.window.Destroy()
	cfg := DefaultTrayConfig()
	cfg.LayoutSwitch = "125"
	w.loadKeyConfig(cfg)
	w.setLayoutComboByValue(w.layout1Combo, cfg.Layout1)
	w.setLayoutComboByValue(w.layout2Combo, cfg.Layout2)
	w.delayBetweenSpin.SetValue(float64(cfg.Delay))
	w.delaySwitchSpin.SetValue(float64(cfg.LayoutSwitchDelay))
	w.setTrayIconMode(TrayIconModeAppWithFlag)
	w.setLanguageSelection(language)
	w.onCancelClicked()
	drainGTK()
	prefs, err := LoadUIPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.Language != "auto" || currentLanguage().Code() != "en" {
		t.Fatal("Cancel applied language")
	}
	w.onApplyClicked()
	deadline := time.Now().Add(3 * time.Second)
	for currentLanguage().Code() != language && time.Now().Before(deadline) {
		drainGTK()
		runtime.Gosched()
	}
	drainGTK()
	if currentLanguage().Code() != language {
		t.Fatal("Apply did not update language")
	}
	prefs, err = LoadUIPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.Language != language || prefs.IconMode != TrayIconModeAppWithFlag {
		t.Fatalf("Apply did not persist UI preferences: %+v", prefs)
	}
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		t.Fatal("UI-only Apply invoked systemctl or pkexec")
	}
	translated, err := NewSettingsWindow(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer translated.window.Destroy()
	title, err := translated.window.GetTitle()
	if err != nil {
		t.Fatal(err)
	}
	if title != wantTitle {
		t.Fatalf("reopened title = %s", title)
	}
	if got := translated.convertKeyCombo.GetActiveText(); got != wantConvert {
		t.Fatalf("untranslated conversion choice %q", got)
	}
	// A write failure must not switch the live UI language.
	path, err := defaultTrayPreferencesPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := w.saveInterfacePreferences(UIPreferences{Language: "en", IconMode: TrayIconModeFlag}); err == nil {
		t.Fatal("invalid storage accepted")
	}
	if currentLanguage().Code() != language {
		t.Fatal("failed save changed language")
	}
}

func drainGTK() {
	for gtk.EventsPending() {
		gtk.MainIterationDo(false)
	}
}

// Optional artifacts render only fixture widgets on an isolated display.
func TestInterfaceLayoutCaptures(t *testing.T) {
	directory := os.Getenv("GSWITCH_UI_ARTIFACT_DIR")
	if os.Getenv("GSWITCH_GTK_TEST") != "1" {
		t.Skip("requires isolated display")
	}
	if directory == "" {
		directory = t.TempDir()
	}
	runtime.LockOSThread()
	gtk.Init(nil)
	old := interfaceLanguage.Load()
	defer interfaceLanguage.Store(old)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, catalog := range i18n.Builtin.Languages() {
		language := catalog.Code
		selectInterfaceLanguage(language)
		w, err := NewSettingsWindow(nil)
		if err != nil {
			t.Fatal(err)
		}
		cfg := DefaultTrayConfig()
		cfg.LayoutSwitch = "125"
		w.loadKeyConfig(cfg)
		w.setLayoutComboByValue(w.layout1Combo, "us")
		w.setLayoutComboByValue(w.layout2Combo, "ru")
		w.setLanguageSelection(language)
		w.setTrayIconMode(TrayIconModeFlag)
		w.window.ShowAll()
		captureFixture(t, filepath.Join(directory, language+"-settings.png"))
		w.conversionShortcutsPopover.Show()
		captureFixture(t, filepath.Join(directory, language+"-shortcuts.png"))
		x, y, err := w.conversionShortcutsPopover.TranslateCoordinates(w.window, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if right := x + w.conversionShortcutsPopover.GetAllocatedWidth(); x < 0 || right > w.window.GetAllocatedWidth() {
			t.Errorf("%s help clipped horizontally: x=%d right=%d window=%d", language, x, right, w.window.GetAllocatedWidth())
		}
		if bottom := y + w.conversionShortcutsPopover.GetAllocatedHeight(); bottom > w.window.GetAllocatedHeight() {
			t.Errorf("%s help clipped: bottom=%d window=%d", language, bottom, w.window.GetAllocatedHeight())
		}
		w.window.Destroy()
		drainGTK()
	}
}

func captureFixture(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(150 * time.Millisecond)
	for time.Now().Before(deadline) {
		drainGTK()
		runtime.Gosched()
	}
	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		t.Fatal(err)
	}
	window, err := screen.GetRootWindow()
	if err != nil {
		t.Fatal(err)
	}
	pixbuf, err := window.PixbufGetFromWindow(0, 0, window.WindowGetWidth(), window.WindowGetHeight())
	if err != nil {
		t.Fatal(err)
	}
	if err := pixbuf.SavePNG(path, 6); err != nil {
		t.Fatal(err)
	}
}

func TestMessageIDsHaveTranslations(t *testing.T) {
	source, err := parser.ParseFile(token.NewFileSet(), "strings.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(source, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range spec.Names {
			if !strings.HasPrefix(name.Name, "str") {
				continue
			}
			literal, ok := spec.Values[i].(*ast.BasicLit)
			if !ok {
				t.Fatalf("message %s has no literal ID", name.Name)
			}
			key, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatal(err)
			}
			if tr(key) == key {
				t.Errorf("message %s has no catalog entry", key)
			}
		}
		return true
	})
}
