package tray

type TrayIconMode string

const (
	TrayIconModeApp         TrayIconMode = "app"
	TrayIconModeFlag        TrayIconMode = "flag"
	TrayIconModeAppWithFlag TrayIconMode = "app-with-flag"
	DefaultTrayIconMode     TrayIconMode = TrayIconModeFlag
)

func normalizeTrayIconMode(mode TrayIconMode) TrayIconMode {
	switch mode {
	case TrayIconModeApp, TrayIconModeFlag, TrayIconModeAppWithFlag:
		return mode
	default:
		return DefaultTrayIconMode
	}
}
