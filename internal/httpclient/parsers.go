package httpclient

import (
	"encoding/json"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/filters"
	"github.com/zalbiraw/traefik-provider/internal/overrides"
)

func parseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, providerConfig *config.HTTPSection, tunnels []config.TunnelConfig) error {
	if providerConfig.Routers != nil && providerConfig.Routers.Discover {
		if routers, ok := raw["routers"]; ok {
			httpConfig.Routers = filters.HTTPRouters(routers, providerConfig.Routers)
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
	if providerConfig.Services != nil && providerConfig.Services.Discover {
		if services, ok := raw["services"]; ok {
			httpConfig.Services = filters.HTTPServices(services, providerConfig.Services)
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
	if providerConfig.Middlewares != nil && providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["middlewares"]; ok {
			httpConfig.Middlewares = filters.HTTPMiddlewares(middlewares, providerConfig.Middlewares)
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
	if providerConfig.ServerTransports.Discover {
		if serversTransports, ok := raw["serversTransports"]; ok {
			httpConfig.ServersTransports = filters.HTTPServerTransports(serversTransports, &providerConfig.ServerTransports)
		}
		// Append extraServerTransports
		for _, extra := range providerConfig.ServerTransports.ExtraServerTransports {
			b, err := json.Marshal(extra)
			if err != nil {
				continue
			}
			var st dynamic.ServersTransport
			if err := json.Unmarshal(b, &st); err != nil {
				continue
			}
			if stName, ok := extra.(map[string]interface{})["name"].(string); ok {
				httpConfig.ServersTransports[stName] = &st
			}
		}
	}
	return nil
}

func parseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, providerConfig *config.TCPSection, tunnels []config.TunnelConfig) error {
	if providerConfig.Routers != nil && providerConfig.Routers.Discover {
		if routers, ok := raw["tcpRouters"]; ok {
			tcpConfig.Routers = filters.TCPRouters(routers, providerConfig.Routers)
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
	if providerConfig.Services != nil && providerConfig.Services.Discover {
		if services, ok := raw["tcpServices"]; ok {
			tcpConfig.Services = filters.TCPServices(services, providerConfig.Services)
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
	if providerConfig.Middlewares != nil && providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["tcpMiddlewares"]; ok {
			tcpConfig.Middlewares = filters.TCPMiddlewares(middlewares, providerConfig.Middlewares)
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
	if providerConfig.Routers != nil && providerConfig.Routers.Discover {
		if routers, ok := raw["udpRouters"]; ok {
			udpConfig.Routers = filters.UDPRouters(routers, providerConfig.Routers)
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
	if providerConfig.Services != nil && providerConfig.Services.Discover {
		if services, ok := raw["udpServices"]; ok {
			udpConfig.Services = filters.UDPServices(services, providerConfig.Services)
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
