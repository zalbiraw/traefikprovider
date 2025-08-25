package httpclient

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/filters"
	"github.com/zalbiraw/traefik-provider/internal/overrides"
)

func convertToTyped[T any](data interface{}) map[string]*T {
	result := make(map[string]*T)
	if dataMap, ok := data.(map[string]interface{}); ok {
		for name, itemData := range dataMap {
			item := new(T)
			if itemMap, ok := itemData.(map[string]interface{}); ok {
				b, _ := json.Marshal(itemMap)
				json.Unmarshal(b, item)
				result[name] = item
			}
		}
	}
	return result
}

func parseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, providerConfig *config.HTTPSection, tunnels []config.TunnelConfig) error {
	if providerConfig.Routers == nil {
		providerConfig.Routers = &config.RoutersConfig{Discover: true}
	}
	if providerConfig.Services == nil {
		providerConfig.Services = &config.ServicesConfig{Discover: true}
	}
	if providerConfig.Middlewares == nil {
		providerConfig.Middlewares = &config.MiddlewaresConfig{Discover: true}
	}

	if providerConfig.Routers.Discover {
		if routers, ok := raw["routers"]; ok {
			typedRouters := convertToTyped[dynamic.Router](routers)
			httpConfig.Routers = filters.HTTPRouters(typedRouters, providerConfig.Routers)
		}
		for _, extra := range providerConfig.Routers.ExtraRoutes {
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
		overrides.OverrideHTTPRouters(httpConfig.Routers, providerConfig.Routers.Overrides)
	}

	if providerConfig.Services.Discover {
		if services, ok := raw["services"]; ok {
			typedServices := convertToTyped[dynamic.Service](services)
			httpConfig.Services = filters.HTTPServices(typedServices, providerConfig.Services)
		}
		for _, extra := range providerConfig.Services.ExtraServices {
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
		overrides.OverrideHTTPServices(httpConfig.Services, providerConfig.Services.Overrides, tunnels)
	}

	if providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["middlewares"]; ok {
			typedMiddlewares := convertToTyped[dynamic.Middleware](middlewares)
			httpConfig.Middlewares = filters.HTTPMiddlewares(typedMiddlewares, providerConfig.Middlewares)
		}
		for _, extra := range providerConfig.Middlewares.ExtraMiddlewares {
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
	}
	return nil
}

func parseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, providerConfig *config.TCPSection, tunnels []config.TunnelConfig) error {
	if providerConfig.Routers == nil {
		providerConfig.Routers = &config.RoutersConfig{Discover: true}
	}
	if providerConfig.Services == nil {
		providerConfig.Services = &config.ServicesConfig{Discover: true}
	}
	if providerConfig.Middlewares == nil {
		providerConfig.Middlewares = &config.MiddlewaresConfig{Discover: true}
	}

	if providerConfig.Routers.Discover {
		if routers, ok := raw["tcpRouters"]; ok {
			typedRouters := convertToTyped[dynamic.TCPRouter](routers)
			tcpConfig.Routers = filters.TCPRouters(typedRouters, providerConfig.Routers)
		}
		for _, extra := range providerConfig.Routers.ExtraRoutes {
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
		overrides.OverrideTCPRouters(tcpConfig.Routers, providerConfig.Routers.Overrides)
	}

	if providerConfig.Services.Discover {
		if services, ok := raw["tcpServices"]; ok {
			typedServices := convertToTyped[dynamic.TCPService](services)
			tcpConfig.Services = filters.TCPServices(typedServices, providerConfig.Services)
		}
		for _, extra := range providerConfig.Services.ExtraServices {
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
		overrides.OverrideTCPServices(tcpConfig.Services, providerConfig.Services.Overrides, tunnels)
	}

	if providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["tcpMiddlewares"]; ok {
			typedMiddlewares := convertToTyped[dynamic.TCPMiddleware](middlewares)
			tcpConfig.Middlewares = filters.TCPMiddlewares(typedMiddlewares, providerConfig.Middlewares)
		}
		for _, extra := range providerConfig.Middlewares.ExtraMiddlewares {
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
	}
	return nil
}

func parseUDPConfig(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, providerConfig *config.UDPSection, tunnels []config.TunnelConfig) error {
	if providerConfig.Routers == nil {
		providerConfig.Routers = &config.UDPRoutersConfig{Discover: true}
	}
	if providerConfig.Services == nil {
		providerConfig.Services = &config.UDPServicesConfig{Discover: true}
	}

	if providerConfig.Routers.Discover {
		if routers, ok := raw["udpRouters"]; ok {
			typedRouters := convertToTyped[dynamic.UDPRouter](routers)
			udpConfig.Routers = filters.UDPRouters(typedRouters, providerConfig.Routers)
		}
		// Append extraRoutes
		for _, extra := range providerConfig.Routers.ExtraRoutes {
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
		overrides.OverrideUDPRouters(udpConfig.Routers, providerConfig.Routers.Overrides)
	}

	if providerConfig.Services.Discover {
		if services, ok := raw["udpServices"]; ok {
			typedServices := convertToTyped[dynamic.UDPService](services)
			udpConfig.Services = filters.UDPServices(typedServices, providerConfig.Services)
		}
		// Append extraServices
		for _, extra := range providerConfig.Services.ExtraServices {
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
		overrides.OverrideUDPServices(udpConfig.Services, providerConfig.Services.Overrides, tunnels)
	}
	return nil
}

func parseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration, providerConfig *config.TLSSection) error {
	if certificates, ok := raw["tlsCertificates"]; ok {
		tlsConfig.Certificates = filters.TLSCertificates(certificates, providerConfig)
	}
	if options, ok := raw["tlsOptions"]; ok {
		tlsConfig.Options = filters.TLSOptions(options, providerConfig)
	}
	if stores, ok := raw["tlsStores"]; ok {
		tlsConfig.Stores = filters.TLSStores(stores, providerConfig)
	}
	return nil
}
