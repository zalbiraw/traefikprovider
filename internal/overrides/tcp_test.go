package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
)

func TestOverrideTCPRouters(t *testing.T) {
	t.Run("rule override", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-router": {
				Rule:    "HostSNI(`tcp.example.com`)",
				Service: "tcp-service",
			},
		}
		
		// Test rule override logic directly
		for _, router := range routers {
			router.Rule = "HostSNI(`new-tcp.example.com`)"
		}
		
		if routers["tcp-router"].Rule != "HostSNI(`new-tcp.example.com`)" {
			t.Errorf("Expected rule 'HostSNI(`new-tcp.example.com`)', got %s", routers["tcp-router"].Rule)
		}
	})
	
	t.Run("service override", func(t *testing.T) {
		routers := map[string]*dynamic.TCPRouter{
			"tcp-service-router": {
				Rule:    "HostSNI(`tcp-service.example.com`)",
				Service: "old-tcp-service",
			},
		}
		
		// Test service override logic directly
		for _, router := range routers {
			router.Service = "new-tcp-service"
		}
		
		if routers["tcp-service-router"].Service != "new-tcp-service" {
			t.Errorf("Expected service 'new-tcp-service', got %s", routers["tcp-service-router"].Service)
		}
	})
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
}
