// Package overrides applies user-defined overrides to matched configs.
package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

// OverrideUDPRouters applies override rules to the given UDP routers map.
func applyOverrideUDP[T any](matched map[string]*dynamic.UDPRouter, value T, apply func(r *dynamic.UDPRouter, v T)) {
	for key, router := range matched {
		apply(router, value)
		matched[key] = router
	}
}

func handleOverrideUDP(
	matched map[string]*dynamic.UDPRouter,
	value interface{},
	applyArray func(r *dynamic.UDPRouter, arr []string),
	applyString func(r *dynamic.UDPRouter, s string),
) {
	switch v := value.(type) {
	case []string:
		applyOverrideUDP(matched, v, applyArray)
	case string:
		applyOverrideUDP(matched, v, applyString)
	}
}

// OverrideUDPRouters applies overrides to the provided UDP routers map.
func OverrideUDPRouters(matched map[string]*dynamic.UDPRouter, overrides config.UDPOverrides) {
	for _, oep := range overrides.Entrypoints {
		handleOverrideUDP(matched, oep.Value,
			func(r *dynamic.UDPRouter, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.UDPRouter, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	for _, osvc := range overrides.Services {
		applyOverrideUDP(matched, osvc.Value, func(r *dynamic.UDPRouter, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}
}

// OverrideUDPServices applies overrides to matched UDP services.
// If a server override specifies a Tunnel, the matched services' servers are
// replaced with the tunnel addresses.
func OverrideUDPServices(matched map[string]*dynamic.UDPService, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
    // Server overrides
    for _, orule := range overrides.Servers {
        // If tunnel is specified and resolves to addresses, prefer it
        if orule.Tunnel != "" {
            addrs := resolveServerURLs(orule.Tunnel, tunnels)
            if len(addrs) > 0 {
                handleUDPServiceOverride(matched, orule.Matcher, addrs,
                    func(s *dynamic.UDPService, v []string) { s.LoadBalancer.Servers = buildUDPServers(v) },
                    func(s *dynamic.UDPService, v string) {},
                )
                continue
            }
        }

        handleUDPServiceOverride(matched, orule.Matcher, orule.Value,
            func(s *dynamic.UDPService, v []string) {
                servers := []dynamic.UDPServer{}
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

// buildUDPServers converts a list of URLs to dynamic.UDPServer slice.
func buildUDPServers(urls []string) []dynamic.UDPServer {
	if len(urls) == 0 {
		return nil
	}
	servers := make([]dynamic.UDPServer, 0, len(urls))
	for _, addr := range urls {
		servers = append(servers, dynamic.UDPServer{Address: addr})
	}
	return servers
}
