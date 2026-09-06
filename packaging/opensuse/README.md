# Experimental package recipe for openSUSE Tumbleweed

This recipe targets x86_64 in an OBS home project. It builds gswitch v0.8.0 from the published commit `9058e19e109507b1db0d85217cb6a7fa5dc9b3a1`. The package is experimental.

## Prepare the sources

You need an OBS account, a home project with an openSUSE Tumbleweed build target, and a `gswitch` package in that project. Install the source preparation tools:

```sh
sudo zypper install osc git go obs-service-tar_scm obs-service-recompress obs-service-go_modules
```

Configure `osc` for your account and check out the package. Replace `YOUR_LOGIN` with the OBS project owner's login:

```sh
osc checkout home:YOUR_LOGIN gswitch
cd home:YOUR_LOGIN/gswitch
```

Copy `_service`, `gswitch.spec`, and `gswitch.changes` from this repository folder into the working copy, then run:

```sh
osc service manualrun
```

The `tar_scm` service fetches the pinned commit from GitHub. Then `recompress` creates `gswitch-0.8.0.tar.gz`, and `go_modules` prepares dependencies in `gswitch-0.8.0-vendor.tar.gz`. The local `third_party/clipboard` replacement remains in the sources and vendor directory. Source preparation needs network access. RPM compilation and Go tests use vendor without network access.

## Submit to your OBS project

Add the recipe and generated archives:

```sh
osc add _service gswitch.spec gswitch.changes gswitch-0.8.0.tar.gz gswitch-0.8.0-vendor.tar.gz
osc status
osc diff
```

Review the changes before running `osc commit`. Do not add the `gswitch/` Git cache directory left by the service. After submission, check the server build result. To investigate a failure, provide the full build log, OBS package URL, and source revision.

The `manual` mode does not update the package automatically when a release appears. For an update, verify the tag's public commit. Update revision/version and archive names in `_service`, source_commit/Version in the spec, and add a `.changes` entry. Repeat source preparation and verification.

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

Go dependencies: https://github.com/openSUSE/obs-service-go_modules

Home/devel rules: https://en.opensuse.org/openSUSE:Package_security_guidelines#Rpmlint_whitelisting_errors_in_home_and_devel_Projects
