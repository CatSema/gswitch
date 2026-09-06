"""Independent UI client for the isolated KWin focus regression; no input injection."""
import gi
import sys
from pathlib import Path

gi.require_version('Gtk', '3.0')
gi.require_version('GdkX11', '3.0')
from gi.repository import Gtk, Gdk, GdkX11, GLib

window = Gtk.Window(title='Focus regression peer')
window.set_default_size(400, 300)
window.add_events(Gdk.EventMask.PROPERTY_CHANGE_MASK)
window.realize()
window.show_all()
window.present_with_time(GdkX11.x11_get_server_time(window.get_window()))

def ready():
    if not window.is_active():
        return True
    Path(sys.argv[1]).write_text('ready')
    return False

GLib.timeout_add(10, ready)
Gtk.main()
