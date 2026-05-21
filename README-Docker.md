# Docker image for Ark Sync

Build the image from this repository (recommended). The container layout follows the upstream Syncthing Docker image conventions; the sync binary is **`arksync`** when built from current `build.go`.

## Build the image

```bash
docker build -t ark-sync:local .
```

For multi-arch or release-aligned builds, use the same `go run build.go` flow as in [doc/08无网页API-only打包指南.md](doc/08无网页API-only打包指南.md) and [.github/workflows/release-ark.yaml](.github/workflows/release-ark.yaml), then adjust the Dockerfile copy step if your artifact is named `arksync-linux-*` instead of `syncthing-linux-*`.

## Volumes and user

Use the `/var/syncthing` volume for synchronized data and configuration (path name kept for compatibility with the upstream image layout).

Ark Sync runs as **UID 1000** and **GID 1000** by default. Override with `PUID` and `PGID`. Set the container hostname with `--hostname=arksync` if needed.

Optional environment variables (same as upstream image):

- **`PCAP`** — extra capabilities without root, e.g. `PCAP=cap_chown,cap_fowner+ep`
- **`UMASK`** — e.g. `UMASK=002`
- **`STGUIADDRESS`** — GUI/API listen address; API-only deployments often use `STGUIADDRESS=` or restrict to localhost in config

## Example usage

**Docker CLI**

```bash
docker build -t ark-sync:local .
docker run --network=host -e STGUIADDRESS= \
  -v /path/to/data:/var/syncthing \
  ark-sync:local
```

**Docker Compose**

```yaml
services:
  arksync:
    image: ark-sync:local
    build: .
    container_name: arksync
    hostname: my-ark-sync
    environment:
      - PUID=1000
      - PGID=1000
      - STGUIADDRESS=
    volumes:
      - /path/to/data:/var/syncthing
    network_mode: host
    restart: unless-stopped
    healthcheck:
      test: curl -fkLsS -m 2 127.0.0.1:8384/rest/noauth/health | grep -o --color=never OK || exit 1
      interval: 1m
      timeout: 10s
      retries: 3
```

## Discovery and networking

Docker’s default bridge network hides real LAN addresses; Ark Sync inside the container may only see `172.17.0.0/16`, which hurts local discovery and transfer speed.

**Use [host network mode](https://docs.docker.com/network/host/)** when possible (as above).

Ark Sync controls which interfaces and ports it listens on; adjust configuration if there are conflicts.

## API / GUI security

By default the image may listen on `0.0.0.0:8384`. The example clears `STGUIADDRESS` so the process uses the config file or GUI settings (often `127.0.0.1`).

If you expose the API externally:

- Enable authentication (API key).
- Enable TLS.
- Prefer probing **`/rest/noauth/health`** for health checks (see [doc/08无网页API-only打包指南.md](doc/08无网页API-only打包指南.md)).

This fork often ships **without bundled Web GUI assets** (`noassets`); plan on REST API or an external client.
