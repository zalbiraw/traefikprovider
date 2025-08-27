package matchers

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestUDPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.UDPRouter
		pattern  string
		expected []string
	}{
		{
			name: "match all UDP routers",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
				"admin": {
					Service: "admin",
				},
			},
			pattern:  ".*",
			expected: []string{"admin", "dns-router", "udp-router"},
		},
		{
			name: "match specific pattern",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-router"},
		},
		{
			name: "match specific pattern",
			routers: map[string]*dynamic.UDPRouter{
				"udp-router": {
					Service: "udp-service",
				},
				"dns-router": {
					Service: "dns-service",
				},
			},
			pattern:  "dns-.*",
			expected: []string{"dns-router"},
		},
		{
			name: "no-match valid rule",
			routers: map[string]*dynamic.UDPRouter{
				"r1": {Service: "s1"},
			},
			pattern:  "does-not-exist",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPRouters(tt.routers, &config.UDPRoutersConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d routers, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected router %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestUDPRouters_CompileError_FromProviderRule(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{"r1": {Service: "s1"}}
	// Invalid provider rule (missing RPAREN) should trigger compile error
	got := UDPRouters(routers, &config.UDPRoutersConfig{Matcher: ""}, "Name(`x`")
	if len(got) != 0 {
		t.Fatalf("expected empty result on provider compile error, got %d", len(got))
	}
}

func TestUDPRouters_CompileError_InvalidToken(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{"r1": {Service: "s1"}}
	// Lone '!' is invalid expression -> compile error
	got := UDPRouters(routers, &config.UDPRoutersConfig{Matcher: "!"}, "")
	if len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestUDPServices_ProviderMatcher(t *testing.T) {
	svcs := map[string]*dynamic.UDPService{
		"s1@p1": {LoadBalancer: &dynamic.UDPServersLoadBalancer{}},
		"s2@p2": {LoadBalancer: &dynamic.UDPServersLoadBalancer{}},
		"s3@p1": {LoadBalancer: &dynamic.UDPServersLoadBalancer{}},
	}
	cfg := &config.UDPServicesConfig{Matcher: ""}
	got := UDPServices(svcs, cfg, "Provider(`p1`)")
	if len(got) != 2 {
		t.Fatalf("expected 2 services from provider p1, got %d", len(got))
	}
	if _, ok := got["s1@p1"]; !ok {
		t.Fatalf("missing s1@p1")
	}
	if _, ok := got["s3@p1"]; !ok {
		t.Fatalf("missing s3@p1")
	}
	if _, ok := got["s2@p2"]; ok {
		t.Fatalf("unexpected s2@p2 included")
	}
}

func TestUDPRouters_ProviderMatcher(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"r1@p1": {Service: "s1"},
		"r2@p2": {Service: "s2"},
		"r3@p1": {Service: "s3"},
	}
	// provider-level matcher only
	cfg := &config.UDPRoutersConfig{Matcher: ""}
	got := UDPRouters(routers, cfg, "Provider(`p1`)")
	if len(got) != 2 {
		t.Fatalf("expected 2 routers from provider p1, got %d", len(got))
	}
	if _, ok := got["r1@p1"]; !ok {
		t.Fatalf("missing r1@p1")
	}
	if _, ok := got["r3@p1"]; !ok {
		t.Fatalf("missing r3@p1")
	}
	if _, ok := got["r2@p2"]; ok {
		t.Fatalf("unexpected r2@p2 included")
	}
}

func TestUDPRouters_CompileErrorReturnsEmpty(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"r1": {Service: "s1"},
	}
	// Invalid rule syntax -> compileRule error path
	got := UDPRouters(routers, &config.UDPRoutersConfig{Matcher: "NameRegexp(`abc`"}, "")
	if len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestUDPRouters_EarlyReturnWhenNoMatcher(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"r1": {Service: "s1"},
	}
	// Empty provider and section matcher -> early return original map
	got := UDPRouters(routers, &config.UDPRoutersConfig{Matcher: ""}, "")
	if len(got) != 1 || got["r1"] != routers["r1"] {
		t.Fatalf("expected original routers map to be returned unchanged")
	}
}

func TestUDPServices_EarlyReturnAndCompileError(t *testing.T) {
	svcs := map[string]*dynamic.UDPService{"s1": {LoadBalancer: &dynamic.UDPServersLoadBalancer{}}}
	// Early return when combined matcher is empty
	if got := UDPServices(svcs, &config.UDPServicesConfig{Matcher: ""}, ""); len(got) != 1 || got["s1"] != svcs["s1"] {
		t.Fatalf("expected original services map to be returned")
	}
	// Compile error path (invalid rule syntax)
	if got := UDPServices(svcs, &config.UDPServicesConfig{Matcher: "NameRegexp(`abc`"}, ""); len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestUDPServices_NoMatchValidRule(t *testing.T) {
	svcs := map[string]*dynamic.UDPService{"s1": {LoadBalancer: &dynamic.UDPServersLoadBalancer{}}}
	// valid rule that matches nothing
	out := UDPServices(svcs, &config.UDPServicesConfig{Matcher: "Name(`does-not-exist`)"}, "")
	if len(out) != 0 {
		t.Fatalf("expected 0 services for no-match valid rule, got %d", len(out))
	}
}

func TestUDPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]*dynamic.UDPService
		pattern  string
		expected []string
	}{
		{
			name: "match all UDP services",
			services: map[string]*dynamic.UDPService{
				"udp-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "udp1:53"},
						},
					},
				},
				"dns-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "dns1:53"},
						},
					},
				},
				"admin": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "admin:53"},
						},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"admin", "dns-service", "udp-service"},
		},
		{
			name: "match specific pattern",
			services: map[string]*dynamic.UDPService{
				"udp-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "udp1:53"},
						},
					},
				},
				"dns-service": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{Address: "dns1:53"},
						},
					},
				},
			},
			pattern:  "udp-.*",
			expected: []string{"udp-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UDPServices(tt.services, &config.UDPServicesConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d services, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected service %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestUDPRoutersAdvancedmatching(t *testing.T) {
	t.Run("match by entrypoints", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service:     "udp-service-1",
				EntryPoints: []string{"udp", "udp-secure"},
			},
			"udp-router-2": {
				Service:     "udp-service-2",
				EntryPoints: []string{"udp-alt"},
			},
		}

		cfg := &config.UDPRoutersConfig{
			Matcher: "Entrypoint(`udp`)",
		}

		result := UDPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["udp-router-1"]; !exists {
			t.Error("Expected udp-router-1 to be in result")
		}
	})

	t.Run("match by service", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service: "udp-service-1",
			},
			"udp-router-2": {
				Service: "udp-service-2",
			},
		}

		cfg := &config.UDPRoutersConfig{
			Matcher: "Service(`udp-service-1`)",
		}

		result := UDPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["udp-router-1"]; !exists {
			t.Error("Expected udp-router-1 to be in result")
		}
	})

	t.Run("invalid service regex", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router-1": {
				Service: "udp-service-1",
			},
		}

		cfg := &config.UDPRoutersConfig{
			Matcher: "ServiceRegexp(`*`)", // Invalid regex
		}

		result := UDPRouters(routers, cfg, "")

		if len(result) != 0 {
			t.Errorf("Expected 0 routers due to invalid regex, got %d", len(result))
		}
	})
}
