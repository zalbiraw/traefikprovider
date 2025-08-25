package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

// OverrideUDPRouters will later apply overrides to filtered UDP routers.

func applyOverrideUDP[T any](filtered map[string]*dynamic.UDPRouter, value T, apply func(r *dynamic.UDPRouter, v T)) {
	for key, router := range filtered {
		apply(router, value)
		filtered[key] = router
	}
}

func handleOverrideUDP(
	filtered map[string]*dynamic.UDPRouter,
	value interface{},
	applyArray func(r *dynamic.UDPRouter, arr []string),
	applyString func(r *dynamic.UDPRouter, s string),
) {
	switch v := value.(type) {
	case []string:
		applyOverrideUDP(filtered, v, applyArray)
	case string:
		applyOverrideUDP(filtered, v, applyString)
	}
}

func OverrideUDPRouters(filtered map[string]*dynamic.UDPRouter, overrides config.UDPOverrides) {
	for _, oep := range overrides.Entrypoints {
		handleOverrideUDP(filtered, oep.Value,
			func(r *dynamic.UDPRouter, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.UDPRouter, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	for _, osvc := range overrides.Services {
		applyOverrideUDP(filtered, osvc.Value, func(r *dynamic.UDPRouter, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}
}

// OverrideUDPServices applies overrides to filtered UDP services.
func OverrideUDPServices(filtered map[string]*dynamic.UDPService, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
	// Server overrides
	for _, orule := range overrides.Servers {
		handleUDPServiceOverride(filtered, orule.Filter, orule.Value,
			func(s *dynamic.UDPService, v []string) {
				servers := []dynamic.UDPServer{}

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
					server := dynamic.UDPServer{Address: addr}
					servers = append(servers, server)
				}

				s.LoadBalancer.Servers = servers
			},
			func(s *dynamic.UDPService, v string) {
				s.LoadBalancer.Servers = append(s.LoadBalancer.Servers, dynamic.UDPServer{Address: v})
			},
		)
	}
}
