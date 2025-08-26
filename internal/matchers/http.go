// Package matchers provides utilities to filter dynamic configuration objects.
package matchers

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/rules"
)

// HTTPRouters filters HTTP routers based on `cfg.Matcher` and optional provider-level matcher.
func HTTPRouters(routers map[string]*dynamic.Router, cfg *config.RoutersConfig, providerMatcher string) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return routers
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, router := range routers {
		ctx := rules.Context{
			Name:        name,
			Provider:    extractProviderFromName(name),
			Entrypoints: router.EntryPoints,
			Service:     router.Service,
		}
		if prog.Match(ctx) {
			// Reset Priority to 0 when discovery of priority is disabled
			if cfg != nil && !cfg.DiscoverPriority {
				r := *router
				r.Priority = 0
				result[name] = &r
			} else {
				result[name] = router
			}
		}
	}
	return result
}

// HTTPServices filters HTTP services based on `cfg.Matcher` and optional provider-level matcher.
func HTTPServices(services map[string]*dynamic.Service, cfg *config.ServicesConfig, providerMatcher string) map[string]*dynamic.Service {
	result := make(map[string]*dynamic.Service)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return services
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, svc := range services {
		ctx := rules.Context{
			Name:     name,
			Provider: extractProviderFromName(name),
		}
		if prog.Match(ctx) {
			result[name] = svc
		}
	}
	return result
}

// HTTPMiddlewares filters HTTP middlewares based on `cfg.Matcher` and optional provider-level matcher.
func HTTPMiddlewares(middlewares map[string]*dynamic.Middleware, cfg *config.MiddlewaresConfig, providerMatcher string) map[string]*dynamic.Middleware {
	result := make(map[string]*dynamic.Middleware)
	combined := combineRules(providerMatcher, cfg.Matcher)
	if combined == "" {
		return middlewares
	}
	prog, err := compileRule(combined)
	if err != nil {
		return result
	}
	for name, mw := range middlewares {
		ctx := rules.Context{
			Name:     name,
			Provider: extractProviderFromName(name),
		}
		if prog.Match(ctx) {
			result[name] = mw
		}
	}
	return result
}
