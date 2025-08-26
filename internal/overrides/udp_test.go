package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestOverrideUDPRouters(t *testing.T) {
	t.Run("entrypoint override with array", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router": {
				Service:     "udp-service",
				EntryPoints: []string{"udp"},
			},
		}

		overrides := config.UDPOverrides{
			Entrypoints: []config.OverrideEntrypoint{
				{
					Matcher: "Name(`udp-router`)",
					Value:   []string{"udp-secure", "udp-alt"},
				},
			},
		}

		OverrideUDPRouters(routers, overrides)

		expected := []string{"udp-secure", "udp-alt"}
		if len(routers["udp-router"].EntryPoints) != 2 {
			t.Errorf("Expected entrypoints %v, got %v", expected, routers["udp-router"].EntryPoints)
		}
	})

	t.Run("entrypoint override with string", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router": {
				Service:     "udp-service",
				EntryPoints: []string{"udp"},
			},
		}

		overrides := config.UDPOverrides{
			Entrypoints: []config.OverrideEntrypoint{
				{
					Matcher: "Name(`udp-router`)",
					Value:   "udp-secure",
				},
			},
		}

		OverrideUDPRouters(routers, overrides)

		if len(routers["udp-router"].EntryPoints) != 2 || routers["udp-router"].EntryPoints[1] != "udp-secure" {
			t.Errorf("Expected entrypoints to include 'udp-secure', got %v", routers["udp-router"].EntryPoints)
		}
	})

	t.Run("service override", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router": {
				Service: "old-udp-service",
			},
		}

		overrides := config.UDPOverrides{
			Services: []config.OverrideService{
				{
					Matcher: "Name(`udp-router`)",
					Value:   "new-udp-service",
				},
			},
		}

		OverrideUDPRouters(routers, overrides)

		if routers["udp-router"].Service != "new-udp-service" {
			t.Errorf("Expected service 'new-udp-service', got %s", routers["udp-router"].Service)
		}
	})

	t.Run("service override with $1 replacement", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-router": {
				Service: "original-service",
			},
		}

		overrides := config.UDPOverrides{
			Services: []config.OverrideService{
				{
					Matcher: "Name(`udp-router`)",
					Value:   "prefix-$1-suffix",
				},
			},
		}

		OverrideUDPRouters(routers, overrides)

		expected := "prefix-original-service-suffix"
		if routers["udp-router"].Service != expected {
			t.Errorf("Expected service '%s', got %s", expected, routers["udp-router"].Service)
		}
	})
}

func TestOverrideUDPServices(t *testing.T) {
	t.Run("server override with array", func(t *testing.T) {
		services := map[string]*dynamic.UDPService{
			"udp-service": {
				LoadBalancer: &dynamic.UDPServersLoadBalancer{
					Servers: []dynamic.UDPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`udp-service`)",
					Value:   []string{"new-server:8080", "backup-server:8080"},
				},
			},
		}

		OverrideUDPServices(services, overrides)

		if len(services["udp-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["udp-service"].LoadBalancer.Servers))
		}

		if services["udp-service"].LoadBalancer.Servers[0].Address != "new-server:8080" {
			t.Errorf("Expected first server 'new-server:8080', got %s", services["udp-service"].LoadBalancer.Servers[0].Address)
		}
	})

	t.Run("server override with string", func(t *testing.T) {
		services := map[string]*dynamic.UDPService{
			"udp-service": {
				LoadBalancer: &dynamic.UDPServersLoadBalancer{
					Servers: []dynamic.UDPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`udp-service`)",
					Value:   "new-server:8080",
				},
			},
		}

		OverrideUDPServices(services, overrides)

		if len(services["udp-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["udp-service"].LoadBalancer.Servers))
		}

		if services["udp-service"].LoadBalancer.Servers[1].Address != "new-server:8080" {
			t.Errorf("Expected second server 'new-server:8080', got %s", services["udp-service"].LoadBalancer.Servers[1].Address)
		}
	})

	t.Run("server override with tunnel", func(t *testing.T) {
		services := map[string]*dynamic.UDPService{
			"tunnel-service": {
				LoadBalancer: &dynamic.UDPServersLoadBalancer{
					Servers: []dynamic.UDPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}

		tunnels := []config.TunnelConfig{
			{
				Name:      "udp-tunnel",
				Addresses: []string{"tunnel1:8080", "tunnel2:8080"},
			},
		}

		overrides := config.ServiceOverrides{
			Servers: []config.OverrideServer{
				{
					Matcher: "Name(`tunnel-service`)",
					Value:   []string{"ignored:8080"},
					Tunnel:  "udp-tunnel",
				},
			},
		}

		OverrideUDPServices(services, overrides)

		// Apply tunnels in a separate pass per new design
		ApplyTunnels(services, overrides, tunnels)

		if len(services["tunnel-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers from tunnel, got %d", len(services["tunnel-service"].LoadBalancer.Servers))
		}

		if services["tunnel-service"].LoadBalancer.Servers[0].Address != "tunnel1:8080" {
			t.Errorf("Expected first server 'tunnel1:8080', got %s", services["tunnel-service"].LoadBalancer.Servers[0].Address)
		}
	})
}
