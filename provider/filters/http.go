package filters

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func HTTPRouters(routers interface{}, config *config.RoutersConfig) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	routersMap, ok := routers.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(routersMap, filters.Name, filters.NameRegex)
	for name, routerMap := range filtered {
		router := &dynamic.Router{}
		if config.DiscoverPriority {
			router.Priority = extractRouterPriority(routerMap, name)
		}
		if err := unmarshalRouter(routerMap, router); err != nil {
			continue
		}
		if len(filters.Entrypoints) > 0 {
			if !routerEntrypointsMatch(router.EntryPoints, filters.Entrypoints) {
				continue
			}
		}
		if filters.Rule != "" && router.Rule != filters.Rule {
			continue
		}
		if filters.RuleRegex != "" {
			matched, err := regexMatch(filters.RuleRegex, router.Rule)
			if err != nil || !matched {
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

func HTTPServices(services interface{}, config *config.ServicesConfig) map[string]*dynamic.Service {
	result := make(map[string]*dynamic.Service)
	servicesMap, ok := services.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(servicesMap, filters.Name, filters.NameRegex)
	for name, serviceMap := range filtered {
		service := &dynamic.Service{}
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

func HTTPMiddlewares(middlewares interface{}, config *config.MiddlewaresConfig) map[string]*dynamic.Middleware {
	result := make(map[string]*dynamic.Middleware)
	middlewaresMap, ok := middlewares.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(middlewaresMap, filters.Name, filters.NameRegex)
	for name, middlewareMap := range filtered {
		middleware := &dynamic.Middleware{}
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

func HTTPServerTransports(serverTransports interface{}, config *config.ServerTransportsConfig) map[string]*dynamic.ServersTransport {
	result := make(map[string]*dynamic.ServersTransport)
	stMap, ok := serverTransports.(map[string]interface{})
	if !ok {
		return result
	}
	filters := config.Filters
	filtered := filterMapByNameAndRegex(stMap, filters.Name, filters.NameRegex)
	for name, stItemMap := range filtered {
		st := &dynamic.ServersTransport{}
		b, err := json.Marshal(stItemMap)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, st); err != nil {
			continue
		}
		result[name] = st
	}
	return result
}
