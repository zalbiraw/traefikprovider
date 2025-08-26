// Package filters provides utilities to filter dynamic configuration objects.
package matchers

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/rules"
)

// tcpRouterMatchesFilter reports whether a TCP router matches the given filter.
func tcpRouterMatchesFilter(prog *rules.Program, name string, router *dynamic.TCPRouter) bool {
	ctx := rules.Context{
		Name:        name,
		Provider:    extractProviderFromName(name),
		Entrypoints: router.EntryPoints,
		Service:     router.Service,
	}
	return prog.Match(ctx)
}

// TCPRouters filters TCP routers based on `cfg.Matcher` and optional provider-level matcher.
func TCPRouters(routers map[string]*dynamic.TCPRouter, cfg *config.RoutersConfig, providerMatcher string) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return routers
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, router := range routers {
		if tcpRouterMatchesFilter(prog, name, router) {
			result[name] = router
		}
	}
	return result
}

// TCPServices filters TCP services based on `cfg.Matcher` and optional provider-level matcher.
func TCPServices(services map[string]*dynamic.TCPService, cfg *config.ServicesConfig, providerMatcher string) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return services
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, service := range services {
		ctx := rules.Context{
			Name:     name,
			Provider: extractProviderFromName(name),
		}
		if prog.Match(ctx) {
			result[name] = service
		}
	}
	return result
}

// TCPMiddlewares filters TCP middlewares based on `cfg.Matcher` and optional provider-level matcher.
func TCPMiddlewares(middlewares map[string]*dynamic.TCPMiddleware, cfg *config.MiddlewaresConfig, providerMatcher string) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return middlewares
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, middleware := range middlewares {
		ctx := rules.Context{
			Name:     name,
			Provider: extractProviderFromName(name),
		}
		if prog.Match(ctx) {
			result[name] = middleware
		}
	}
	return result
}
