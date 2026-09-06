# Go release binaries are stripped explicitly; do not create empty debugsource RPMs.
%global debug_package %{nil}
Name:           gswitch
Version:        0.8.0
Release:        0
Summary:        Keyboard layout correction on explicit request
License:        MIT
URL:            https://github.com/arumata/gswitch
# Generated from the newest reachable public release tag by _service.
Source0:        gswitch-%{version}.tar.gz
Source1:        vendor.tar.gz
BuildRequires:  golang(API) >= 1.25
BuildRequires:  gcc
BuildRequires:  pkgconfig
BuildRequires:  pkgconfig(gtk+-3.0)
BuildRequires:  pkgconfig(x11)
BuildRequires:  pkgconfig(xkbcommon)
BuildRequires:  systemd-rpm-macros
BuildRequires:  desktop-file-utils
Requires:       systemd
Requires:       udev
Requires:       xkeyboard-config
Requires:       /usr/bin/pkexec
Requires:       hicolor-icon-theme >= 0.18
Recommends:     wl-clipboard
Recommends:     desktop-file-utils
ExclusiveArch:  x86_64

%description
Correct text typed in the wrong keyboard layout on an explicit keyboard
trigger. Includes the graphical-session user daemon and GTK tray/settings.
The service receives keyboard events and emits corrected input.
Device access is granted to the active local graphical session.

%prep
%setup -q
%setup -q -T -D -a 1
# Keep bundled Go dependency notices in the installed package.
mkdir third-party-licenses
find vendor -type f \( -name LICENSE -o -name NOTICE -o -name COPYING \) | while read -r file; do
    dest="third-party-licenses/${file#vendor/}"
    mkdir -p "$(dirname "$dest")"
    cp "$file" "$dest"
done

%build
export CGO_ENABLED=1
export GOTOOLCHAIN=local
export GOPROXY=off
export GOSUMDB=off
export GOFLAGS='-mod=vendor -buildvcs=false'
export GOCACHE="${GOCACHE:-$PWD/.gocache}"
go build -trimpath -ldflags '-s -w -X main.version=%{version}' -o build/gswitch ./cmd/gswitch
go build -trimpath -ldflags '-s -w -X main.version=%{version}' -o build/gswitch-tray ./cmd/gswitch-tray

%install
install -D -m 0755 build/gswitch %{buildroot}%{_bindir}/gswitch
install -D -m 0755 build/gswitch-tray %{buildroot}%{_bindir}/gswitch-tray
install -D -m 0644 configs/default.conf %{buildroot}%{_sysconfdir}/gswitch/default.conf
install -D -m 0644 configs/gswitch.service %{buildroot}%{_userunitdir}/gswitch.service
install -D -m 0644 configs/70-gswitch.rules %{buildroot}%{_udevrulesdir}/70-gswitch.rules
install -D -m 0644 configs/gswitch-tray.desktop %{buildroot}%{_sysconfdir}/xdg/autostart/gswitch-tray.desktop
install -D -m 0644 configs/gswitch-daemon.desktop %{buildroot}%{_sysconfdir}/xdg/autostart/gswitch-daemon.desktop
install -D -m 0644 configs/gswitch-tray.desktop %{buildroot}%{_datadir}/applications/gswitch-tray.desktop
install -D -m 0644 configs/com.github.arumata.gswitch.policy %{buildroot}%{_datadir}/polkit-1/actions/com.github.arumata.gswitch.policy
mkdir -p %{buildroot}%{_datadir}/icons
cp -a assets/icons/hicolor %{buildroot}%{_datadir}/icons/

%check
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off CGO_ENABLED=1
export GOFLAGS='-mod=vendor -buildvcs=false'
export GOCACHE="${GOCACHE:-$PWD/.gocache}"
go test -count=1 -trimpath ./...
desktop-file-validate configs/gswitch-tray.desktop configs/gswitch-daemon.desktop

# Embed the published upstream lifecycle hooks; do not enable a root service.
# User-manager/udev effects still need verification on a real Tumbleweed session.
%pre -f scripts/preinstall.sh
%post -f scripts/postinstall.sh
%preun -f scripts/preremove.sh
%postun -f scripts/postremove.sh

%files
%license LICENSE
%license third-party-licenses
%doc README.md
%{_bindir}/gswitch
%{_bindir}/gswitch-tray
%dir %{_sysconfdir}/gswitch
%config(noreplace) %{_sysconfdir}/gswitch/default.conf
%{_userunitdir}/gswitch.service
%{_udevrulesdir}/70-gswitch.rules
%config(noreplace) %{_sysconfdir}/xdg/autostart/gswitch-tray.desktop
%config(noreplace) %{_sysconfdir}/xdg/autostart/gswitch-daemon.desktop
%{_datadir}/applications/gswitch-tray.desktop
%{_datadir}/polkit-1/actions/com.github.arumata.gswitch.policy
%{_datadir}/icons/hicolor/*/apps/gswitch.*

%changelog
* Sun Sep 06 2026 arumata <arumata@gmail.com> - 0.8.0-0
- Initial experimental source package for Tumbleweed.
