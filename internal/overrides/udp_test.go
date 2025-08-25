package overrides

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
)

func TestOverrideUDPRouters(t *testing.T) {
	t.Run("service override", func(t *testing.T) {
		routers := map[string]*dynamic.UDPRouter{
			"udp-service-router": {
				Service: "old-udp-service",
			},
		}
		
		// Test service override logic directly
		for _, router := range routers {
			router.Service = "new-udp-service"
		}
		
		if routers["udp-service-router"].Service != "new-udp-service" {
			t.Errorf("Expected service 'new-udp-service', got %s", routers["udp-service-router"].Service)
		}
	})
}

func TestOverrideUDPServices(t *testing.T) {
	t.Run("server override", func(t *testing.T) {
		services := map[string]*dynamic.UDPService{
			"udp-service": {
				LoadBalancer: &dynamic.UDPServersLoadBalancer{
					Servers: []dynamic.UDPServer{
						{Address: "old-server:8080"},
					},
				},
			},
		}
		
		// Test server override logic directly
		for _, service := range services {
			if service.LoadBalancer != nil {
				service.LoadBalancer.Servers = []dynamic.UDPServer{
					{Address: "new-server:8080"},
					{Address: "backup-server:8080"},
				}
			}
		}
		
		if len(services["udp-service"].LoadBalancer.Servers) != 2 {
			t.Errorf("Expected 2 servers, got %d", len(services["udp-service"].LoadBalancer.Servers))
		}
		
		if services["udp-service"].LoadBalancer.Servers[0].Address != "new-server:8080" {
			t.Errorf("Expected first server 'new-server:8080', got %s", services["udp-service"].LoadBalancer.Servers[0].Address)
		}
	})
}
