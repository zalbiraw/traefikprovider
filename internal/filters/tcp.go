package filters

import (
	"encoding/json"
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TCPRouters(routers interface{}, config *config.RoutersConfig) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameRegex(routersMap, filters.Name)
	for name, routerMap := range filtered {
		router := &dynamic.TCPRouter{}
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
		if filters.Rule != "" {
			matched, err := regexMatch(filters.Rule, router.Rule)
			if err != nil || !matched {
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

func TCPServices(services interface{}, config *config.ServicesConfig) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameRegex(servicesMap, filters.Name)
	for name, serviceMap := range filtered {
		service := &dynamic.TCPService{}
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

func TCPMiddlewares(middlewares interface{}, config *config.MiddlewaresConfig) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	middlewaresMap, ok := middlewares.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameRegex(middlewaresMap, filters.Name)
	for name, middlewareMap := range filtered {
		middleware := &dynamic.TCPMiddleware{}
		b, err := json.Marshal(middlewareMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, middleware); err != nil {
			continue
		}
		result[name] = middleware
	}
	return result
}
