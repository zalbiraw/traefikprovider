package overrides

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/matchers"
)

func applyRouterOverride[T any](matched map[string]*dynamic.Router, matcher string, value T, apply func(r *dynamic.Router, v T)) {
	rc := &config.RoutersConfig{Matcher: matcher}
	for key, router := range matchers.HTTPRouters(matched, rc, "") {
		apply(router, value)
		matched[key] = router
	}
}

func handleRouterOverride(
	matched map[string]*dynamic.Router,
	matcher string,
	value interface{},
	applyArray func(r *dynamic.Router, arr []string),
	applyString func(r *dynamic.Router, s string),
) {
	switch v := value.(type) {
	case []string:
		applyRouterOverride(matched, matcher, v, applyArray)
	case string:
		applyRouterOverride(matched, matcher, v, applyString)
	}
}

func applyServiceOverride[T any](matched map[string]*dynamic.Service, matcher string, value T, apply func(r *dynamic.Service, v T)) {
	rc := &config.ServicesConfig{Matcher: matcher}
	for key, service := range matchers.HTTPServices(matched, rc, "") {
		apply(service, value)
		matched[key] = service
	}
}

func handleServiceOverride(
	matched map[string]*dynamic.Service,
	matcher string,
	value interface{},
	applyArray func(r *dynamic.Service, arr []string),
	applyString func(r *dynamic.Service, s string),
) {
	switch v := value.(type) {
	case []string:
		applyServiceOverride(matched, matcher, v, applyArray)
	case string:
		applyServiceOverride(matched, matcher, v, applyString)
	}
}

func applyTCPServiceOverride[T any](matched map[string]*dynamic.TCPService, matcher string, value T, apply func(r *dynamic.TCPService, v T)) {
	rc := &config.ServicesConfig{Matcher: matcher}
	for key, service := range matchers.TCPServices(matched, rc, "") {
		apply(service, value)
		matched[key] = service
	}
}

func handleTCPServiceOverride(
	matched map[string]*dynamic.TCPService,
	matcher string,
	value interface{},
	applyArray func(r *dynamic.TCPService, arr []string),
	applyString func(r *dynamic.TCPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyTCPServiceOverride(matched, matcher, v, applyArray)
	case string:
		applyTCPServiceOverride(matched, matcher, v, applyString)
	}
}

func applyUDPServiceOverride[T any](matched map[string]*dynamic.UDPService, matcher string, value T, apply func(r *dynamic.UDPService, v T)) {
	rc := &config.UDPServicesConfig{Matcher: matcher}
	for key, service := range matchers.UDPServices(matched, rc, "") {
		apply(service, value)
		matched[key] = service
	}
}

func handleUDPServiceOverride(
	matched map[string]*dynamic.UDPService,
	matcher string,
	value interface{},
	applyArray func(r *dynamic.UDPService, arr []string),
	applyString func(r *dynamic.UDPService, s string),
) {
	switch v := value.(type) {
	case []string:
		applyUDPServiceOverride(matched, matcher, v, applyArray)
	case string:
		applyUDPServiceOverride(matched, matcher, v, applyString)
	}
}

func resolveServerURLs(tunnelName string, tunnels []config.TunnelConfig) []string {
	if tunnelName != "" {
		for _, t := range tunnels {
			if t.Name == tunnelName {
				return t.Addresses
			}
		}
	}
	return []string{}
}

// ApplyTunnels applies tunnels after other overrides. If a server override rule
// specifies a Tunnel, the matched services' servers are replaced with the tunnel addresses.
func ApplyTunnels(matched interface{}, overrides config.ServiceOverrides, tunnels []config.TunnelConfig) {
	for _, orule := range overrides.Servers {
		if orule.Tunnel == "" {
			continue
		}
		addrs := resolveServerURLs(orule.Tunnel, tunnels)
		if len(addrs) == 0 {
			continue
		}

		switch m := matched.(type) {
		case map[string]*dynamic.Service:
			handleServiceOverride(m, orule.Matcher, addrs,
				func(s *dynamic.Service, v []string) { s.LoadBalancer.Servers = buildServers(v) },
				func(s *dynamic.Service, v string) {},
			)
		case map[string]*dynamic.TCPService:
			handleTCPServiceOverride(m, orule.Matcher, addrs,
				func(s *dynamic.TCPService, v []string) { s.LoadBalancer.Servers = buildTCPServers(v) },
				func(s *dynamic.TCPService, v string) {},
			)
		case map[string]*dynamic.UDPService:
			handleUDPServiceOverride(m, orule.Matcher, addrs,
				func(s *dynamic.UDPService, v []string) { s.LoadBalancer.Servers = buildUDPServers(v) },
				func(s *dynamic.UDPService, v string) {},
			)
		}
	}
}
