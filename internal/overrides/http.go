// Package overrides applies user-defined overrides to filtered configs.
package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

// OverrideHTTPRouters applies override rules to the given HTTP routers map.
func OverrideHTTPRouters(filtered map[string]*dynamic.Router, overrides config.RouterOverrides) {
	// Rule overrides
	for _, orule := range overrides.Rules {
		applyRouterOverride(filtered, orule.Filter, orule.Value, func(r *dynamic.Router, v string) {
			if strings.Contains(v, "$1") {
				r.Rule = strings.ReplaceAll(v, "$1", r.Rule)
			} else {
				r.Rule = v
			}
		})
	}

	// Entrypoint overrides
	for _, oep := range overrides.Entrypoints {
		handleRouterOverride(filtered, oep.Filter, oep.Value,
			func(r *dynamic.Router, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.Router, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	// Service overrides
	for _, osvc := range overrides.Services {
		applyRouterOverride(filtered, osvc.Filter, osvc.Value, func(r *dynamic.Router, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}

	// Middlewares overrides
	for _, omw := range overrides.Middlewares {
		handleRouterOverride(filtered, omw.Filter, omw.Value,
			func(r *dynamic.Router, arr []string) { r.Middlewares = arr },
			func(r *dynamic.Router, s string) { r.Middlewares = append(r.Middlewares, s) },
		)
	}
}

// OverrideHTTPServices applies overrides to filtered HTTP services.
func OverrideHTTPServices(filtered map[string]*dynamic.Service, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
	// Server overrides
	for _, orule := range overrides.Servers {
		handleServiceOverride(filtered, orule.Filter, orule.Value,
			func(s *dynamic.Service, v []string) {
				urls := resolveServerURLs(orule.Tunnel, tunnels, v)
				s.LoadBalancer.Servers = buildServers(urls)
			},
			func(s *dynamic.Service, v string) {
				s.LoadBalancer.Servers = append(s.LoadBalancer.Servers, dynamic.Server{URL: v})
			},
		)
	}
	// Healthcheck overrides
	for _, ohc := range overrides.Healthchecks {
		applyServiceOverride(filtered, ohc.Filter, ohc, func(s *dynamic.Service, hc config.OverrideHealthcheck) {
			applyHealthcheck(s, hc)
		})
	}
}

// resolveServerURLs returns the list of server URLs to use, taking tunnels into account.
func resolveServerURLs(tunnelName string, tunnels []config.TunnelConfig, defaults []string) []string {
	if tunnelName == "" {
		return defaults
	}
	for _, t := range tunnels {
		if t.Name == tunnelName {
			return t.Addresses
		}
	}
	return defaults
}

// buildServers converts a list of URLs to dynamic.Server slice.
func buildServers(urls []string) []dynamic.Server {
	if len(urls) == 0 {
		return nil
	}
	servers := make([]dynamic.Server, 0, len(urls))
	for _, addr := range urls {
		servers = append(servers, dynamic.Server{URL: addr})
	}
	return servers
}

// applyHealthcheck applies non-empty healthcheck fields.
func applyHealthcheck(s *dynamic.Service, hc config.OverrideHealthcheck) {
	if s.LoadBalancer == nil || s.LoadBalancer.HealthCheck == nil {
		return
	}
	if hc.Path != "" {
		s.LoadBalancer.HealthCheck.Path = hc.Path
	}
	if hc.Interval != "" {
		s.LoadBalancer.HealthCheck.Interval = hc.Interval
	}
	if hc.Timeout != "" {
		s.LoadBalancer.HealthCheck.Timeout = hc.Timeout
	}
}
