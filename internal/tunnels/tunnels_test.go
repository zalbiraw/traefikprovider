package tunnels

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/traefik/genconf/dynamic/types"
	"github.com/zalbiraw/traefikprovider/config"
)

//nolint:gocognit,gocyclo // test covers multiple branches in a single scenario for clarity
func TestApplyHTTPTunnels_StripsRouterTLSOptions_DefaultTrue(t *testing.T) {
	httpCfg := &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"r1": {Service: "svc1", TLS: &dynamic.RouterTLSConfig{Options: "tls@file"}},
			"r2": {Service: "svc1", TLS: &dynamic.RouterTLSConfig{Options: "tls@file", CertResolver: "le", Domains: []types.Domain{{Main: "example.com"}}}},
			"r3": {Service: "svc2", TLS: &dynamic.RouterTLSConfig{Options: "tls@file"}},
			"r4": {Service: "svc3", TLS: &dynamic.RouterTLSConfig{Options: "tls@file"}},
		},
		Services: map[string]*dynamic.Service{
			"svc1": {LoadBalancer: &dynamic.ServersLoadBalancer{}},
			"svc2": {LoadBalancer: &dynamic.ServersLoadBalancer{}},
			// svc3 intentionally has nil LoadBalancer to cover creation path
			"svc3": {},
		},
	}

	tunnel := config.TunnelConfig{
		Addresses: []string{"http://tunnel:443"},
		Matcher:   "Name(`svc1`) || Name(`svc3`)",
		MTLS:      &config.MTLSConfig{CAFile: "ca.crt", CertFile: "crt.pem", KeyFile: "key.pem"},
	}

	ApplyHTTPTunnels(httpCfg, "", []config.TunnelConfig{tunnel})

	// Service updated
	svc := httpCfg.Services["svc1"]
	if svc == nil || svc.LoadBalancer == nil || len(svc.LoadBalancer.Servers) != 1 || svc.LoadBalancer.Servers[0].URL != "http://tunnel:443" {
		t.Fatalf("expected svc1 to have 1 server set to tunnel, got: %+v", svc)
	}
	if svc.LoadBalancer.ServersTransport == "" {
		t.Fatalf("expected ServersTransport to be set when MTLS provided")
	}

	// ServersTransport created
	if httpCfg.ServersTransports == nil || httpCfg.ServersTransports[svc.LoadBalancer.ServersTransport] == nil {
		t.Fatalf("expected servers transport %q to exist", svc.LoadBalancer.ServersTransport)
	}

	// r1: only Options -> TLS removed
	if httpCfg.Routers["r1"].TLS != nil {
		t.Fatalf("expected r1 TLS to be removed, got: %+v", httpCfg.Routers["r1"].TLS)
	}
	// r2: Options cleared but other fields kept
	r2tls := httpCfg.Routers["r2"].TLS
	if r2tls == nil || r2tls.Options != "" || r2tls.CertResolver != "le" || len(r2tls.Domains) != 1 || r2tls.Domains[0].Main != "example.com" {
		t.Fatalf("unexpected r2 TLS after strip: %+v", r2tls)
	}
	// r3: different service -> untouched
	if httpCfg.Routers["r3"].TLS == nil || httpCfg.Routers["r3"].TLS.Options != "tls@file" {
		t.Fatalf("expected r3 TLS unchanged, got: %+v", httpCfg.Routers["r3"].TLS)
	}
	// r4/svc3: LoadBalancer should have been created and server set
	svc3 := httpCfg.Services["svc3"]
	if svc3 == nil || svc3.LoadBalancer == nil || len(svc3.LoadBalancer.Servers) != 1 {
		t.Fatalf("expected svc3 loadBalancer created with one server, got: %+v", svc3)
	}
}

func TestApplyHTTPTunnels_NoStrip_WhenFlagFalse(t *testing.T) {
	httpCfg := &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"r1": {Service: "svc1", TLS: &dynamic.RouterTLSConfig{Options: "tls@file"}},
		},
		Services: map[string]*dynamic.Service{
			"svc1": {LoadBalancer: &dynamic.ServersLoadBalancer{}},
		},
	}

	strip := false
	tunnel := config.TunnelConfig{
		Addresses: []string{"http://tunnel:443"},
		Matcher:   "Name(`svc1`)",
		MTLS:      &config.MTLSConfig{StripRouterTLSOptions: &strip},
	}

	ApplyHTTPTunnels(httpCfg, "", []config.TunnelConfig{tunnel})

	if httpCfg.Routers["r1"].TLS == nil || httpCfg.Routers["r1"].TLS.Options != "tls@file" {
		t.Fatalf("expected r1 TLS untouched when strip=false, got: %+v", httpCfg.Routers["r1"].TLS)
	}
}

func TestShouldStripRouterTLSOptions(t *testing.T) {
	if !shouldStripRouterTLSOptions(nil) {
		t.Fatal("nil config should default to true")
	}
	b := true
	if !shouldStripRouterTLSOptions(&config.MTLSConfig{StripRouterTLSOptions: &b}) {
		t.Fatal("expected true when flag=true")
	}
	b = false
	if shouldStripRouterTLSOptions(&config.MTLSConfig{StripRouterTLSOptions: &b}) {
		t.Fatal("expected false when flag=false")
	}
}

func TestBuildHTTPServers(t *testing.T) {
	if res := buildHTTPServers(nil); res != nil {
		t.Fatalf("expected nil for empty input, got: %+v", res)
	}
	res := buildHTTPServers([]string{"http://a", "http://b"})
	if len(res) != 2 || res[0].URL != "http://a" || res[1].URL != "http://b" {
		t.Fatalf("unexpected servers: %+v", res)
	}
}

func TestApplyTCPTunnels_Basic(t *testing.T) {
	tcpCfg := &dynamic.TCPConfiguration{
		Routers: map[string]*dynamic.TCPRouter{
			"tr":  {Service: "tsvc"},
			"tr2": {Service: "tsvc2"},
		},
		Services: map[string]*dynamic.TCPService{
			"tsvc": {LoadBalancer: &dynamic.TCPServersLoadBalancer{}},
			// tsvc2 without LoadBalancer to cover creation path
			"tsvc2": {},
		},
	}
	tunnel := config.TunnelConfig{Addresses: []string{"tcp://tunnel:8443"}, Matcher: "Name(`tsvc`) || Name(`tsvc2`)"}
	ApplyTCPTunnels(tcpCfg, "", []config.TunnelConfig{tunnel})

	svc := tcpCfg.Services["tsvc"]
	if svc == nil || svc.LoadBalancer == nil || len(svc.LoadBalancer.Servers) != 1 || svc.LoadBalancer.Servers[0].Address != "tcp://tunnel:8443" {
		t.Fatalf("unexpected tcp service after apply: %+v", svc)
	}
	svc2 := tcpCfg.Services["tsvc2"]
	if svc2 == nil || svc2.LoadBalancer == nil || len(svc2.LoadBalancer.Servers) != 1 || svc2.LoadBalancer.Servers[0].Address != "tcp://tunnel:8443" {
		t.Fatalf("unexpected tcp service after apply for tsvc2: %+v", svc2)
	}
}

func TestBuildServersTransportFromMTLS(t *testing.T) {
	m := &config.MTLSConfig{CAFile: "ca.pem", CertFile: "crt.pem", KeyFile: "key.pem"}
	st := buildServersTransportFromMTLS(m)
	if st == nil {
		t.Fatal("expected servers transport not nil")
	}
	if len(st.RootCAs) != 1 || st.RootCAs[0] != "ca.pem" {
		t.Fatalf("unexpected RootCAs: %+v", st.RootCAs)
	}
	if len(st.Certificates) != 1 || st.Certificates[0].CertFile != "crt.pem" || st.Certificates[0].KeyFile != "key.pem" {
		t.Fatalf("unexpected Certificates: %+v", st.Certificates)
	}
}

func TestServersTransportNameStableAndDifferent(t *testing.T) {
	a := serversTransportName("Name(`a`)")
	b := serversTransportName("Name(`b`)")
	if a == "" || b == "" {
		t.Fatal("expected non-empty names")
	}
	if a == b {
		t.Fatalf("expected different names for different matchers: %s vs %s", a, b)
	}
	// Stability
	if a != serversTransportName("Name(`a`)") {
		t.Fatal("expected stable deterministic name")
	}
}

func TestStripRoutersTLSForService_EdgeCases(t *testing.T) {
	// Nil config
	stripRoutersTLSForService(nil, "svc")

	// No routers
	stripRoutersTLSForService(&dynamic.HTTPConfiguration{}, "svc")

	// Router without match or TLS remains unchanged
	h := &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"r": {Service: "other"},
		},
	}
	stripRoutersTLSForService(h, "svc")
	if h.Routers["r"].TLS != nil {
		t.Fatalf("expected TLS to remain nil: %+v", h.Routers["r"].TLS)
	}
}

func TestApplyHTTPTunnels_NoAddressesOrNilSections(t *testing.T) {
	// Nil Services -> no panic, no effect
	ApplyHTTPTunnels(&dynamic.HTTPConfiguration{}, "", []config.TunnelConfig{{Addresses: nil, Matcher: "Name(`x`)"}})

	// With Services but empty addresses -> ignore
	h := &dynamic.HTTPConfiguration{Services: map[string]*dynamic.Service{
		"svc": {LoadBalancer: &dynamic.ServersLoadBalancer{}},
	}}
	ApplyHTTPTunnels(h, "", []config.TunnelConfig{{Addresses: nil, Matcher: "Name(`svc`)"}})
	if h.Services["svc"].LoadBalancer != nil && len(h.Services["svc"].LoadBalancer.Servers) != 0 {
		t.Fatalf("expected no servers updated when addresses empty: %+v", h.Services["svc"].LoadBalancer.Servers)
	}
}

func TestEnsureHTTPServersTransports(t *testing.T) {
	h := &dynamic.HTTPConfiguration{}
	ensureHTTPServersTransports(h)
	if h.ServersTransports == nil {
		t.Fatal("expected serversTransports to be initialized")
	}
}

func TestBuildTCPServers(t *testing.T) {
	if res := buildTCPServers(nil); res != nil {
		t.Fatalf("expected nil for empty input, got: %+v", res)
	}
	res := buildTCPServers([]string{"a", "b"})
	if len(res) != 2 || res[0].Address != "a" || res[1].Address != "b" {
		t.Fatalf("unexpected tcp servers: %+v", res)
	}
}

func TestApplyTCPTunnels_NilServicesEarlyReturn(t *testing.T) {
	// No panic and no effect when Services is nil
	tcpCfg := &dynamic.TCPConfiguration{}
	ApplyTCPTunnels(tcpCfg, "", []config.TunnelConfig{{Addresses: []string{"tcp://x"}, Matcher: "Name(`svc`)"}})
	if tcpCfg.Services != nil {
		t.Fatalf("expected Services to remain nil, got: %+v", tcpCfg.Services)
	}
}

func TestApplyTCPTunnels_EmptyAddressesIgnored(t *testing.T) {
	tcpCfg := &dynamic.TCPConfiguration{
		Services: map[string]*dynamic.TCPService{
			"svc": {LoadBalancer: &dynamic.TCPServersLoadBalancer{}},
		},
	}
	// Empty addresses -> no server updates
	ApplyTCPTunnels(tcpCfg, "", []config.TunnelConfig{{Addresses: nil, Matcher: "Name(`svc`)"}})
	if lb := tcpCfg.Services["svc"].LoadBalancer; lb != nil && len(lb.Servers) != 0 {
		t.Fatalf("expected no servers updated when tunnel addresses empty, got: %+v", lb.Servers)
	}
}
