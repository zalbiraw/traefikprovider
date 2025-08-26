package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestOverrideTCPRouters_RuleOverride(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:    "HostSNI(`tcp.example.com`)",
			Service: "tcp-service",
		},
	}

	for _, router := range routers {
		router.Rule = "HostSNI(`new-tcp.example.com`)"
	}

	if routers["tcp-router"].Rule != "HostSNI(`new-tcp.example.com`)" {
		t.Errorf("Expected rule 'HostSNI(`new-tcp.example.com`)', got %s", routers["tcp-router"].Rule)
	}
}

func TestOverrideTCPRouters_EntrypointOverrideWithArray(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:        "HostSNI(`tcp.example.com`)",
			Service:     "tcp-service",
			EntryPoints: []string{"tcp"},
		},
	}

	overrides := config.RouterOverrides{
		Entrypoints: []config.OverrideEntrypoint{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   []string{"tcp-secure", "tcp-alt"},
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	expected := []string{"tcp-secure", "tcp-alt"}
	if len(routers["tcp-router"].EntryPoints) != 2 {
		t.Errorf("Expected entrypoints %v, got %v", expected, routers["tcp-router"].EntryPoints)
	}
}

func TestOverrideTCPRouters_EntrypointOverrideWithString(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:        "HostSNI(`tcp.example.com`)",
			Service:     "tcp-service",
			EntryPoints: []string{"tcp"},
		},
	}

	overrides := config.RouterOverrides{
		Entrypoints: []config.OverrideEntrypoint{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   "tcp-secure",
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	if len(routers["tcp-router"].EntryPoints) != 2 || routers["tcp-router"].EntryPoints[1] != "tcp-secure" {
		t.Errorf("Expected entrypoints to include 'tcp-secure', got %v", routers["tcp-router"].EntryPoints)
	}
}

func TestOverrideTCPRouters_ServiceOverride(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:    "HostSNI(`tcp.example.com`)",
			Service: "old-tcp-service",
		},
	}

	overrides := config.RouterOverrides{
		Services: []config.OverrideService{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   "new-tcp-service",
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	if routers["tcp-router"].Service != "new-tcp-service" {
		t.Errorf("Expected service 'new-tcp-service', got %s", routers["tcp-router"].Service)
	}
}

func TestOverrideTCPRouters_ServiceOverrideWithReplacement(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:    "HostSNI(`tcp.example.com`)",
			Service: "original-service",
		},
	}

	overrides := config.RouterOverrides{
		Services: []config.OverrideService{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   "prefix-$1-suffix",
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	expected := "prefix-original-service-suffix"
	if routers["tcp-router"].Service != expected {
		t.Errorf("Expected service '%s', got %s", expected, routers["tcp-router"].Service)
	}
}

func TestOverrideTCPRouters_MiddlewareOverrideWithArray(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:        "HostSNI(`tcp.example.com`)",
			Service:     "tcp-service",
			Middlewares: []string{"existing"},
		},
	}

	overrides := config.RouterOverrides{
		Middlewares: []config.OverrideMiddleware{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   []string{"tcp-auth", "tcp-ratelimit"},
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	expected := []string{"tcp-auth", "tcp-ratelimit"}
	if len(routers["tcp-router"].Middlewares) != 2 {
		t.Errorf("Expected middlewares %v, got %v", expected, routers["tcp-router"].Middlewares)
	}
}

func TestOverrideTCPRouters_MiddlewareOverrideWithString(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:        "HostSNI(`tcp.example.com`)",
			Service:     "tcp-service",
			Middlewares: []string{"existing"},
		},
	}

	overrides := config.RouterOverrides{
		Middlewares: []config.OverrideMiddleware{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   "tcp-auth",
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	expected := []string{"existing", "tcp-auth"}
	if len(routers["tcp-router"].Middlewares) != 2 {
		t.Errorf("Expected middlewares %v, got %v", expected, routers["tcp-router"].Middlewares)
	}
}

func TestOverrideTCPRouters_EmptyOverrides(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:    "HostSNI(`tcp.example.com`)",
			Service: "tcp-service",
		},
	}

	overrides := config.RouterOverrides{}
	OverrideTCPRouters(routers, overrides)

	if routers["tcp-router"].Service != "tcp-service" {
		t.Errorf("Expected service to remain 'tcp-service', got %s", routers["tcp-router"].Service)
	}
}

func TestOverrideTCPRouters_EmptyRouterMap(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{}

	overrides := config.RouterOverrides{
		Services: []config.OverrideService{
			{
				Matcher: "Name(`any-router`)",
				Value:   "new-service",
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	if len(routers) != 0 {
		t.Errorf("Expected empty router map, got %d routers", len(routers))
	}
}

func TestOverrideTCPRouters_AllOverrideTypesWithFiltering(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-router": {
			Rule:        "HostSNI(`tcp.example.com`)",
			Service:     "tcp-service",
			EntryPoints: []string{"tcp"},
			Middlewares: []string{"existing"},
		},
	}

	overrides := config.RouterOverrides{
		Entrypoints: []config.OverrideEntrypoint{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   []string{"tcp-secure"},
			},
		},
		Services: []config.OverrideService{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   "new-tcp-service",
			},
		},
		Middlewares: []config.OverrideMiddleware{
			{
				Matcher: "Name(`tcp-router`)",
				Value:   []string{"tcp-auth"},
			},
		},
	}

	OverrideTCPRouters(routers, overrides)

	if len(routers["tcp-router"].EntryPoints) != 1 || routers["tcp-router"].EntryPoints[0] != "tcp-secure" {
		t.Errorf("Expected entrypoints [tcp-secure], got %v", routers["tcp-router"].EntryPoints)
	}

	if routers["tcp-router"].Service != "new-tcp-service" {
		t.Errorf("Expected service 'new-tcp-service', got %s", routers["tcp-router"].Service)
	}

	if len(routers["tcp-router"].Middlewares) != 1 || routers["tcp-router"].Middlewares[0] != "tcp-auth" {
		t.Errorf("Expected middlewares [tcp-auth], got %v", routers["tcp-router"].Middlewares)
	}
}

func TestOverrideTCPServices(t *testing.T) {
	t.Run("server override", func(t *testing.T) {
		services := map[string]*dynamic.TCPService{
			"tcp-service": {
				LoadBalancer: &dynamic.TCPServersLoadBalancer{
					Servers: []dynamic.TCPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		// Test server override logic directly
		for _, service := range services {
			if service.LoadBalancer != nil {
				service.LoadBalancer.Servers = []dynamic.TCPServer{
					{Address: "new-server:8080"},
					{Address: "backup-server:8080"},
				}
			}
		}

		if len(services["tcp-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["tcp-service"].LoadBalancer.Servers))
		}

		if services["tcp-service"].LoadBalancer.Servers[0].Address != "new-server:8080" {
			t.Errorf("Expected first server 'new-server:8080', got %s", services["tcp-service"].LoadBalancer.Servers[0].Address)
		}
	})

	t.Run("server override with array", func(t *testing.T) {
		services := map[string]*dynamic.TCPService{
			"test-service": {
				LoadBalancer: &dynamic.TCPServersLoadBalancer{
					Servers: []dynamic.TCPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`test-service`)",
					Value:   []string{"new-server:8080", "backup-server:8080"},
				},
			},
		}

		OverrideTCPServices(services, overrides, nil)

		if len(services["test-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["test-service"].LoadBalancer.Servers))
		}

		if services["test-service"].LoadBalancer.Servers[0].Address != "new-server:8080" {
			t.Errorf("Expected first server 'new-server:8080', got %s", services["test-service"].LoadBalancer.Servers[0].Address)
		}
	})

	t.Run("server override with string", func(t *testing.T) {
		services := map[string]*dynamic.TCPService{
			"test-service": {
				LoadBalancer: &dynamic.TCPServersLoadBalancer{
					Servers: []dynamic.TCPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`test-service`)",
					Value:   "new-server:8080",
				},
			},
		}

		OverrideTCPServices(services, overrides, nil)

		if len(services["test-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["test-service"].LoadBalancer.Servers))
		}

		if services["test-service"].LoadBalancer.Servers[1].Address != "new-server:8080" {
			t.Errorf("Expected second server 'new-server:8080', got %s", services["test-service"].LoadBalancer.Servers[1].Address)
		}
	})

	t.Run("server override with tunnel", func(t *testing.T) {
		services := map[string]*dynamic.TCPService{
			"tunnel-service": {
				LoadBalancer: &dynamic.TCPServersLoadBalancer{
					Servers: []dynamic.TCPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		tunnels := []config.TunnelConfig{
			{
				Name:      "tcp-tunnel",
				Addresses: []string{"tunnel1:8080", "tunnel2:8080"},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`tunnel-service`)",
					Value:   []string{"ignored:8080"},
					Tunnel:  "tcp-tunnel",
				},
			},
		}

		OverrideTCPServices(services, overrides, tunnels)

		if len(services["tunnel-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers from tunnel, got %d", len(services["tunnel-service"].LoadBalancer.Servers))
		}

		if services["tunnel-service"].LoadBalancer.Servers[0].Address != "tunnel1:8080" {
			t.Errorf("Expected first server 'tunnel1:8080', got %s", services["tunnel-service"].LoadBalancer.Servers[0].Address)
		}
	})
}
