package matchers

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestMatchTCPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.TCPRouter
		pattern  string
		expected []string
	}{
		{
			name: "match all routers",
			routers: map[string]*dynamic.TCPRouter{
				"tcp-router": {
					Rule:    "HostSNI(`tcp.example.com`)",
					Service: "tcp-service",
				},
				"mysql-router": {
					Rule:    "HostSNI(`mysql.example.com`)",
					Service: "mysql-service",
				},
			},
			pattern:  ".*",
			expected: []string{"mysql-router", "tcp-router"},
		},
		{
			name: "match specific pattern",
			routers: map[string]*dynamic.TCPRouter{
				"tcp-router": {
					Rule:    "HostSNI(`tcp.example.com`)",
					Service: "tcp-service",
				},
				"mysql-router": {
					Rule:    "HostSNI(`mysql.example.com`)",
					Service: "mysql-service",
				},
			},
			pattern:  "tcp-.*",
			expected: []string{"tcp-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPRouters(tt.routers, &config.RoutersConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

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

func TestTCPRouters_EarlyReturnWhenNoMatcher(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"r1": {Service: "s1"},
	}
	// Empty provider and section matcher -> early return original map
	got := TCPRouters(routers, &config.RoutersConfig{Matcher: ""}, "")
	if len(got) != 1 || got["r1"] != routers["r1"] {
		t.Fatalf("expected original routers map to be returned unchanged")
	}
}

func TestTCPServices_CompileErrorSyntaxReturnsEmpty(t *testing.T) {
	svcs := map[string]*dynamic.TCPService{"s1": {}}
	got := TCPServices(svcs, &config.ServicesConfig{Matcher: "Name(`unterminated"}, "")
	if len(got) != 0 {
		t.Fatalf("expected empty result on compile syntax error, got %d", len(got))
	}
}

func TestTCPMiddlewares_CompileErrorSyntaxReturnsEmpty(t *testing.T) {
	mws := map[string]*dynamic.TCPMiddleware{"m1": {}}
	got := TCPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: "Name(`unterminated"}, "")
	if len(got) != 0 {
		t.Fatalf("expected empty result on compile syntax error, got %d", len(got))
	}
}

func TestTCPServices_EarlyReturnAndCompileError(t *testing.T) {
	svcs := map[string]*dynamic.TCPService{"s1": {LoadBalancer: &dynamic.TCPServersLoadBalancer{}}}
	// Early return when combined matcher is empty
	if got := TCPServices(svcs, &config.ServicesConfig{Matcher: ""}, ""); len(got) != 1 || got["s1"] != svcs["s1"] {
		t.Fatalf("expected original services map to be returned")
	}
	// Compile error path
	if got := TCPServices(svcs, &config.ServicesConfig{Matcher: "NameRegexp(`[` )"}, ""); len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestTCPMiddlewares_EarlyReturnAndCompileError(t *testing.T) {
	mws := map[string]*dynamic.TCPMiddleware{"m1": {}}
	// Early return when combined matcher is empty
	if got := TCPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: ""}, ""); len(got) != 1 || got["m1"] != mws["m1"] {
		t.Fatalf("expected original middlewares map to be returned")
	}
	// Compile error path
	if got := TCPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: "NameRegexp(`[` )"}, ""); len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestTCPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]*dynamic.TCPService
		pattern  string
		expected []string
	}{
		{
			name: "match all TCP services",
			services: map[string]*dynamic.TCPService{
				"tcp-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "tcp1:3306"},
						},
					},
				},
				"mysql-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "mysql1:3306"},
						},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"mysql-service", "tcp-service"},
		},
		{
			name: "match specific pattern",
			services: map[string]*dynamic.TCPService{
				"tcp-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "tcp1:3306"},
						},
					},
				},
				"mysql-service": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{Address: "mysql1:3306"},
						},
					},
				},
			},
			pattern:  "tcp-.*",
			expected: []string{"tcp-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPServices(tt.services, &config.ServicesConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

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

func TestTCPMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		middlewares map[string]*dynamic.TCPMiddleware
		pattern     string
		expected    []string
	}{
		{
			name: "match all TCP middlewares",
			middlewares: map[string]*dynamic.TCPMiddleware{
				"tcp-auth": {
					IPWhiteList: &dynamic.TCPIPWhiteList{
						SourceRange: []string{"192.168.1.0/24"},
					},
				},
				"tcp-limiter": {
					InFlightConn: &dynamic.TCPInFlightConn{
						Amount: 100,
					},
				},
			},
			pattern:  ".*",
			expected: []string{"tcp-auth", "tcp-limiter"},
		},
		{
			name: "match specific pattern",
			middlewares: map[string]*dynamic.TCPMiddleware{
				"tcp-auth": {
					IPWhiteList: &dynamic.TCPIPWhiteList{
						SourceRange: []string{"192.168.1.0/24"},
					},
				},
				"tcp-limiter": {
					InFlightConn: &dynamic.TCPInFlightConn{
						Amount: 100,
					},
				},
			},
			pattern:  "^tcp-auth$",
			expected: []string{"tcp-auth"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCPMiddlewares(tt.middlewares, &config.MiddlewaresConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d middlewares, got %d", len(tt.expected), len(result))
				return
			}

			for _, expectedName := range tt.expected {
				if _, found := result[expectedName]; !found {
					t.Errorf("Expected middleware %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestTCPRoutersAdvancedmatching(t *testing.T) {
	t.Run("match by entrypoints", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:        "HostSNI(`tcp1.example.com`)",
				Service:     "tcp-service-1",
				EntryPoints: []string{"tcp", "tcp-secure"},
			},
			"tcp-router-2": {
				Rule:        "HostSNI(`tcp2.example.com`)",
				Service:     "tcp-service-2",
				EntryPoints: []string{"tcp-alt"},
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "Entrypoint(`tcp`)",
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["tcp-router-1"]; !exists {
			t.Error("Expected tcp-router-1 to be in result")
		}
	})

	t.Run("match by rule", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:    "HostSNI(`tcp1.example.com`)",
				Service: "tcp-service-1",
			},
			"tcp-router-2": {
				Rule:    "HostSNI(`tcp2.example.com`)",
				Service: "tcp-service-2",
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "NameRegexp(`.*tcp-router-1.*`)",
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["tcp-router-1"]; !exists {
			t.Error("Expected tcp-router-1 to be in result")
		}
	})

	t.Run("match by service", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:    "HostSNI(`tcp1.example.com`)",
				Service: "tcp-service-1",
			},
			"tcp-router-2": {
				Rule:    "HostSNI(`tcp2.example.com`)",
				Service: "tcp-service-2",
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "ServiceRegexp(`.*tcp-service-1.*`)",
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["tcp-router-1"]; !exists {
			t.Error("Expected tcp-router-1 to be in result")
		}
	})

	t.Run("match by service name", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:    "HostSNI(`tcp1.example.com`)",
				Service: "tcp-service-1",
			},
			"tcp-router-2": {
				Rule:    "HostSNI(`tcp2.example.com`)",
				Service: "tcp-service-2",
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "Service(`tcp-service-1`)",
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 router, got %d", len(result))
		}

		if _, exists := result["tcp-router-1"]; !exists {
			t.Error("Expected tcp-router-1 to be in result")
		}
	})

	t.Run("invalid rule regex", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:    "HostSNI(`tcp1.example.com`)",
				Service: "tcp-service-1",
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "NameRegexp(`[`)", // Invalid regex
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 0 {
			t.Errorf("Expected 0 routers due to invalid regex, got %d", len(result))
		}
	})

	t.Run("invalid service regex", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router-1": {
				Rule:    "HostSNI(`tcp1.example.com`)",
				Service: "tcp-service-1",
			},
		}

		cfg := &config.RoutersConfig{
			Matcher: "ServiceRegexp(`*`)", // Invalid regex
		}

		result := TCPRouters(routers, cfg, "")

		if len(result) != 0 {
			t.Errorf("Expected 0 routers due to invalid regex, got %d", len(result))
		}
	})
}
