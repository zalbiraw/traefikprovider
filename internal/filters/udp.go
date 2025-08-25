package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func UDPRouters(routers map[string]*dynamic.UDPRouter, cfg *config.UDPRoutersConfig, pf config.ProviderFilter) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" && len(filter.Entrypoints) == 0 && filter.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.UDPRouter, *dynamic.UDPRouter](routers, filter.Name, filter.Provider)
	for name, router := range filtered {
		if len(filter.Entrypoints) > 0 {
			if !routerEntrypointsMatch(router.EntryPoints, filter.Entrypoints) {
				continue
			}
		}
		if filter.Service != "" {
			matched, err := regexMatch(filter.Service, router.Service)
			if err != nil || !matched {
				continue
			}
		}
		result[name] = router
	}
	return result
}

func UDPServices(services map[string]*dynamic.UDPService, cfg *config.UDPServicesConfig, pf config.ProviderFilter) map[string]*dynamic.UDPService {
	result := make(map[string]*dynamic.UDPService)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.UDPService, *dynamic.UDPService](services, filter.Name, filter.Provider)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}
