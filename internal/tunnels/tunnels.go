// Package tunnels contains helpers to apply tunnel configuration to services.
package tunnels

import (
	"fmt"
	"hash/fnv"

	"github.com/traefik/genconf/dynamic"
	dtls "github.com/traefik/genconf/dynamic/tls"
	"github.com/zalbiraw/traefikprovider/config"
	"github.com/zalbiraw/traefikprovider/internal/matchers"
)

// ApplyHTTPTunnels finds services matching each tunnel's matcher and:
// - replaces servers with tunnel addresses
// - if tunnel has mTLS, creates a ServersTransport and references it from the service.
func ApplyHTTPTunnels(httpConfig *dynamic.HTTPConfiguration, providerMatcher string, tns []config.TunnelConfig) {
	if httpConfig == nil || httpConfig.Services == nil {
		return
	}
	for _, t := range tns {
		if len(t.Addresses) == 0 {
			continue
		}
		// Match services using the tunnel's matcher
		sel := &config.ServicesConfig{Matcher: t.Matcher}
		// Important: ignore provider-level matcher here because provider suffixes may
		// have been stripped from names during overrides; use only the tunnel matcher.
		matched := matchers.HTTPServices(httpConfig.Services, sel, "")

		// Prepare optional ServersTransport if MTLS is provided
		var transportName string
		if t.MTLS != nil {
			ensureHTTPServersTransports(httpConfig)
			transportName = serversTransportName(t.Matcher)
			if _, exists := httpConfig.ServersTransports[transportName]; !exists {
				httpConfig.ServersTransports[transportName] = buildServersTransportFromMTLS(t.MTLS)
			}
		}

		// Update matched services
		for name, svc := range matched {
			// Ensure LoadBalancer exists
			if svc.LoadBalancer == nil {
				svc.LoadBalancer = &dynamic.ServersLoadBalancer{}
			}
			svc.LoadBalancer.Servers = buildHTTPServers(t.Addresses)
			if transportName != "" {
				svc.LoadBalancer.ServersTransport = transportName
			}
			httpConfig.Services[name] = svc

			// Optionally strip router TLS options for routers referencing this service
			if t.MTLS != nil && shouldStripRouterTLSOptions(t.MTLS) {
				stripRoutersTLSForService(httpConfig, name)
			}
		}
	}
}

func ensureHTTPServersTransports(httpConfig *dynamic.HTTPConfiguration) {
	if httpConfig.ServersTransports == nil {
		httpConfig.ServersTransports = map[string]*dynamic.ServersTransport{}
	}
}

// buildHTTPServers converts URLs to []dynamic.Server.
func buildHTTPServers(urls []string) []dynamic.Server {
	if len(urls) == 0 {
		return nil
	}
	servers := make([]dynamic.Server, 0, len(urls))
	for _, addr := range urls {
		servers = append(servers, dynamic.Server{URL: addr})
	}
	return servers
}

// buildServersTransportFromMTLS maps provider MTLS config to dynamic.ServersTransport.
func buildServersTransportFromMTLS(m *config.MTLSConfig) *dynamic.ServersTransport {
	st := &dynamic.ServersTransport{}
	if m.CAFile != "" {
		st.RootCAs = []string{m.CAFile}
	}
	if m.CertFile != "" && m.KeyFile != "" {
		st.Certificates = dtls.Certificates{
			{CertFile: m.CertFile, KeyFile: m.KeyFile},
		}
	}
	return st
}

// serversTransportName generates a stable name from the tunnel matcher.
func serversTransportName(matcher string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(matcher))
	return fmt.Sprintf("st-%08x", h.Sum32())
}

// shouldStripRouterTLSOptions returns true when the MTLS config indicates TLS options
// should be stripped from routers (default true when unset).
func shouldStripRouterTLSOptions(m *config.MTLSConfig) bool {
	if m == nil || m.StripRouterTLSOptions == nil {
		return true
	}
	return *m.StripRouterTLSOptions
}

// stripRoutersTLSForService inspects routers that use the given service name and
// mutates their TLS config according to the rule:
// - if TLS only has Options set (no certResolver, no domains), remove the TLS block
// - else clear only the Options field
func stripRoutersTLSForService(httpConfig *dynamic.HTTPConfiguration, serviceName string) {
	if httpConfig == nil || httpConfig.Routers == nil {
		return
	}
	for rname, r := range httpConfig.Routers {
		if r == nil || r.Service != serviceName || r.TLS == nil {
			continue
		}
		// Determine if only Options is set
		onlyOptions := r.TLS != nil && r.TLS.Options != "" && r.TLS.CertResolver == "" && len(r.TLS.Domains) == 0
		if onlyOptions {
			r.TLS = nil
		} else {
			// Remove just the Options
			r.TLS.Options = ""
		}
		httpConfig.Routers[rname] = r
	}
}

// ApplyTCPTunnels finds TCP services matching each tunnel's matcher and replaces servers with tunnel addresses.
func ApplyTCPTunnels(tcpConfig *dynamic.TCPConfiguration, providerMatcher string, tns []config.TunnelConfig) {
	if tcpConfig == nil || tcpConfig.Services == nil {
		return
	}
	for _, t := range tns {
		if len(t.Addresses) == 0 {
			continue
		}
		sel := &config.ServicesConfig{Matcher: t.Matcher}
		matched := matchers.TCPServices(tcpConfig.Services, sel, "")
		for name, svc := range matched {
			if svc.LoadBalancer == nil {
				svc.LoadBalancer = &dynamic.TCPServersLoadBalancer{}
			}
			svc.LoadBalancer.Servers = buildTCPServers(t.Addresses)
			tcpConfig.Services[name] = svc
		}
	}
}

// buildTCPServers converts addresses to []dynamic.TCPServer.
func buildTCPServers(urls []string) []dynamic.TCPServer {
	if len(urls) == 0 {
		return nil
	}
	servers := make([]dynamic.TCPServer, 0, len(urls))
	for _, addr := range urls {
		servers = append(servers, dynamic.TCPServer{Address: addr})
	}
	return servers
}
