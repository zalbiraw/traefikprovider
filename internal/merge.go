// Package internal provides internal utilities for merging configurations.
package internal

import (
	"github.com/traefik/genconf/dynamic"
	tlstypes "github.com/traefik/genconf/dynamic/tls"
)

// MergeConfigurations merges multiple dynamic.Configurations into one.
func MergeConfigurations(configs ...*dynamic.Configuration) *dynamic.Configuration {
	merged := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{Routers: map[string]*dynamic.Router{}, Services: map[string]*dynamic.Service{}, Middlewares: map[string]*dynamic.Middleware{}},
		TCP:  &dynamic.TCPConfiguration{Routers: map[string]*dynamic.TCPRouter{}, Services: map[string]*dynamic.TCPService{}, Middlewares: map[string]*dynamic.TCPMiddleware{}},
		UDP:  &dynamic.UDPConfiguration{Routers: map[string]*dynamic.UDPRouter{}, Services: map[string]*dynamic.UDPService{}},
		TLS:  &dynamic.TLSConfiguration{Certificates: []*tlstypes.CertAndStores{}, Options: map[string]tlstypes.Options{}, Stores: map[string]tlstypes.Store{}},
	}
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		mergeHTTP(merged, cfg)
		mergeTCP(merged, cfg)
		mergeUDP(merged, cfg)
		mergeTLS(merged, cfg)
	}
	return merged
}

func mergeHTTP(dst, src *dynamic.Configuration) {
	if src.HTTP == nil {
		return
	}
	for k, v := range src.HTTP.Routers {
		dst.HTTP.Routers[k] = v
	}
	for k, v := range src.HTTP.Services {
		dst.HTTP.Services[k] = v
	}
	for k, v := range src.HTTP.Middlewares {
		dst.HTTP.Middlewares[k] = v
	}
}

func mergeTCP(dst, src *dynamic.Configuration) {
	if src.TCP == nil {
		return
	}
	for k, v := range src.TCP.Routers {
		dst.TCP.Routers[k] = v
	}
	for k, v := range src.TCP.Services {
		dst.TCP.Services[k] = v
	}
	for k, v := range src.TCP.Middlewares {
		dst.TCP.Middlewares[k] = v
	}
}

func mergeUDP(dst, src *dynamic.Configuration) {
	if src.UDP == nil {
		return
	}
	for k, v := range src.UDP.Routers {
		dst.UDP.Routers[k] = v
	}
	for k, v := range src.UDP.Services {
		dst.UDP.Services[k] = v
	}
}

func mergeTLS(dst, src *dynamic.Configuration) {
	if src.TLS == nil {
		return
	}
	dst.TLS.Certificates = append(dst.TLS.Certificates, src.TLS.Certificates...)
	for k, v := range src.TLS.Options {
		dst.TLS.Options[k] = v
	}
	for k, v := range src.TLS.Stores {
		dst.TLS.Stores[k] = v
	}
}
