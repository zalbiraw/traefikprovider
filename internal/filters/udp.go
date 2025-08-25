package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func UDPRouters(routers map[string]*dynamic.UDPRouter, config *config.UDPRoutersConfig) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
	filter := config.Filter

	if filter.Name == "" && len(filter.Entrypoints) == 0 && filter.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.UDPRouter, *dynamic.UDPRouter](routers, filter.Name)
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

func UDPServices(services map[string]*dynamic.UDPService, config *config.UDPServicesConfig) map[string]*dynamic.UDPService {
	result := make(map[string]*dynamic.UDPService)
	filter := config.Filter

	if filter.Name == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.UDPService, *dynamic.UDPService](services, filter.Name)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}
