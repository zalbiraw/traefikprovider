// Package overrides applies user-defined overrides to filtered configs.
package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func applyOverrideTCP[T any](filtered map[string]*dynamic.TCPRouter, value T, apply func(r *dynamic.TCPRouter, v T)) {
	for key, router := range filtered {
		apply(router, value)
		filtered[key] = router
	}
}

func handleOverrideTCP(
	filtered map[string]*dynamic.TCPRouter,
	value interface{},
	applyArray func(r *dynamic.TCPRouter, arr []string),
	applyString func(r *dynamic.TCPRouter, s string),
) {
	switch v := value.(type) {
	case []string:
		applyOverrideTCP(filtered, v, applyArray)
	case string:
		applyOverrideTCP(filtered, v, applyString)
	}
}

// OverrideTCPRouters applies override rules to the given TCP routers map.
func OverrideTCPRouters(filtered map[string]*dynamic.TCPRouter, overrides config.RouterOverrides) {
	for _, oep := range overrides.Entrypoints {
		handleOverrideTCP(filtered, oep.Value,
			func(r *dynamic.TCPRouter, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.TCPRouter, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	for _, osvc := range overrides.Services {
		applyOverrideTCP(filtered, osvc.Value, func(r *dynamic.TCPRouter, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}

	for _, omw := range overrides.Middlewares {
		handleOverrideTCP(filtered, omw.Value,
			func(r *dynamic.TCPRouter, arr []string) { r.Middlewares = arr },
			func(r *dynamic.TCPRouter, s string) { r.Middlewares = append(r.Middlewares, s) },
		)
	}
}

// OverrideTCPServices applies overrides to filtered TCP services.
func OverrideTCPServices(filtered map[string]*dynamic.TCPService, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
	// Server overrides
	for _, orule := range overrides.Servers {
		handleTCPServiceOverride(filtered, orule.Filter, orule.Value,
			func(s *dynamic.TCPService, v []string) {
				servers := []dynamic.TCPServer{}

				// If tunnel is specified, use tunnel addresses instead of v
				if orule.Tunnel != "" {
					for _, tunnel := range tunnels {
						if tunnel.Name == orule.Tunnel {
							v = tunnel.Addresses
						}
					}
				}

				// Use the provided addresses (or tunnel addresses)
				for _, addr := range v {
					server := dynamic.TCPServer{Address: addr}
					servers = append(servers, server)
				}

				s.LoadBalancer.Servers = servers
			},
			func(s *dynamic.TCPService, v string) {
				s.LoadBalancer.Servers = append(s.LoadBalancer.Servers, dynamic.TCPServer{Address: v})
			},
		)
	}
}
