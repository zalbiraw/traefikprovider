package filters

import (
	"encoding/json"
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func UDPRouters(routers interface{}, config *config.UDPRoutersConfig) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(routersMap, filters.Name, filters.NameRegex)
	for name, routerMap := range filtered {
		router := &dynamic.UDPRouter{}
		b, err := json.Marshal(routerMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, router); err != nil {
			continue
		}
		if len(filters.Entrypoints) > 0 {
			if !routerEntrypointsMatch(router.EntryPoints, filters.Entrypoints) {
				continue
			}
		}
		if filters.Service != "" && router.Service != filters.Service {
			continue
		}
		if filters.ServiceRegex != "" {
			matched, err := regexMatch(filters.ServiceRegex, router.Service)
			if err != nil || !matched {
				continue
			}
		}
		result[name] = router
	}
	return result
}

func UDPServices(services interface{}, config *config.UDPServicesConfig) map[string]*dynamic.UDPService {
	result := make(map[string]*dynamic.UDPService)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(servicesMap, filters.Name, filters.NameRegex)
	for name, serviceMap := range filtered {
		service := &dynamic.UDPService{}
		b, err := json.Marshal(serviceMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, service); err != nil {
			continue
		}
		result[name] = service
	}
	return result
}
