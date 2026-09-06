package tray

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>
*/
import "C"

// gotk3 does not bind gtk_widget_set_default_direction. Call on the GTK thread
// before constructing translated widgets, including dialogs and popovers.
func setInterfaceDirection() {
	direction := C.GTK_TEXT_DIR_LTR
	if currentLanguage().Direction() == "rtl" {
		direction = C.GTK_TEXT_DIR_RTL
	}
	C.gtk_widget_set_default_direction(C.GtkTextDirection(direction))
}
