package httpclient

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func parseHTTPConfig(raw map[string]interface{}, httpConfig *dynamic.HTTPConfiguration, providerConfig *config.HTTPSection) error {
	if providerConfig.Routers.Discover {
		if routers, ok := raw["routers"]; ok {
			httpConfig.Routers = filterRouters(routers, providerConfig.Routers)
		}
	}
	if providerConfig.Services.Discover {
		if services, ok := raw["services"]; ok {
			httpConfig.Services = filterServices(services, providerConfig.Services)
		}
	}
	if providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["middlewares"]; ok {
			httpConfig.Middlewares = filterMiddlewares(middlewares, providerConfig.Middlewares)
		}
	}
	if providerConfig.ServerTransports.Discover {
		if serversTransports, ok := raw["serversTransports"]; ok {
			httpConfig.ServersTransports = filterServerTransports(serversTransports, &providerConfig.ServerTransports)
		}
	}
	return nil
}

func parseTCPConfig(raw map[string]interface{}, tcpConfig *dynamic.TCPConfiguration, providerConfig *config.TCPSection) error {
	if providerConfig.Routers.Discover {
		if routers, ok := raw["tcpRouters"]; ok {
			tcpConfig.Routers = filterTCPRouters(routers, providerConfig.Routers)
		}
	}
	if providerConfig.Services.Discover {
		if services, ok := raw["tcpServices"]; ok {
			tcpConfig.Services = filterTCPServices(services, providerConfig.Services)
		}
	}
	if providerConfig.Middlewares.Discover {
		if middlewares, ok := raw["tcpMiddlewares"]; ok {
			tcpConfig.Middlewares = filterTCPMiddlewares(middlewares, providerConfig.Middlewares)
		}
	}
	return nil
}

func parseUDPConfig(raw map[string]interface{}, udpConfig *dynamic.UDPConfiguration, providerConfig *config.UDPSection) error {
	if providerConfig.Routers.Discover {
		if routers, ok := raw["udpRouters"]; ok {
			udpConfig.Routers = filterUDPRouters(routers, providerConfig.Routers)
		}
	}
	if providerConfig.Services.Discover {
		if services, ok := raw["udpServices"]; ok {
			udpConfig.Services = filterUDPServices(services, providerConfig.Services)
		}
	}
	return nil
}

func parseTLSConfig(raw map[string]interface{}, tlsConfig *dynamic.TLSConfiguration, providerConfig *config.TLSSection) error {
	if certificates, ok := raw["tlsCertificates"]; ok {
		tlsConfig.Certificates = filterTLSCertificates(certificates, providerConfig)
	}
	if options, ok := raw["tlsOptions"]; ok {
		tlsConfig.Options = filterTLSOptions(options, providerConfig)
	}
	if stores, ok := raw["tlsStores"]; ok {
		tlsConfig.Stores = filterTLSStores(stores, providerConfig)
	}
	return nil
}
