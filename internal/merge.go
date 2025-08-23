package internal

import (
	"github.com/traefik/genconf/dynamic"

	tlstypes "github.com/traefik/genconf/dynamic/tls"
)

// MergeConfigurations merges multiple dynamic.Configurations into one.
func MergeConfigurations(configs ...*dynamic.Configuration) *dynamic.Configuration {
	merged := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{Routers: map[string]*dynamic.Router{}, Services: map[string]*dynamic.Service{}, Middlewares: map[string]*dynamic.Middleware{}, ServersTransports: map[string]*dynamic.ServersTransport{}},
		TCP:  &dynamic.TCPConfiguration{Routers: map[string]*dynamic.TCPRouter{}, Services: map[string]*dynamic.TCPService{}, Middlewares: map[string]*dynamic.TCPMiddleware{}},
		UDP:  &dynamic.UDPConfiguration{Routers: map[string]*dynamic.UDPRouter{}, Services: map[string]*dynamic.UDPService{}},
		TLS:  &dynamic.TLSConfiguration{Certificates: []*tlstypes.CertAndStores{}, Options: map[string]tlstypes.Options{}, Stores: map[string]tlstypes.Store{}},
	}
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		if cfg.HTTP != nil {
			for k, v := range cfg.HTTP.Routers {
				merged.HTTP.Routers[k] = v
			}
			for k, v := range cfg.HTTP.Services {
				merged.HTTP.Services[k] = v
			}
			for k, v := range cfg.HTTP.Middlewares {
				merged.HTTP.Middlewares[k] = v
			}
			for k, v := range cfg.HTTP.ServersTransports {
				merged.HTTP.ServersTransports[k] = v
			}
		}
		if cfg.TCP != nil {
			for k, v := range cfg.TCP.Routers {
				merged.TCP.Routers[k] = v
			}
			for k, v := range cfg.TCP.Services {
				merged.TCP.Services[k] = v
			}
			for k, v := range cfg.TCP.Middlewares {
				merged.TCP.Middlewares[k] = v
			}
		}
		if cfg.UDP != nil {
			for k, v := range cfg.UDP.Routers {
				merged.UDP.Routers[k] = v
			}
			for k, v := range cfg.UDP.Services {
				merged.UDP.Services[k] = v
			}
		}
		if cfg.TLS != nil {
			merged.TLS.Certificates = append(merged.TLS.Certificates, cfg.TLS.Certificates...)
			for k, v := range cfg.TLS.Options {
				merged.TLS.Options[k] = v
			}
			for k, v := range cfg.TLS.Stores {
				merged.TLS.Stores[k] = v
			}
		}
	}
	return merged
}
