# Ark Sync

[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](LICENSE)

**Ark Sync** is a continuous file synchronization program, forked from [Syncthing](https://github.com/syncthing/syncthing) and customized for the Ark stack. It synchronizes files between two or more computers.

This repository ships the **`arksync`** binary (API-only / `noassets` builds by default). See [doc/08无网页API-only打包指南.md](doc/08无网页API-only打包指南.md) and [doc/01Syncthing完整指南.md](doc/01Syncthing完整指南.md) for usage in Chinese.

## Goals

We strive to fulfill the goals below (see [GOALS.md](GOALS.md) for the full commentary). Ark Sync should be:

1. **Safe From Data Loss** — Protecting user data is paramount.
2. **Secure Against Attackers** — Data must not be exposed to unauthorized parties.
3. **Easy to Use** — Approachable and understandable (often via REST API or a host client).
4. **Automatic** — User interaction only when necessary.
5. **Universally Available** — Runs on common desktop and server platforms.
6. **For Individuals** — Empowering users with safe, secure synchronization.
7. **Everything Else** — Without conflicting with the goals above.

## Getting Started

- Build and run locally: [doc/08无网页API-only打包指南.md](doc/08无网页API-only打包指南.md)
- REST API: [doc/02REST_API参考.md](doc/02REST_API参考.md)
- Example systemd and desktop files: [etc/](etc/)

The main process is intended to be launched by **Ark Sync Client** (`arksync_client`); direct runs from a shell require `ARKSYNC_SKIP_PARENT_CHECK=1` unless you use the approved parent process.

## Docker

To run Ark Sync in Docker, see [README-Docker.md](README-Docker.md).

## Releases

Prebuilt binaries for Linux, Windows, and macOS are published via GitHub Actions when you push a version tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Workflow: [.github/workflows/release-ark.yaml](.github/workflows/release-ark.yaml). Asset names look like `arksync-linux-amd64-v1.0.0-noupgrade.tar.gz`.

Automatic in-app upgrade is disabled on the **`noupgrade`** branch (`-no-upgrade` / `noupgrade` build tag).

## Getting in Touch

- **Issues**: [github.com/arkcore-share/ark-sync/issues](https://github.com/arkcore-share/ark-sync/issues)
- Upstream Syncthing forum/docs remain useful for protocol behavior: [docs.syncthing.net](https://docs.syncthing.net/) (this fork may differ in GUI, upgrades, and packaging).

## Building

From a git checkout on branch **`noupgrade`** (or your development branch):

```bash
# API-only binary
go run build.go -no-upgrade -tags "noassets" -build-out ./bin/arksync build syncthing

# Linux release tarball (needs CGO / toolchain as in CI)
go run build.go -no-upgrade -tags "sqlite_omit_load_extension sqlite_dbstat noassets" tar syncthing
```

## License

All code is licensed under the [Mozilla Public License 2.0](LICENSE).
