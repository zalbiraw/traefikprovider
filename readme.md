# Traefik Provider: HTTP/TCP/UDP/TLS Merger with Matchers, Overrides, and Tunnels

This project is a Traefik Provider plugin that dynamically fetches, filters, merges, and enriches configurations from one or more upstream sources. It focuses on:

- Discovering HTTP/TCP/UDP/TLS resources via matchers
- Applying per-resource overrides
- Injecting HTTP/TCP tunnels (with optional mTLS) and wiring services to ServersTransports
- Producing a single merged Traefik dynamic configuration payload

It is designed to work both in Traefik plugin mode and in local development mode.

## Features

- Discover routers, services, middlewares, and TLS objects selectively per provider
- Powerful matchers to include/exclude resources by provider, name, and other rules
- Overrides to mutate discovered objects safely (rules, entrypoints, services, middlewares, etc.)
- Tunnels that:
  - Replace service server lists with tunnel addresses
  - Optionally create and attach HTTP ServersTransports with mTLS
- Automatic provider-suffix stripping on names and cross-references (e.g., `service@file` → `service`)
- Merge multiple providers into one consistent dynamic configuration (including ServersTransports)

## Architecture

- Core plugin entry: `traefikprovider/`
  - Entry point: `traefikprovider.go` (`Provider.Provide` sends merged JSON payloads periodically)
- Data model and config: `config/`
- Parsing pipeline: `internal/parsers/`
  - HTTP/TCP/UDP/TLS parse stages
  - Name cleanup and overrides: `internal/overrides/`
- Matchers: `internal/matchers/`
- Tunnels: `internal/tunnels/` (creates HTTP ServersTransports from mTLS, rewrites services)
- Merge: `internal/merge.go` (merges routers/services/middlewares/ServersTransports etc.)
- HTTP client: `internal/httpclient/` (fetches upstream raw JSON, parses into `dynamic.Configuration`)

## Installation

Use either GitHub plugin mode or local mode.

- Static config (plugin mode):
  - `.traefik.yml` already describes this provider plugin for Traefik Hub/Plugin Catalog
  - In Traefik static config, enable the plugin provider with this module name:
    - `github.com/zalbiraw/traefikprovider`

- Static config (local mode, recommended for development):
  - Use `experimental.localPlugins` and mount this repo inside Traefik:

```yaml
experimental:
  localPlugins:
    traefik:
      moduleName: github.com/zalbiraw/traefikprovider

providers:
  plugin:
    traefik:
      pollInterval: "5s"
      providers:
        # see Configuration Reference below
```

See `test/configs/traefik-main/traefik.yml` for a complete example.

## Quick Start (with the included test stack)

- Requirements: Docker + Docker Compose, Make
- Commands:
  - `make traefik-up` — brings up Traefik and the demo upstreams
  - `make traefik-down` — tears down
  - `make traefik-restart` — restarts the test stack

These use `test/docker-compose.yml` and supporting test configs under `test/configs/`.

## Configuration Reference

Root provider config (static):

- Path: Traefik static config (e.g., `traefik.yml`)
- Section:
  - `providers.plugin.traefik.pollInterval` string (Go duration, e.g. `"5s"`)
  - `providers.plugin.traefik.providers[]` array of upstream ProviderConfigs

ProviderConfig model (`config/config.go`):

- `name` string — descriptive name
- `matcher` string — provider-level matcher (e.g., `Provider('file')`)
- `connection`:
  - `host` string (required)
  - `port` int (required)
  - `path` string (required)
  - `timeout` string (Go duration)
  - `headers` map[string]string
  - `mTLS` (optional) for calling the upstream provider:
    - `caFile`, `certFile`, `keyFile` (paths)
- `http` `HTTPSection` (see below)
- `tcp` `TCPSection`
- `udp` `UDPSection`
- `tls` `TLSSection`
- `tunnels` []`TunnelConfig` (see Tunnels)

HTTPSection (`config/sections.go`):

- `discover` bool (default: true)
- `routers` `RoutersConfig`
- `middlewares` `MiddlewaresConfig`
- `services` `ServicesConfig`

TCPSection (`config/sections.go`):

- `discover` bool (default: true)
- `routers` `RoutersConfig`
- `middlewares` `MiddlewaresConfig`
- `services` `ServicesConfig`

UDPSection (`config/sections.go`):

- `discover` bool (default: true)
- `routers` `UDPRoutersConfig`
- `services` `UDPServicesConfig`

TLSSection (`config/sections.go`):

- `discover` bool (default: true)

RoutersConfig (`config/routers.go`):

- `discover` bool
- `discoverPriority` bool — keep discovered priorities when true; otherwise reset to 0
- `matcher` string — matcher to select routers (e.g., by name/provider)
- `stripServiceProvider` bool — if true, strip `@provider` from router.service
- `overrides` `RouterOverrides`
- `extraRoutes` []any — extra router definitions (raw)

RouterOverrides (`config/routers.go`):

- `name` string — rename matching routers
- `rules` []`OverrideRule`:
  - `value` string (rule)
  - `matcher` string
- `entrypoints` []`OverrideEntrypoint`:
  - `value` any (string or []string)
  - `matcher` string
- `services` []`OverrideService`:
  - `value` string (service name)
  - `matcher` string
- `middlewares` []`OverrideMiddleware`:
  - `value` any (string or []string)
  - `matcher` string

ServicesConfig, MiddlewaresConfig, and UDP configs follow the same pattern (discover, matcher, overrides, extra definitions). See files in `config/` for exact shapes.

### Tunnels (`config/config.go`, `internal/tunnels/tunnels.go`)

- `tunnels` array per provider config
- `TunnelConfig`:
  - `matcher` string — matches services to modify
  - `addresses` []string — replace service servers with these addresses
  - `mTLS` (optional) — when present, a HTTP ServersTransport is created:
    - `rootCAs` from `caFile`
    - `certificates` from `certFile` + `keyFile`
- For HTTP services only, when mTLS is provided:
  - The plugin creates a `ServersTransport` named `st-<hash of matcher>`
  - The matched services’ `loadBalancer.serversTransport` is set to that name
- Notes:
  - `ServersTransport` objects are stored under `http.serversTransports` in the resulting dynamic config.
  - This plugin does not set `serverName` automatically. If your upstream cert CN/SAN does not match the host in `addresses`, use a matching hostname or request an explicit `serverName` option.

## How Matching, Overrides, and Name Cleanup Work

- Names may include `@provider` suffixes (e.g., `serviceA@file`). The plugin strips `@provider` in keys and cross-references to make merging consistent across sources.
  - Code: `internal/overrides/names.go`
- Matching:
  - Provider-level matcher filters at the provider scope (e.g., only resources from a source)
  - Per-section matchers further filter resources
  - See `internal/matchers/` for matcher implementation
- Overrides:
  - Applied after parsing and name normalization
  - Routers: adjust rules, entrypoints, service, middlewares, and optional name rename
  - Services and Middlewares support similar override patterns (see `config/` and `internal/overrides/`)

## Merging Behavior

- The plugin polls all configured providers, builds a `*dynamic.Configuration` per provider, and merges them.
- Merge implementation: `internal/merge.go`
  - HTTP: merges `routers`, `services`, `middlewares`, and `serversTransports`
  - TCP/UDP/TLS: merges corresponding maps/arrays
- Later providers override earlier ones on identical keys.

## Example Static Configuration (local plugin mode)

```yaml
api:
  insecure: true
  dashboard: true

experimental:
  localPlugins:
    traefik:
      moduleName: github.com/zalbiraw/traefikprovider

entryPoints:
  web:
    address: ":80"

accessLog: {}

providers:
  plugin:
    traefik:
      pollInterval: "5s"
      providers:
        - name: provider1
          connection:
            host: traefik-provider1
            port: 8080
            path: /api/rawdata
            timeout: "10s"
          matcher: "Provider(`file`)"
          tunnels:
            - matcher: "NameRegexp(`.*`)"
              addresses:
                - "https://traefik-provider1:443"
              mTLS:
                caFile: "/etc/traefik/certs/ca.crt"
                certFile: "/etc/traefik/certs/client.crt"
                keyFile: "/etc/traefik/certs/client.key"
        - name: provider2
          connection:
            host: traefik-provider2
            port: 8080
            path: /api/rawdata
            timeout: "10s"
          matcher: "Provider(`file`)"
          tunnels:
            - matcher: "NameRegexp(`.*`)"
              addresses:
                - "https://traefik-provider2:443"
              mTLS:
                caFile: "/etc/traefik/certs/ca.crt"
                certFile: "/etc/traefik/certs/client.crt"
                keyFile: "/etc/traefik/certs/client.key"
```

## Use Cases

- Aggregate multiple upstream providers (e.g., File + CRDs + custom API) into one coherent config
- Enforce global patterns via matchers and overrides (e.g., entrypoints, rules)
- Front services via tunneling with mTLS from Traefik to upstream services without modifying upstream definitions
- Strip provider suffixes for fully merged names and references

## Troubleshooting

- Missing routes or 404s:
  - Ensure router rules and entrypoints match your request (e.g., Host header, PathPrefix)
  - Verify the merged config via Traefik dashboard or logs
- TLS/mTLS issues when tunneling:
  - Confirm `rootCAs` and client certs/keys paths are correct and mounted
  - If upstream certificate CN/SAN does not match tunnel hostname, use a matching hostname or request an explicit `serverName` option
- Name collisions:
  - Since provider postfixes are stripped, ensure you don’t unintentionally collide distinct services/routers with the same base name

## Development

- Code layout:
  - `traefikprovider.go` — plugin entry points and background polling
  - `config/` — user-facing configuration
  - `internal/parsers/` — parsing pipeline
  - `internal/overrides/` — name stripping + override application
  - `internal/tunnels/` — tunnels and ServersTransports creation/wiring
  - `internal/merge.go` — merging logic
  - `internal/httpclient/` — HTTP client to fetch upstream raw JSON
  - `test/` — docker-compose test stack with example upstreams
- Make targets:
  - `make traefik-up`, `make traefik-down`, `make traefik-restart`

## Security Notes

- Do not embed secrets directly into configs. Mount certs/keys via volumes.
- Be careful when enabling `discover` widely; scope with matchers.
- Only grant read access to upstream endpoints exposed via `connection`.

## License

See `LICENSE`.
