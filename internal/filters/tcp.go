package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TCPRouters(routers map[string]*dynamic.TCPRouter, config *config.RoutersConfig) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	filter := config.Filter

	if filter.Name == "" && len(filter.Entrypoints) == 0 && filter.Rule == "" && filter.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.TCPRouter, *dynamic.TCPRouter](routers, filter.Name)
	for name, router := range filtered {
		if len(filter.Entrypoints) > 0 {
			if !routerEntrypointsMatch(router.EntryPoints, filter.Entrypoints) {
				continue
			}
		}
		if filter.Rule != "" {
			matched, err := regexMatch(filter.Rule, router.Rule)
			if err != nil || !matched {
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

func TCPServices(services map[string]*dynamic.TCPService, config *config.ServicesConfig) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	filter := config.Filter

	if filter.Name == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.TCPService, *dynamic.TCPService](services, filter.Name)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}

func TCPMiddlewares(middlewares map[string]*dynamic.TCPMiddleware, config *config.MiddlewaresConfig) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	filter := config.Filter

	if filter.Name == "" {
		return middlewares
	}

	filtered := filterMapByNameRegex[dynamic.TCPMiddleware, *dynamic.TCPMiddleware](middlewares, filter.Name)
	for name, middleware := range filtered {
		result[name] = middleware
	}
	return result
}
