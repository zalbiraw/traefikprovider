// Package httpclient fetches and parses dynamic configuration from remote providers.
package httpclient

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/filters"
	"github.com/zalbiraw/traefik-provider/internal/overrides"
)

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

func parseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, providerConfig *config.HTTPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	ensureHTTPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processHTTPRouters(raw, httpConfig, providerConfig, pf)
	}
	if providerConfig.Services.Discover {
		processHTTPServices(raw, httpConfig, providerConfig, pf, tunnels)
	}
	if providerConfig.Middlewares.Discover {
		processHTTPMiddlewares(raw, httpConfig, providerConfig, pf)
	}
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

func processHTTPRouters(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, pf config.ProviderFilter) {
	if routers, ok := raw["routers"]; ok {
		typedRouters := convertToTyped[dynamic.Router](routers)
		httpConfig.Routers = filters.HTTPRouters(typedRouters, pc.Routers, pf)
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

func processHTTPServices(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	if services, ok := raw["services"]; ok {
		typedServices := convertToTyped[dynamic.Service](services)
		httpConfig.Services = filters.HTTPServices(typedServices, pc.Services, pf)
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
	overrides.OverrideHTTPServices(httpConfig.Services, pc.Services.Overrides, tunnels)
}

func processHTTPMiddlewares(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, pc *config.HTTPSection, pf config.ProviderFilter) {
	if middlewares, ok := raw["middlewares"]; ok {
		typedMiddlewares := convertToTyped[dynamic.Middleware](middlewares)
		httpConfig.Middlewares = filters.HTTPMiddlewares(typedMiddlewares, pc.Middlewares, pf)
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

func parseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, providerConfig *config.TCPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	ensureTCPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processTCPRouters(raw, tcpConfig, providerConfig, pf)
	}
	if providerConfig.Services.Discover {
		processTCPServices(raw, tcpConfig, providerConfig, pf, tunnels)
	}
	if providerConfig.Middlewares.Discover {
		processTCPMiddlewares(raw, tcpConfig, providerConfig, pf)
	}
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

func processTCPRouters(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, pf config.ProviderFilter) {
	if routers, ok := raw["tcpRouters"]; ok {
		typedRouters := convertToTyped[dynamic.TCPRouter](routers)
		tcpConfig.Routers = filters.TCPRouters(typedRouters, pc.Routers, pf)
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

func processTCPServices(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	if services, ok := raw["tcpServices"]; ok {
		typedServices := convertToTyped[dynamic.TCPService](services)
		tcpConfig.Services = filters.TCPServices(typedServices, pc.Services, pf)
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
	overrides.OverrideTCPServices(tcpConfig.Services, pc.Services.Overrides, tunnels)
}

func processTCPMiddlewares(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, pc *config.TCPSection, pf config.ProviderFilter) {
	if middlewares, ok := raw["tcpMiddlewares"]; ok {
		typedMiddlewares := convertToTyped[dynamic.TCPMiddleware](middlewares)
		tcpConfig.Middlewares = filters.TCPMiddlewares(typedMiddlewares, pc.Middlewares, pf)
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

func parseUDPConfig(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, providerConfig *config.UDPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	ensureUDPDefaults(providerConfig)
	if providerConfig.Routers.Discover {
		processUDPRouters(raw, udpConfig, providerConfig, pf)
	}
	if providerConfig.Services.Discover {
		processUDPServices(raw, udpConfig, providerConfig, pf, tunnels)
	}
}

func ensureUDPDefaults(pc *config.UDPSection) {
	if pc.Routers == nil {
		pc.Routers = &config.UDPRoutersConfig{Discover: true}
	}
	if pc.Services == nil {
		pc.Services = &config.UDPServicesConfig{Discover: true}
	}
}

func processUDPRouters(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, pc *config.UDPSection, pf config.ProviderFilter) {
	if routers, ok := raw["udpRouters"]; ok {
		typedRouters := convertToTyped[dynamic.UDPRouter](routers)
		udpConfig.Routers = filters.UDPRouters(typedRouters, pc.Routers, pf)
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

func processUDPServices(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, pc *config.UDPSection, pf config.ProviderFilter, tunnels []config.TunnelConfig) {
	if services, ok := raw["udpServices"]; ok {
		typedServices := convertToTyped[dynamic.UDPService](services)
		udpConfig.Services = filters.UDPServices(typedServices, pc.Services, pf)
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
	overrides.OverrideUDPServices(udpConfig.Services, pc.Services.Overrides, tunnels)
}

func parseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration, providerConfig *config.TLSSection) {
	if certificates, ok := raw["tlsCertificates"]; ok {
		tlsConfig.Certificates = filters.TLSCertificates(certificates, providerConfig)
	}
	if options, ok := raw["tlsOptions"]; ok {
		tlsConfig.Options = filters.TLSOptions(options, providerConfig)
	}
	if stores, ok := raw["tlsStores"]; ok {
		tlsConfig.Stores = filters.TLSStores(stores, providerConfig)
	}
}
