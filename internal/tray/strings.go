package tray

// Stable message IDs. Translations live in internal/i18n/locales.

// Window titles
const (
	strWindowTitle = "strWindowTitle"
)

// Section titles
const (
	strSectionService = "strSectionService"
	strSectionKeys    = "strSectionKeys"
	strSectionLayouts = "strSectionLayouts"
	strSectionTray    = "strSectionTray"
	strSectionDelays  = "strSectionDelays"
	strSectionDevices = "strSectionDevices"
)

// Service section
const (
	strLabelStatus        = "strLabelStatus"
	strStatusUnknown      = "strStatusUnknown"
	strStatusRunning      = "strStatusRunning"
	strStatusStopped      = "strStatusStopped"
	strStatusFailed       = "strStatusFailed"
	strStatusNotInstalled = "strStatusNotInstalled"
	strButtonRestart      = "strButtonRestart"
	strButtonStart        = "strButtonStart"
	strButtonStop         = "strButtonStop"
	strCheckAutostart     = "strCheckAutostart"
)

// Tray menu - service status
const (
	strTrayServiceRunning      = "strTrayServiceRunning"
	strTrayServiceStopped      = "strTrayServiceStopped"
	strTrayServiceFailed       = "strTrayServiceFailed"
	strTrayServiceNotInstalled = "strTrayServiceNotInstalled"
	strTrayServiceUnknown      = "strTrayServiceUnknown"
)

// Tray tooltips and detection status
const (
	strTooltipOK             = "strTooltipOK" // e.g. "gswitch: Alt+Shift (xkb)"
	strTooltipNeedsConfig    = "strTooltipNeedsConfig"
	strTooltipServiceError   = "strTooltipServiceError"
	strTooltipServiceStopped = "strTooltipServiceStopped" // For tooltip composition (without "gswitch:" prefix)
	strTooltipServiceFailed  = "strTooltipServiceFailed"  // For tooltip composition (without "gswitch:" prefix)
	strTooltipDetectError    = "strTooltipDetectError"
	strTooltipClickConfig    = "strTooltipClickConfig"
	strTrayRefresh           = "strTrayRefresh"
)

// First-run dialog (GNOME)
const (
	strFirstRunTitle   = "strFirstRunTitle"
	strFirstRunMessage = "strFirstRunMessage"
	strFirstRunConfig  = "strFirstRunConfig"
	strFirstRunDismiss = "strFirstRunDismiss"
)

// Error messages
const (
	strErrorTitle          = "strErrorTitle"
	strErrorRestartFailed  = "strErrorRestartFailed"
	strErrorStartFailed    = "strErrorStartFailed"
	strErrorStopFailed     = "strErrorStopFailed"
	strErrorEnableFailed   = "strErrorEnableFailed"
	strErrorDisableFailed  = "strErrorDisableFailed"
	strErrorValidation     = "strErrorValidation"
	strErrorSaveFailed     = "strErrorSaveFailed"
	strErrorTraySaveFailed = "strErrorTraySaveFailed"
)

// Warning messages
const (
	strWarnRestartFailed = "strWarnRestartFailed"
)

// Keys section
const (
	strLabelLayoutSwitch           = "strLabelLayoutSwitch"
	strLabelConvertKey             = "strLabelConvertKey"
	strLayoutSwitchHint            = "strLayoutSwitchHint"
	strLabelConversionModifiers    = "strLabelConversionModifiers"
	strConversionShortcutsTitle    = "strConversionShortcutsTitle"
	strShortcutsChanged            = "strShortcutsChanged"
	strShortcutsTyped              = "strShortcutsTyped"
	strShortcutsSelected           = "strShortcutsSelected"
	strShortcutsWord               = "strShortcutsWord"
	strShortcutsLine               = "strShortcutsLine"
	strShortcutsLayout             = "strShortcutsLayout"
	strShortcutsCase               = "strShortcutsCase"
	strShortcutsUndoTitle          = "strShortcutsUndoTitle"
	strShortcutsUndo               = "strShortcutsUndo"
	strShortcutsDoubleTap          = "strShortcutsDoubleTap"
	strShortcutsOtherShift         = "strShortcutsOtherShift"
	strShortcutsInfoIcon           = "strShortcutsInfoIcon"
	strConversionModifiersStandard = "strConversionModifiersStandard"
	strConversionModifiersPunto    = "strConversionModifiersPunto"
	strAutoDetect                  = "strAutoDetect"
)

// Layout switch options
// Order: Auto, Alt+Shift, Ctrl+Shift, Super+Space, Caps Lock, Custom
// Note: "Super (Win)" removed - single key is ambiguous for layout switching
var layoutSwitchOptions = []struct {
	Label string
	Value string
}{
	{strAutoDetect, "auto"},
	{"LAlt+LShift (56+42)", "56+42"},
	{"LCtrl+LShift (29+42)", "29+42"},
	{"LSuper+Space (125+57)", "125+57"},
	{"Caps Lock (58)", "58"},
	{strOther, "custom"},
}

// Convert key options
var convertKeyOptions = []struct {
	Label string
	Value string
}{
	{strDoubleShift, "0"},
	{"Pause/Break", "119"},
	{"Scroll Lock", "70"},
	{strOther, "custom"},
}

// Layouts section
const (
	strCheckAutoDetect = "strCheckAutoDetect"
	strLabelLayout1    = "strLabelLayout1"
	strLabelLayout2    = "strLabelLayout2"
)

// Layout options
var layoutOptions = []string{
	"us", "ru", "ua", "de", "fr", "es", "it", "pt", "pl", "cz", "sk", "hu",
}

// Tray section
const strLabelTrayIcon = "strLabelTrayIcon"

var trayIconModeOptions = []struct {
	Label string
	Value TrayIconMode
}{
	{strIconApp, TrayIconModeApp},
	{strIconAppFlag, TrayIconModeAppWithFlag},
	{strIconFlag, TrayIconModeFlag},
}

// Delays section
const (
	strLabelDelayBetween = "strLabelDelayBetween"
	strLabelDelaySwitch  = "strLabelDelaySwitch"
	strLabelMs           = "strLabelMs"
)

// Devices section
const (
	strDevicesLoading    = "strDevicesLoading"
	strDevicesNoAccess   = "strDevicesNoAccess"
	strDevicesEmpty      = "strDevicesEmpty"
	strDevicesAllBlocked = "strDevicesAllBlocked"
	strDeviceBlocked     = "strDeviceBlocked"
)

// Buttons
const (
	strButtonCancel = "strButtonCancel"
	strButtonApply  = "strButtonApply"
	strButtonOK     = "strButtonOK"
)

// Key picker dialog
const (
	strKeyPickerTitleLayoutSwitch   = "strKeyPickerTitleLayoutSwitch"
	strKeyPickerTitleConvertKey     = "strKeyPickerTitleConvertKey"
	strKeyPickerInstruction         = "strKeyPickerInstruction"
	strKeyPickerScancode            = "strKeyPickerScancode"
	strKeyPickerHint                = "strKeyPickerHint"
	strKeyPickerHintConvert         = "strKeyPickerHintConvert"
	strConversionKeyRecovery        = "strConversionKeyRecovery"
	strKeyPickerModifierRejected    = "strKeyPickerModifierRejected"
	strKeyPickerCombinationRejected = "strKeyPickerCombinationRejected"

	strKeyPickerCurrentLayoutSwitch = "strKeyPickerCurrentLayoutSwitch"
	strKeyPickerCurrentConvertKey   = "strKeyPickerCurrentConvertKey"
)

// Detection status (settings window)
const (
	strDetecting          = "strDetecting"
	strDetectedOK         = "strDetectedOK"
	strDetectedWarning    = "strDetectedWarning"
	strDetectNeedsConfig  = "strDetectNeedsConfig"
	strDetectError        = "strDetectError"
	strDetectErrorUnknown = "strDetectErrorUnknown"
)

// Auto-detect validation dialog (shown when saving with failed detection)
const (
	strAutoDetectFailTitle   = "strAutoDetectFailTitle"
	strAutoDetectFailReason  = "strAutoDetectFailReason"
	strAutoDetectFailDefault = "strAutoDetectFailDefault"
	strAutoDetectFailSave    = "strAutoDetectFailSave"
	strAutoDetectFailManual  = "strAutoDetectFailManual"
	strAutoDetectFailSources = "strAutoDetectFailSources"
)

const (
	strSectionInterface    = "strSectionInterface"
	strLabelLanguage       = "strLabelLanguage"
	strLanguageAuto        = "strLanguageAuto"
	strOther               = "strOther"
	strDoubleShift         = "strDoubleShift"
	strIconApp             = "strIconApp"
	strIconAppFlag         = "strIconAppFlag"
	strIconFlag            = "strIconFlag"
	strMenuSettings        = "strMenuSettings"
	strMenuQuit            = "strMenuQuit"
	strMenuServiceHint     = "strMenuServiceHint"
	strMenuActionHint      = "strMenuActionHint"
	strMenuRefreshHint     = "strMenuRefreshHint"
	strMenuSettingsHint    = "strMenuSettingsHint"
	strMenuQuitHint        = "strMenuQuitHint"
	strInvalidTrayIcon     = "strInvalidTrayIcon"
	strInvalidLanguage     = "strInvalidLanguage"
	strInvalidLayoutSwitch = "strInvalidLayoutSwitch"
	strInvalidConvertKey   = "strInvalidConvertKey"
	strInvalidLayout1      = "strInvalidLayout1"
	strInvalidLayout2      = "strInvalidLayout2"
	strInvalidDelay        = "strInvalidDelay"
	strInvalidSwitchDelay  = "strInvalidSwitchDelay"
	strMissingLayoutSwitch = "strMissingLayoutSwitch"
	strMissingConvertKey   = "strMissingConvertKey"
	strInvalidScancode     = "strInvalidScancode"
	strErrorUISaveFailed   = "strErrorUISaveFailed"
	strErrorUILoadFailed   = "strErrorUILoadFailed"
)

const (
	strAppTooltip             = "strAppTooltip"
	strGswitchMissing         = "strGswitchMissing"
	strResponseInvalid        = "strResponseInvalid"
	strGswitchRunFailed       = "strGswitchRunFailed"
	strConvertModifierInvalid = "strConvertModifierInvalid"
)

const (
	strCustomKey = "strCustomKey"
	strLayoutUS  = "strLayoutUS"
	strLayoutRU  = "strLayoutRU"
	strLayoutDE  = "strLayoutDE"
	strLayoutFR  = "strLayoutFR"
	strLayoutES  = "strLayoutES"
	strLayoutIT  = "strLayoutIT"
	strLayoutPT  = "strLayoutPT"
	strLayoutPL  = "strLayoutPL"
	strLayoutUA  = "strLayoutUA"
	strLayoutBY  = "strLayoutBY"
	strLayoutKZ  = "strLayoutKZ"
	strLayoutGB  = "strLayoutGB"
)
