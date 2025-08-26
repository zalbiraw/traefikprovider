# Traefik Provider Preview Environment

This directory contains Docker resources for testing the traefikprovider with multiple Traefik instances.

## Quick Start

```bash
cd cmd/preview
docker-compose up -d
```

## Architecture

- **traefik-main** (port 8080): Main instance with HTTP/TCP/UDP/TLS configs
- **traefik-provider1** (port 8081): Secondary instance for provider merging
- **traefik-provider2** (port 8082): Third instance for provider merging
- **Backend services**: whoami, nginx, tcp-echo, udp-echo, redis

## API Endpoints

- Main Traefik: http://localhost:8080/api/rawdata
- Provider1: http://localhost:8081/api/rawdata  
- Provider2: http://localhost:8082/api/rawdata

## Configuration

Each Traefik instance has its configuration in `configs/`:
- `traefik-main/dynamic.yml`: Comprehensive HTTP/TCP/UDP/TLS setup
- `traefik-provider1/dynamic.yml`: Provider1 specific configs
- `traefik-provider2/dynamic.yml`: Provider2 specific configs

## Testing

Use the main.go in this directory to test the provider against these instances:

```bash
go run main.go
```

## Cleanup

```bash
docker-compose down -v
```
