package internal

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	tlstypes "github.com/traefik/genconf/dynamic/tls"
)

func assertHTTPInitialized(t *testing.T, http *dynamic.HTTPConfiguration) {
	t.Helper()
	if http == nil {
		t.Fatal("Expected HTTP section to be initialized")
	}
	if http.Routers == nil || http.Services == nil || http.Middlewares == nil {
		t.Error("Expected HTTP maps to be initialized")
	}
}

func assertTCPInitialized(t *testing.T, tcp *dynamic.TCPConfiguration) {
	t.Helper()
	if tcp == nil {
		t.Fatal("Expected TCP section to be initialized")
	}
	if tcp.Routers == nil || tcp.Services == nil || tcp.Middlewares == nil {
		t.Error("Expected TCP maps to be initialized")
	}
}

func assertUDPInitialized(t *testing.T, udp *dynamic.UDPConfiguration) {
	t.Helper()
	if udp == nil {
		t.Fatal("Expected UDP section to be initialized")
	}
	if udp.Routers == nil || udp.Services == nil {
		t.Error("Expected UDP maps to be initialized")
	}
}

func assertTLSInitialized(t *testing.T, tls *dynamic.TLSConfiguration) {
	t.Helper()
	if tls == nil {
		t.Fatal("Expected TLS section to be initialized")
	}
	if tls.Certificates == nil || tls.Options == nil || tls.Stores == nil {
		t.Error("Expected TLS fields to be initialized")
	}
}

func assertConfigInitialized(t *testing.T, result *dynamic.Configuration) {
	t.Helper()
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	assertHTTPInitialized(t, result.HTTP)
	assertTCPInitialized(t, result.TCP)
	assertUDPInitialized(t, result.UDP)
	assertTLSInitialized(t, result.TLS)
}

func TestMergeConfigurations_Empty(t *testing.T) {
	result := MergeConfigurations()
	assertConfigInitialized(t, result)
}

func TestMergeConfigurations_Single(t *testing.T) {
	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router1": {Rule: "Host(`example.com`)", Service: "service1"},
			},
			Services: map[string]*dynamic.Service{
				"service1": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://backend:8080"}}}},
			},
		},
	}
	result := MergeConfigurations(cfg)
	assertConfigInitialized(t, result)
	if len(result.HTTP.Routers) != 1 || len(result.HTTP.Services) != 1 {
		t.Errorf("Expected 1 router and 1 service, got %d/%d", len(result.HTTP.Routers), len(result.HTTP.Services))
	}
}

func TestMergeConfigurations_MultipleHTTP(t *testing.T) {
	c1 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{
		Routers:  map[string]*dynamic.Router{"router1": {Rule: "Host(`api.example.com`)", Service: "api-service"}},
		Services: map[string]*dynamic.Service{"api-service": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://api-backend:8080"}}}}},
	}}
	c2 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{
		Routers:  map[string]*dynamic.Router{"router2": {Rule: "Host(`web.example.com`)", Service: "web-service"}},
		Services: map[string]*dynamic.Service{"web-service": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://web-backend:8080"}}}}},
	}}
	result := MergeConfigurations(c1, c2)
	assertConfigInitialized(t, result)
	if len(result.HTTP.Routers) != 2 || len(result.HTTP.Services) != 2 {
		t.Errorf("Expected 2 routers and 2 services, got %d/%d", len(result.HTTP.Routers), len(result.HTTP.Services))
	}
}

func TestMergeConfigurations_TCP(t *testing.T) {
	c1 := &dynamic.Configuration{TCP: &dynamic.TCPConfiguration{
		Routers:  map[string]*dynamic.TCPRouter{"tcp-router1": {Rule: "HostSNI(`tcp.example.com`)", Service: "tcp-service1"}},
		Services: map[string]*dynamic.TCPService{"tcp-service1": {LoadBalancer: &dynamic.TCPServersLoadBalancer{Servers: []dynamic.TCPServer{{Address: "tcp-backend:8081"}}}}},
	}}
	c2 := &dynamic.Configuration{TCP: &dynamic.TCPConfiguration{
		Routers:  map[string]*dynamic.TCPRouter{"tcp-router2": {Rule: "HostSNI(`tcp2.example.com`)", Service: "tcp-service2"}},
		Services: map[string]*dynamic.TCPService{"tcp-service2": {LoadBalancer: &dynamic.TCPServersLoadBalancer{Servers: []dynamic.TCPServer{{Address: "tcp2-backend:8081"}}}}},
	}}
	result := MergeConfigurations(c1, c2)
	assertConfigInitialized(t, result)
	if len(result.TCP.Routers) != 2 || len(result.TCP.Services) != 2 {
		t.Errorf("Expected 2 TCP routers and 2 services, got %d/%d", len(result.TCP.Routers), len(result.TCP.Services))
	}
}

func TestMergeConfigurations_UDP(t *testing.T) {
	c := &dynamic.Configuration{UDP: &dynamic.UDPConfiguration{
		Routers:  map[string]*dynamic.UDPRouter{"udp-router1": {Service: "udp-service1"}},
		Services: map[string]*dynamic.UDPService{"udp-service1": {LoadBalancer: &dynamic.UDPServersLoadBalancer{Servers: []dynamic.UDPServer{{Address: "udp-backend:8082"}}}}},
	}}
	result := MergeConfigurations(c)
	assertConfigInitialized(t, result)
	if len(result.UDP.Routers) != 1 || len(result.UDP.Services) != 1 {
		t.Errorf("Expected 1 UDP router and 1 service, got %d/%d", len(result.UDP.Routers), len(result.UDP.Services))
	}
}

func TestMergeConfigurations_TLS(t *testing.T) {
	c1 := &dynamic.Configuration{TLS: &dynamic.TLSConfiguration{
		Certificates: []*tlstypes.CertAndStores{{Certificate: tlstypes.Certificate{CertFile: "/path/to/cert1.pem", KeyFile: "/path/to/key1.pem"}}},
		Options:      map[string]tlstypes.Options{"default": {MinVersion: "VersionTLS12"}},
		Stores:       map[string]tlstypes.Store{"default": {DefaultCertificate: &tlstypes.Certificate{CertFile: "/path/to/default.pem", KeyFile: "/path/to/default-key.pem"}}},
	}}
	c2 := &dynamic.Configuration{TLS: &dynamic.TLSConfiguration{
		Certificates: []*tlstypes.CertAndStores{{Certificate: tlstypes.Certificate{CertFile: "/path/to/cert2.pem", KeyFile: "/path/to/key2.pem"}}},
		Options:      map[string]tlstypes.Options{"custom": {MinVersion: "VersionTLS13"}},
	}}
	result := MergeConfigurations(c1, c2)
	assertConfigInitialized(t, result)
	if len(result.TLS.Certificates) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(result.TLS.Certificates))
	}
}

func TestMergeConfigurations_NilSections(t *testing.T) {
	c1 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{Routers: map[string]*dynamic.Router{"router1": {Rule: "Host(`example.com`)", Service: "service1"}}}}
	c2 := &dynamic.Configuration{TCP: &dynamic.TCPConfiguration{Routers: map[string]*dynamic.TCPRouter{"tcp-router1": {Rule: "HostSNI(`tcp.example.com`)", Service: "tcp-service1"}}}}
	result := MergeConfigurations(c1, c2)
	assertConfigInitialized(t, result)
}

func TestMergeConfigurations_NilConfigurations(t *testing.T) {
	c1 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{Routers: map[string]*dynamic.Router{"router1": {Rule: "Host(`example.com`)", Service: "service1"}}}}
	c2 := (*dynamic.Configuration)(nil)
	c3 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{Routers: map[string]*dynamic.Router{"router2": {Rule: "Host(`example2.com`)", Service: "service2"}}}}
	result := MergeConfigurations(c1, c2, c3)
	assertConfigInitialized(t, result)
	if len(result.HTTP.Routers) != 2 {
		t.Errorf("Expected 2 HTTP routers, got %d", len(result.HTTP.Routers))
	}
}

func TestMergeConfigurations_OverlappingKeys(t *testing.T) {
	c1 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{
		Routers:  map[string]*dynamic.Router{"router1": {Rule: "Host(`example.com`)", Service: "service1"}},
		Services: map[string]*dynamic.Service{"shared-service": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://backend1:8080"}}}}},
	}}
	c2 := &dynamic.Configuration{HTTP: &dynamic.HTTPConfiguration{
		Routers:  map[string]*dynamic.Router{"router1": {Rule: "Host(`updated.example.com`)", Service: "updated-service"}},
		Services: map[string]*dynamic.Service{"shared-service": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://backend2:8080"}}}}},
	}}
	result := MergeConfigurations(c1, c2)
	assertConfigInitialized(t, result)
	if result.HTTP.Routers["router1"].Rule != "Host(`updated.example.com`)" {
		t.Error("Expected router1 rule to be overwritten")
	}
	if result.HTTP.Services["shared-service"].LoadBalancer.Servers[0].URL != "http://backend2:8080" {
		t.Error("Expected shared-service to be overwritten")
	}
}

func TestMergeConfigurations_AllProtocols(t *testing.T) {
	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     map[string]*dynamic.Router{"http-router": {Rule: "Host(`http.example.com`)", Service: "http-service"}},
			Services:    map[string]*dynamic.Service{"http-service": {LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "http://http-backend:8080"}}}}},
			Middlewares: map[string]*dynamic.Middleware{"auth": {BasicAuth: &dynamic.BasicAuth{Users: []string{"user:password"}}}},
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:     map[string]*dynamic.TCPRouter{"tcp-router": {Rule: "HostSNI(`tcp.example.com`)", Service: "tcp-service"}},
			Services:    map[string]*dynamic.TCPService{"tcp-service": {LoadBalancer: &dynamic.TCPServersLoadBalancer{Servers: []dynamic.TCPServer{{Address: "tcp-backend:8081"}}}}},
			Middlewares: map[string]*dynamic.TCPMiddleware{"tcp-auth": {IPWhiteList: &dynamic.TCPIPWhiteList{SourceRange: []string{"192.168.1.0/24"}}}},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  map[string]*dynamic.UDPRouter{"udp-router": {Service: "udp-service"}},
			Services: map[string]*dynamic.UDPService{"udp-service": {LoadBalancer: &dynamic.UDPServersLoadBalancer{Servers: []dynamic.UDPServer{{Address: "udp-backend:8082"}}}}},
		},
		TLS: &dynamic.TLSConfiguration{
			Certificates: []*tlstypes.CertAndStores{{Certificate: tlstypes.Certificate{CertFile: "/path/to/cert.pem", KeyFile: "/path/to/key.pem"}}},
			Options:      map[string]tlstypes.Options{"default": {MinVersion: "VersionTLS12"}},
			Stores:       map[string]tlstypes.Store{"default": {DefaultCertificate: &tlstypes.Certificate{CertFile: "/path/to/default.pem", KeyFile: "/path/to/default-key.pem"}}},
		},
	}
	result := MergeConfigurations(cfg)
	assertConfigInitialized(t, result)
}

func TestMergeConfigurationsOverwrite(t *testing.T) {
	// Test that later configurations overwrite earlier ones for the same keys
	config1 := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"shared-router": {
					Rule:    "Host(`old.example.com`)",
					Service: "old-service",
				},
			},
			Services: map[string]*dynamic.Service{
				"shared-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: "http://old-backend:8080"},
						},
					},
				},
			},
		},
	}

	config2 := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"shared-router": {
					Rule:    "Host(`new.example.com`)",
					Service: "new-service",
				},
			},
			Services: map[string]*dynamic.Service{
				"shared-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: "http://new-backend:8080"},
						},
					},
				},
			},
		},
	}

	result := MergeConfigurations(config1, config2)

	// Verify the second configuration overwrote the first
	if result.HTTP.Routers["shared-router"].Rule != "Host(`new.example.com`)" {
		t.Error("Expected router to be overwritten by later configuration")
	}
	if result.HTTP.Routers["shared-router"].Service != "new-service" {
		t.Error("Expected router service to be overwritten by later configuration")
	}

	if result.HTTP.Services["shared-service"].LoadBalancer.Servers[0].URL != "http://new-backend:8080" {
		t.Error("Expected service to be overwritten by later configuration")
	}
}

func TestMergeConfigurationsTLSCertificateAppend(t *testing.T) {
	// Test that TLS certificates are appended, not overwritten
	config1 := &dynamic.Configuration{
		TLS: &dynamic.TLSConfiguration{
			Certificates: []*tlstypes.CertAndStores{
				{
					Certificate: tlstypes.Certificate{
						CertFile: "/path/to/cert1.pem",
						KeyFile:  "/path/to/key1.pem",
					},
				},
			},
		},
	}

	config2 := &dynamic.Configuration{
		TLS: &dynamic.TLSConfiguration{
			Certificates: []*tlstypes.CertAndStores{
				{
					Certificate: tlstypes.Certificate{
						CertFile: "/path/to/cert2.pem",
						KeyFile:  "/path/to/key2.pem",
					},
				},
			},
		},
	}

	result := MergeConfigurations(config1, config2)

	// Verify certificates are appended
	if len(result.TLS.Certificates) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(result.TLS.Certificates))
	}

	if result.TLS.Certificates[0].Certificate.CertFile != "/path/to/cert1.pem" {
		t.Error("Expected first certificate to be preserved")
	}
	if result.TLS.Certificates[1].Certificate.CertFile != "/path/to/cert2.pem" {
		t.Error("Expected second certificate to be appended")
	}
}
