// Package overrides applies user-defined overrides to matched configs.
package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func applyOverrideTCP[T any](matched map[string]*dynamic.TCPRouter, value T, apply func(r *dynamic.TCPRouter, v T)) {
	for key, router := range matched {
		apply(router, value)
		matched[key] = router
	}
}

func handleOverrideTCP(
	matched map[string]*dynamic.TCPRouter,
	value interface{},
	applyArray func(r *dynamic.TCPRouter, arr []string),
	applyString func(r *dynamic.TCPRouter, s string),
) {
	switch v := value.(type) {
	case []string:
		applyOverrideTCP(matched, v, applyArray)
	case string:
		applyOverrideTCP(matched, v, applyString)
	}
}

// OverrideTCPRouters applies override rules to the given TCP routers map.
func OverrideTCPRouters(matched map[string]*dynamic.TCPRouter, overrides config.RouterOverrides) {
	for _, oep := range overrides.Entrypoints {
		handleOverrideTCP(matched, oep.Value,
			func(r *dynamic.TCPRouter, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.TCPRouter, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	for _, osvc := range overrides.Services {
		applyOverrideTCP(matched, osvc.Value, func(r *dynamic.TCPRouter, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}

	for _, omw := range overrides.Middlewares {
		handleOverrideTCP(matched, omw.Value,
			func(r *dynamic.TCPRouter, arr []string) { r.Middlewares = arr },
			func(r *dynamic.TCPRouter, s string) { r.Middlewares = append(r.Middlewares, s) },
		)
	}
}

// OverrideTCPServices applies overrides to matched TCP services.
// If a server override specifies a Tunnel, the matched services' servers are
// replaced with the tunnel addresses.
func OverrideTCPServices(matched map[string]*dynamic.TCPService, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
    // Server overrides
    for _, orule := range overrides.Servers {
        // If tunnel is specified and resolves to addresses, prefer it
        if orule.Tunnel != "" {
            addrs := resolveServerURLs(orule.Tunnel, tunnels)
            if len(addrs) > 0 {
                handleTCPServiceOverride(matched, orule.Matcher, addrs,
                    func(s *dynamic.TCPService, v []string) { s.LoadBalancer.Servers = buildTCPServers(v) },
                    func(s *dynamic.TCPService, v string) {},
                )
                continue
            }
        }

        handleTCPServiceOverride(matched, orule.Matcher, orule.Value,
            func(s *dynamic.TCPService, v []string) {
                servers := []dynamic.TCPServer{}
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

// buildTCPServers converts a list of URLs to dynamic.TCPServer slice.
func buildTCPServers(urls []string) []dynamic.TCPServer {
	if len(urls) == 0 {
		return nil
	}
	servers := make([]dynamic.TCPServer, 0, len(urls))
	for _, addr := range urls {
		servers = append(servers, dynamic.TCPServer{Address: addr})
	}
	return servers
}
