# Experimental package recipe for openSUSE Tumbleweed

This recipe targets x86_64 in an OBS home project. It builds the newest gswitch release tag reachable from the default branch, currently v0.8.0. The package is experimental.

## Prepare the sources

You need an OBS account and a home project with an openSUSE Tumbleweed build target. Install the source preparation tools:

```sh
sudo zypper install osc git go obs-service-obs_scm obs-service-tar obs-service-set_version obs-service-recompress obs-service-go_modules
```

Configure `osc` for your account. Create a regular OBS-managed package, not an SCM-synchronized package. Replace `YOUR_LOGIN` with the OBS project owner's login:

```sh
osc meta pkg -e home:YOUR_LOGIN gswitch
```

The command opens the package XML in an editor and creates the package if it does not exist. Keep the generated package metadata, but do not add a `<scmsync>` element or a Git repository URL. The `_service` file below already fetches the upstream Git sources. If the package metadata currently contains `<scmsync>`, remove that element before continuing; otherwise `osc commit` will skip the package as SCM-managed.

Check out the empty package:

```sh
osc checkout home:YOUR_LOGIN gswitch
cd home:YOUR_LOGIN/gswitch
```

Copy `_service`, `gswitch.spec`, and `gswitch.changes` from this repository folder into the working copy, then run:

```sh
osc service manualrun
osc addremove
```

The `obs_scm` service finds the newest reachable tag matching `v*`, checks out that exact tag, and records its version in OBS metadata. `set_version` reads that metadata and updates the spec, `tar` and `recompress` create the versioned source archive, and `go_modules` reads the generated `gswitch-*.obscpio` rather than guessing among tar archives. The local `third_party/clipboard` replacement remains in the sources and vendor directory. Source preparation needs network access. RPM compilation and Go tests use vendor without network access.

## Submit to your OBS project

Add the recipe and generated service files:

```sh
osc addremove
osc status
osc diff
```

Review the changes before running `osc commit`. Do not add the `gswitch/` Git cache directory left by the service. After submission, check the server build result. To investigate a failure, provide the full build log, OBS package URL, and source revision.

The `manual` mode does not start by itself when a release appears. For an update, verify that the new public release tag is reachable from the default branch and add a `.changes` entry. Remove the previous generated outputs before every run so wildcard inputs cannot see more than one release:

```sh
rm -f gswitch-*.obscpio gswitch-*.tar.gz gswitch.obsinfo vendor.tar.gz
osc service manualrun
osc addremove
osc status
osc diff
```

No version or revision edit is needed: the services select the newest matching tag, derive the version from OBS metadata, update the spec, and regenerate the source and vendor archives. The checked-out tag, rather than the development branch HEAD, remains the reproducible source boundary.

## Verification status

The local OBS services ran through the installed osc dispatcher in a Tumbleweed container. The RPM built without network access; Go tests and desktop file validation passed. Installation and binary dependencies were checked. A modified configuration survived reinstallation of the same version and remained as `.rpmsave` after removal.

The `osc service manualrun` command in an authenticated OBS working copy and the server build have not been verified. Container checks also do not establish that the tray, user service, or text correction work in a Tumbleweed graphical session. Verification of udev ACLs and the polkit dialog remains open. Official Tumbleweed support is not yet claimed.

After a successful build, check the user service and tray in your graphical session. Test word correction and configuration persistence after logging in again. Run the daemon as the graphical-session user, never as root.

## Packaging findings

rpmlint reports `polkit-untracked-privilege`: the action `com.github.arumata.gswitch.write-config` is absent from the polkit-default-privs profiles. The policy requires administrator authentication and permits temporary retention of that authorization for an active session (`auth_admin_keep`). No filter has been added for this finding.

openSUSE deliberately assigns low badness to these findings in home/devel projects. A first home-project build does not need to wait for an audit. Factory inclusion requires a security team audit and whitelisting of the action. Always check the full build log for your project.

The remaining rpmlint warnings concern PIE for both binaries, local source archive names, and identical Go module license texts. The spec retains lifecycle hooks from the published release, including migration of the old root service. Inclusion in the official openSUSE repository will require a separate packaging review.

## Documentation

OBS: https://openbuildservice.org/help/manuals/obs-user-guide/art-obs-bg

Creating package metadata with osc: https://openbuildservice.org/help/manuals/obs-user-guide/cha-obs-osc

SCM synchronization is a separate workflow: https://openbuildservice.org/help/manuals/obs-user-guide/cha-obs-scm-bridge

Go dependencies: https://github.com/openSUSE/obs-service-go_modules

Home/devel rules: https://en.opensuse.org/openSUSE:Package_security_guidelines#Rpmlint_whitelisting_errors_in_home_and_devel_Projects
