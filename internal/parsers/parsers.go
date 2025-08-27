// Package parsers provides functions to parse raw provider data into Traefik dynamic config.
package parsers

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/matchers"
	"github.com/zalbiraw/traefikprovider/internal/overrides"
	"github.com/zalbiraw/traefikprovider/internal/tunnels"
)

// convertToTyped converts a loosely-typed map to a map of typed pointers.
//
//nolint:nestif // deeply nested due to JSON shape handling
func convertToTyped[T any](data interface{}) map[string]*T {
	result := make(map[string]*T)
	if dataMap, ok := data.(map[string]interface{}); ok {
		for name, itemData := range dataMap {
			item := new(T)
			if itemMap, ok := itemData.(map[string]interface{}); ok {
				b, err := json.Marshal(itemMap)
				if err != nil {
					continue
				}
				if err := json.Unmarshal(b, item); err != nil {
					continue
				}
				result[name] = item
			}
		}
	}
	return result
}

func ensureHTTPDefaults(pc *config.HTTPSection) {
	if pc.Routers == nil {
		pc.Routers = &config.RoutersConfig{Discover: true}
	}
	if pc.Services == nil {
		pc.Services = &config.ServicesConfig{Discover: true}
	}
	if pc.Middlewares == nil {
		pc.Middlewares = &config.MiddlewaresConfig{Discover: true}
	}
}

// ParseHTTPConfig fills httpConfig from raw data according to providerConfig and tunnels.
func ParseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, providerConfig *config.HTTPSection, providerMatcher string, tns []config.TunnelConfig) {
	ensureHTTPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processHTTPRouters(raw, httpConfig, providerConfig, providerMatcher)
	}
	if providerConfig.Services.Discover {
		processHTTPServices(raw, httpConfig, providerConfig, providerMatcher, tns)
	}
	if providerConfig.Middlewares.Discover {
		processHTTPMiddlewares(raw, httpConfig, providerConfig, providerMatcher)
	}
}

func processHTTPRouters(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, providerMatcher string) {
	if routers, ok := raw["routers"]; ok {
		typedRouters := convertToTyped[dynamic.Router](routers)
		httpConfig.Routers = matchers.HTTPRouters(typedRouters, pc.Routers, providerMatcher)
	}
	for _, extra := range pc.Routers.ExtraRoutes {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var router dynamic.Router
		if err := json.Unmarshal(b, &router); err != nil {
			continue
		}
		if routerName, ok := extra.(map[string]interface{})["name"].(string); ok {
			httpConfig.Routers[routerName] = &router
		}
	}
	overrides.StripProvidersHTTP(httpConfig)
	overrides.OverrideHTTPRouters(httpConfig.Routers, pc.Routers.Overrides)

	if pc.Routers != nil && !pc.Routers.DiscoverPriority {
		for _, r := range httpConfig.Routers {
			r.Priority = 0
		}
	}
}

func processHTTPServices(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, providerMatcher string, tns []config.TunnelConfig) {
	if services, ok := raw["services"]; ok {
		typedServices := convertToTyped[dynamic.Service](services)
		httpConfig.Services = matchers.HTTPServices(typedServices, pc.Services, providerMatcher)
	}
	for _, extra := range pc.Services.ExtraServices {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var service dynamic.Service
		if err := json.Unmarshal(b, &service); err != nil {
			continue
		}
		if serviceName, ok := extra.(map[string]interface{})["name"].(string); ok {
			httpConfig.Services[serviceName] = &service
		}
	}
	overrides.StripProvidersHTTP(httpConfig)
	overrides.OverrideHTTPServices(httpConfig.Services, pc.Services.Overrides, tns)

	// Apply tunnels by matcher after overrides
	tunnels.ApplyHTTPTunnels(httpConfig, providerMatcher, tns)
}

func processHTTPMiddlewares(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, providerMatcher string) {
	if middlewares, ok := raw["middlewares"]; ok {
		typedMiddlewares := convertToTyped[dynamic.Middleware](middlewares)
		httpConfig.Middlewares = matchers.HTTPMiddlewares(typedMiddlewares, pc.Middlewares, providerMatcher)
	}
	for _, extra := range pc.Middlewares.ExtraMiddlewares {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var middleware dynamic.Middleware
		if err := json.Unmarshal(b, &middleware); err != nil {
			continue
		}
		if middlewareName, ok := extra.(map[string]interface{})["name"].(string); ok {
			httpConfig.Middlewares[middlewareName] = &middleware
		}
	}
	overrides.StripProvidersHTTP(httpConfig)
}

func ensureTCPDefaults(pc *config.TCPSection) {
	if pc.Routers == nil {
		pc.Routers = &config.RoutersConfig{Discover: true}
	}
	if pc.Services == nil {
		pc.Services = &config.ServicesConfig{Discover: true}
	}
	if pc.Middlewares == nil {
		pc.Middlewares = &config.MiddlewaresConfig{Discover: true}
	}
}

// ParseTCPConfig fills tcpConfig from raw data according to providerConfig and tunnels.
func ParseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, providerConfig *config.TCPSection, providerMatcher string, tns []config.TunnelConfig) {
	ensureTCPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processTCPRouters(raw, tcpConfig, providerConfig, providerMatcher)
	}
	if providerConfig.Services.Discover {
		processTCPServices(raw, tcpConfig, providerConfig, providerMatcher, tns)
	}
	if providerConfig.Middlewares.Discover {
		processTCPMiddlewares(raw, tcpConfig, providerConfig, providerMatcher)
	}
}

func processTCPRouters(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, providerMatcher string) {
	if routers, ok := raw["tcpRouters"]; ok {
		typedRouters := convertToTyped[dynamic.TCPRouter](routers)
		tcpConfig.Routers = matchers.TCPRouters(typedRouters, pc.Routers, providerMatcher)
	}
	for _, extra := range pc.Routers.ExtraRoutes {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var router dynamic.TCPRouter
		if err := json.Unmarshal(b, &router); err != nil {
			continue
		}
		if routerName, ok := extra.(map[string]interface{})["name"].(string); ok {
			tcpConfig.Routers[routerName] = &router
		}
	}
	overrides.StripProvidersTCP(tcpConfig)
	overrides.OverrideTCPRouters(tcpConfig.Routers, pc.Routers.Overrides)

	if pc.Routers != nil && !pc.Routers.DiscoverPriority {
		for _, r := range tcpConfig.Routers {
			r.Priority = 0
		}
	}
}

func processTCPServices(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, providerMatcher string, tns []config.TunnelConfig) {
	if services, ok := raw["tcpServices"]; ok {
		typedServices := convertToTyped[dynamic.TCPService](services)
		tcpConfig.Services = matchers.TCPServices(typedServices, pc.Services, providerMatcher)
	}
	for _, extra := range pc.Services.ExtraServices {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var service dynamic.TCPService
		if err := json.Unmarshal(b, &service); err != nil {
			continue
		}
		if serviceName, ok := extra.(map[string]interface{})["name"].(string); ok {
			tcpConfig.Services[serviceName] = &service
		}
	}
	overrides.StripProvidersTCP(tcpConfig)
	overrides.OverrideTCPServices(tcpConfig.Services, pc.Services.Overrides, tns)

	// Apply tunnels by matcher after overrides
	tunnels.ApplyTCPTunnels(tcpConfig, providerMatcher, tns)
}

func processTCPMiddlewares(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, providerMatcher string) {
	if middlewares, ok := raw["tcpMiddlewares"]; ok {
		typedMiddlewares := convertToTyped[dynamic.TCPMiddleware](middlewares)
		tcpConfig.Middlewares = matchers.TCPMiddlewares(typedMiddlewares, pc.Middlewares, providerMatcher)
	}
	for _, extra := range pc.Middlewares.ExtraMiddlewares {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var middleware dynamic.TCPMiddleware
		if err := json.Unmarshal(b, &middleware); err != nil {
			continue
		}
		if middlewareName, ok := extra.(map[string]interface{})["name"].(string); ok {
			tcpConfig.Middlewares[middlewareName] = &middleware
		}
	}
	overrides.StripProvidersTCP(tcpConfig)
}

func ensureUDPDefaults(pc *config.UDPSection) {
	if pc.Routers == nil {
		pc.Routers = &config.UDPRoutersConfig{Discover: true}
	}
	if pc.Services == nil {
		pc.Services = &config.UDPServicesConfig{Discover: true}
	}
}

// ParseUDPConfig fills udpConfig from raw data according to providerConfig.
func ParseUDPConfig(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, providerConfig *config.UDPSection, providerMatcher string) {
	ensureUDPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processUDPRouters(raw, udpConfig, providerConfig, providerMatcher)
	}
	if providerConfig.Services.Discover {
		processUDPServices(raw, udpConfig, providerConfig, providerMatcher)
	}
}

func processUDPRouters(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, pc *config.UDPSection, providerMatcher string) {
	if routers, ok := raw["udpRouters"]; ok {
		typedRouters := convertToTyped[dynamic.UDPRouter](routers)
		udpConfig.Routers = matchers.UDPRouters(typedRouters, pc.Routers, providerMatcher)
	}
	for _, extra := range pc.Routers.ExtraRoutes {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var router dynamic.UDPRouter
		if err := json.Unmarshal(b, &router); err != nil {
			continue
		}
		if routerName, ok := extra.(map[string]interface{})["name"].(string); ok {
			udpConfig.Routers[routerName] = &router
		}
	}
	overrides.StripProvidersUDP(udpConfig)
	overrides.OverrideUDPRouters(udpConfig.Routers, pc.Routers.Overrides)
}

func processUDPServices(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, pc *config.UDPSection, providerMatcher string) {
	if services, ok := raw["udpServices"]; ok {
		typedServices := convertToTyped[dynamic.UDPService](services)
		udpConfig.Services = matchers.UDPServices(typedServices, pc.Services, providerMatcher)
	}
	for _, extra := range pc.Services.ExtraServices {
		b, err := json.Marshal(extra)
		if err != nil {
			continue
		}
		var service dynamic.UDPService
		if err := json.Unmarshal(b, &service); err != nil {
			continue
		}
		if serviceName, ok := extra.(map[string]interface{})["name"].(string); ok {
			udpConfig.Services[serviceName] = &service
		}
	}
	overrides.StripProvidersUDP(udpConfig)
	overrides.OverrideUDPServices(udpConfig.Services, pc.Services.Overrides)
}

// ParseTLSConfig fills tlsConfig from raw data according to providerConfig.
func ParseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration, providerConfig *config.TLSSection) {
	if certificates, ok := raw["tlsCertificates"]; ok {
		tlsConfig.Certificates = matchers.TLSCertificates(certificates, providerConfig)
	}
	if options, ok := raw["tlsOptions"]; ok {
		tlsConfig.Options = matchers.TLSOptions(options, providerConfig)
	}
	if stores, ok := raw["tlsStores"]; ok {
		tlsConfig.Stores = matchers.TLSStores(stores, providerConfig)
	}
}

func parseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration, providerConfig *config.TLSSection) {
	ParseTLSConfig(raw, tlsConfig, providerConfig)
}

// parseDynamicConfiguration mirrors the behavior used by the httpclient to build a dynamic.Configuration
// from raw JSON and a provider configuration. This is used by tests to validate nil-section handling.
func parseDynamicConfiguration(body []byte, providerCfg *config.ProviderConfig) (*dynamic.Configuration, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return &dynamic.Configuration{}, err
	}

	httpConfig := &dynamic.HTTPConfiguration{}
	tcpConfig := &dynamic.TCPConfiguration{}
	udpConfig := &dynamic.UDPConfiguration{}
	tlsConfig := &dynamic.TLSConfiguration{}

	// Ensure defaults similar to httpclient.ensureProviderDefaults
	if providerCfg.HTTP == nil {
		providerCfg.HTTP = &config.HTTPSection{Discover: true}
	}
	if providerCfg.TCP == nil {
		providerCfg.TCP = &config.TCPSection{Discover: true}
	}
	if providerCfg.UDP == nil {
		providerCfg.UDP = &config.UDPSection{Discover: true}
	}
	if providerCfg.TLS == nil {
		providerCfg.TLS = &config.TLSSection{Discover: true}
	}
	if providerCfg.Tunnels == nil {
		providerCfg.Tunnels = []config.TunnelConfig{}
	}

	if providerCfg.HTTP.Discover {
		ParseHTTPConfig(raw, httpConfig, providerCfg.HTTP, providerCfg.Matcher, providerCfg.Tunnels)
	}
	if providerCfg.TCP.Discover {
		ParseTCPConfig(raw, tcpConfig, providerCfg.TCP, providerCfg.Matcher, providerCfg.Tunnels)
	}
	if providerCfg.UDP.Discover {
		ParseUDPConfig(raw, udpConfig, providerCfg.UDP, providerCfg.Matcher)
	}
	if providerCfg.TLS.Discover {
		ParseTLSConfig(raw, tlsConfig, providerCfg.TLS)
	}

	return &dynamic.Configuration{
		HTTP: httpConfig,
		TCP:  tcpConfig,
		UDP:  udpConfig,
		TLS:  tlsConfig,
	}, nil
}
