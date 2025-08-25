package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TCPRouters(routers map[string]*dynamic.TCPRouter, config *config.RoutersConfig) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	filters := config.Filters
	filtered := filterMapByNameRegex[dynamic.TCPRouter, *dynamic.TCPRouter](routers, filters.Name)
	for name, router := range filtered {
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

func TCPServices(services map[string]*dynamic.TCPService, config *config.ServicesConfig) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	filters := config.Filters
	filtered := filterMapByNameRegex[dynamic.TCPService, *dynamic.TCPService](services, filters.Name)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}

func TCPMiddlewares(middlewares map[string]*dynamic.TCPMiddleware, config *config.MiddlewaresConfig) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	filters := config.Filters
	filtered := filterMapByNameRegex[dynamic.TCPMiddleware, *dynamic.TCPMiddleware](middlewares, filters.Name)
	for name, middleware := range filtered {
		result[name] = middleware
	}
	return result
}
