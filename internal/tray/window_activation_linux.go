//nolint:gocritic // The cgo pseudo-package is misidentified as a duplicate package import.
package tray

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>
#include <gdk/gdkx.h>

static guint32 settings_request_time(GtkWidget *widget) {
    GdkDisplay *display = gtk_widget_get_display(widget);
    if (!GDK_IS_X11_DISPLAY(display))
        return gtk_get_current_event_time();

    // DBus tray clicks carry no GTK input event. Obtain an X server timestamp
    // before loading settings, not the stale timestamp of the last GTK event.
    gtk_widget_add_events(widget, GDK_PROPERTY_CHANGE_MASK);
    gtk_widget_realize(widget);
    return gdk_x11_get_server_time(gtk_widget_get_window(widget));
}
*/
import "C" //nolint:gocritic // C is the cgo pseudo-package, not a duplicate Go import.

import (
	"runtime"
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

// settingsRequestTime must run on the GTK thread at the start of a user request.
func settingsRequestTime(window *gtk.Window) uint32 {
	timestamp := uint32(C.settings_request_time((*C.GtkWidget)(unsafe.Pointer(window.GObject))))
	runtime.KeepAlive(window)
	return timestamp
}
