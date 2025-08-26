// Package matchers provides utilities to filter dynamic configuration objects
// (HTTP/TCP/UDP/TLS) using a rules-based matcher DSL.
package matchers

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/rules"
)

// UDPRouters filters UDP routers based on `cfg.Matcher` and optional provider-level matcher.
func UDPRouters(routers map[string]*dynamic.UDPRouter, cfg *config.UDPRoutersConfig, providerMatcher string) map[string]*dynamic.UDPRouter {
	result := make(map[string]*dynamic.UDPRouter)
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
			result[name] = router
		}
	}
	return result
}

// UDPServices filters UDP services based on `cfg.Matcher` and optional provider-level matcher.
func UDPServices(services map[string]*dynamic.UDPService, cfg *config.UDPServicesConfig, providerMatcher string) map[string]*dynamic.UDPService {
	result := make(map[string]*dynamic.UDPService)
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
