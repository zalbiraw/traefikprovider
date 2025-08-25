package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func UDPRouters(routers map[string]*dynamic.UDPRouter, config *config.UDPRoutersConfig) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
	filters := config.Filters
	filtered := filterMapByNameRegex[dynamic.UDPRouter, *dynamic.UDPRouter](routers, filters.Name)
	for name, router := range filtered {
		if len(filters.Entrypoints) > 0 {
			if !routerEntrypointsMatch(router.EntryPoints, filters.Entrypoints) {
				continue
			}
		}
		if filters.Service != "" {
			matched, err := regexMatch(filters.Service, router.Service)
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
	filters := config.Filters
	filtered := filterMapByNameRegex[dynamic.UDPService, *dynamic.UDPService](services, filters.Name)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}
