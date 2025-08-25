package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func HTTPRouters(routers map[string]*dynamic.Router, config *config.RoutersConfig) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	filters := config.Filter

	if filters.Name == "" && len(filters.Entrypoints) == 0 && filters.Rule == "" && filters.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.Router, *dynamic.Router](routers, filters.Name)
	for name, router := range filtered {
		if !config.DiscoverPriority {
			router.Priority = 0
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

func HTTPServices(services map[string]*dynamic.Service, config *config.ServicesConfig) map[string]*dynamic.Service {
	result := make(map[string]*dynamic.Service)
	filter := config.Filter

	if filter.Name == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.Service, *dynamic.Service](services, filter.Name)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}

func HTTPMiddlewares(middlewares map[string]*dynamic.Middleware, config *config.MiddlewaresConfig) map[string]*dynamic.Middleware {
	result := make(map[string]*dynamic.Middleware)
	filter := config.Filter

	if filter.Name == "" {
		return middlewares
	}

	filtered := filterMapByNameRegex[dynamic.Middleware, *dynamic.Middleware](middlewares, filter.Name)
	for name, middleware := range filtered {
		result[name] = middleware
	}
	return result
}
