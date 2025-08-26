# Test Environment

This folder spins up a 3-Traefik setup to exercise the provider plugin:

- traefik-main: loads the local plugin and discovers config from provider1/provider2
- traefik-provider1, traefik-provider2: expose dynamic configs under /api/rawdata (file provider)
- whoami, tcp-echo, udp-echo, redis: backends for HTTP/TCP/UDP

Key ideas covered:
- matchers-based discovery for routers/services/middlewares
- overrides on routers/services (rules, entrypoints, middlewares, servers, healthchecks)
- tunneling: replace service backends with "tunnel" addresses that point at other Traefik entrypoints

## Files
- test/docker-compose.yml
- test/configs/traefik-main/traefik.yml
- test/configs/traefik-provider1/dynamic.yml
- test/configs/traefik-provider2/dynamic.yml
- test/scripts/smoke.sh

## Run

Bring the environment up:

```bash
docker compose -f test/docker-compose.yml up -d
```

Run a simple smoke test (HTTP/TCP/UDP):

```bash
bash test/scripts/smoke.sh
```

Tear down:

```bash
docker compose -f test/docker-compose.yml down -v
```

## What to Expect
- HTTP discovery via plugin: routes defined in provider1/2 dynamic.yml should be reachable via traefik-main
- Tunnels: for selected services (e.g., provider1-service) traffic will be sent via traefik-provider1/2 entrypoints defined in traefik-main config
- TCP and UDP routing examples exposed on entrypoints tcp-ep and udp-ep
