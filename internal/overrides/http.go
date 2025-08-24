package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func OverrideHTTPRouters(filtered map[string]*dynamic.Router, overrides config.RouterOverrides) {
	// Rule overrides
	for _, orule := range overrides.Rules {
		applyRouterOverride(filtered, orule.Filters, orule.Value, func(r *dynamic.Router, v string) {
			if strings.Contains(v, "$1") {
				r.Rule = strings.ReplaceAll(v, "$1", r.Rule)
			} else {
				r.Rule = v
			}
		})
	}

	// Entrypoint overrides
	for _, oep := range overrides.Entrypoints {
		handleRouterOverride(filtered, oep.Filters, oep.Value,
			func(r *dynamic.Router, arr []string) { r.EntryPoints = arr },
			func(r *dynamic.Router, s string) { r.EntryPoints = append(r.EntryPoints, s) },
		)
	}

	// Service overrides
	for _, osvc := range overrides.Services {
		applyRouterOverride(filtered, osvc.Filters, osvc.Value, func(r *dynamic.Router, v string) {
			if strings.Contains(v, "$1") {
				r.Service = strings.ReplaceAll(v, "$1", r.Service)
			} else {
				r.Service = v
			}
		})
	}

	// Middlewares overrides
	for _, omw := range overrides.Middlewares {
		handleRouterOverride(filtered, omw.Filters, omw.Value,
			func(r *dynamic.Router, arr []string) { r.Middlewares = arr },
			func(r *dynamic.Router, s string) { r.Middlewares = append(r.Middlewares, s) },
		)
	}
}

// OverrideHTTPServices applies overrides to filtered HTTP services.
func OverrideHTTPServices(filtered map[string]*dynamic.Service, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
	// Server overrides
	for _, orule := range overrides.Servers {
		handleServiceOverride(filtered, orule.Filters, orule.Value,
			func(s *dynamic.Service, v []string) {
				servers := []dynamic.Server{}

				// If tunnel is specified, use tunnel addresses instead of v
				if orule.Tunnel != "" {
					for _, tunnel := range tunnels {
						if tunnel.Name == orule.Tunnel {
							v = tunnel.Addresses
						}
					}
				}
				for _, addr := range v {
					server := dynamic.Server{URL: addr}
					servers = append(servers, server)
				}

				s.LoadBalancer.Servers = servers
			},
			func(s *dynamic.Service, v string) {
				s.LoadBalancer.Servers = append(s.LoadBalancer.Servers, dynamic.Server{URL: v})
			},
		)
	}
	// Healthcheck overrides
	for _, ohc := range overrides.Healthchecks {
		applyServiceOverride(filtered, ohc.Filters, ohc, func(s *dynamic.Service, hc config.OverrideHealthcheck) {
			if s.LoadBalancer != nil && s.LoadBalancer.HealthCheck != nil {
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
		})
	}
}
