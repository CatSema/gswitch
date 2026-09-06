package tray

import (
	"os"
	"sync/atomic"

	"github.com/arumata/gswitch/internal/i18n"
)

// Catalogs are immutable. The pointer swaps atomically when Apply commits a
// language; tray workers and GTK callbacks never share mutable message maps.
var interfaceLanguage atomic.Pointer[i18n.Translator]

func currentLanguage() *i18n.Translator {
	if value := interfaceLanguage.Load(); value != nil {
		return value
	}
	return i18n.Builtin.Select("en", os.Getenv)
}

func tr(key string) string { return currentLanguage().Text(key) }

func selectInterfaceLanguage(preference string) {
	interfaceLanguage.Store(i18n.Builtin.Select(preference, os.Getenv))
}

func localizedLayoutName(name string) string {
	keys := map[string]string{
		"English (US)": strLayoutUS,
		"Russian":      strLayoutRU,
		"German":       strLayoutDE,
		"French":       strLayoutFR,
		"Spanish":      strLayoutES,
		"Italian":      strLayoutIT,
		"Portuguese":   strLayoutPT,
		"Polish":       strLayoutPL,
		"Ukrainian":    strLayoutUA,
		"Belarusian":   strLayoutBY,
		"Kazakh":       strLayoutKZ,
		"English (UK)": strLayoutGB,
		"Unknown":      strStatusUnknown,
	}
	if key, ok := keys[name]; ok {
		return tr(key)
	}
	return name // Names supplied by the desktop/input method remain verbatim.
}
