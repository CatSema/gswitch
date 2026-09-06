package tray

import (
	"bytes"
	"image/png"
	"testing"
)

func TestNormalIconSelection(t *testing.T) {
	if got := GetNormalIcon(TrayIconModeApp, "US"); !bytes.Equal(got, appIcon) {
		t.Fatal("app mode did not return the embedded application icon")
	}
	flag, err := png.Decode(bytes.NewReader(GetNormalIcon(TrayIconModeFlag, "US")))
	if err != nil {
		t.Fatalf("decode display flag: %v", err)
	}
	if got := flag.Bounds().Size(); got.X != iconSize || got.Y != 33 {
		t.Fatalf("display flag size = %v, want 48x33", got)
	}
	if got := GetNormalIcon(TrayIconModeFlag, "ZZ"); !bytes.Equal(got, appIcon) {
		t.Fatal("unknown layout did not fall back to the application icon")
	}
}

func TestOverlayIconContainsApplicationAndFlag(t *testing.T) {
	overlay, err := png.Decode(bytes.NewReader(GetNormalIcon(TrayIconModeAppWithFlag, "US")))
	if err != nil {
		t.Fatalf("decode overlay: %v", err)
	}
	base, err := png.Decode(bytes.NewReader(appIcon))
	if err != nil {
		t.Fatalf("decode application icon: %v", err)
	}
	flag, err := png.Decode(bytes.NewReader(GetFlagIcon("US")))
	if err != nil {
		t.Fatalf("decode flag icon: %v", err)
	}
	scaledFlag := scaleNearest(flag, 22, 15)

	if !sameColor(overlay.At(1, 1), base.At(1, 1)) {
		t.Fatal("overlay does not retain the application icon outside the flag")
	}
	if !sameColor(overlay.At(iconSize-1, iconSize-1), scaledFlag.At(21, 14)) {
		t.Fatal("overlay does not contain the flag in the bottom-right corner")
	}
}

func sameColor(left, right interface {
	RGBA() (uint32, uint32, uint32, uint32)
}) bool {
	lr, lg, lb, la := left.RGBA()
	rr, rg, rb, ra := right.RGBA()
	return lr == rr && lg == rg && lb == rb && la == ra
}

func TestOfferedLayoutsHaveFlags(t *testing.T) {
	for _, code := range []string{"us", "ru", "de", "fr", "es", "it", "pt", "pl", "ua", "by", "kz", "gb"} {
		t.Run(code, func(t *testing.T) {
			if !HasFlagIcon(code) {
				t.Fatal("offered layout has no embedded flag")
			}
			flag, err := png.Decode(bytes.NewReader(GetFlagIcon(code)))
			if err != nil {
				t.Fatal(err)
			}
			if size := flag.Bounds().Size(); size.X != 16 || size.Y != 11 {
				t.Fatalf("flag size = %v, want consistent 16x11 assets", size)
			}
			for _, mode := range []TrayIconMode{TrayIconModeFlag, TrayIconModeAppWithFlag} {
				data := GetNormalIcon(mode, code)
				if bytes.Equal(data, appIcon) {
					t.Fatal("known layout fell back to application icon")
				}
				if _, err := png.Decode(bytes.NewReader(data)); err != nil {
					t.Fatalf("invalid flag mode icon: %v", err)
				}
			}
		})
	}
}
