package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterMapByNameRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]*struct{}
		pattern  string
		expected []string
	}{
		{
			name: "match all services",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
				"database":    {},
			},
			pattern:  ".*",
			expected: []string{"api-service", "database", "web-app"},
		},
		{
			name: "match specific pattern",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
				"database":    {},
			},
			pattern:  "web-.*",
			expected: []string{"web-app"},
		},
		{
			name: "include all services",
			input: map[string]*struct{}{
				"web-service":    {},
				"api-service":    {},
				"admin@internal": {},
				"debug@internal": {},
			},
			pattern:  ".*",
			expected: []string{"admin@internal", "api-service", "debug@internal", "web-service"},
		},
		{
			name: "no matches",
			input: map[string]*struct{}{
				"web-app":     {},
				"api-service": {},
			},
			pattern:  "nonexistent",
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    map[string]*struct{}{},
			pattern:  ".*",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterMapByNameRegex[struct{}, *struct{}](tt.input, tt.pattern, "")

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d results, got %d", len(tt.expected), len(result))
				return
			}

			// Convert result to slice for comparison
			resultSlice := make([]string, 0, len(result))
			for key := range result {
				resultSlice = append(resultSlice, key)
			}

			// Sort both slices for comparison
			for _, expected := range tt.expected {
				found := false
				for _, actual := range resultSlice {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected key %s not found in result", expected)
				}
			}
		})
	}
}

func TestFilterMapByNameRegexInvalidPattern(t *testing.T) {
	input := map[string]*struct{}{
		"service1": {},
	}

	// Invalid regex pattern should return empty result
	result := filterMapByNameRegex[struct{}, *struct{}](input, "[", "")
	if len(result) != 0 {
		t.Errorf("Expected empty result for invalid regex, got %d items", len(result))
	}
}

func TestRegexMatchInvalidPattern(t *testing.T) {
	// Test invalid regex pattern
	matched, err := regexMatch("[", "test")
	if err == nil {
		t.Error("Expected error for invalid regex pattern")
	}
	if matched {
		t.Error("Expected no match for invalid regex")
	}
}

func TestRegexMatchEmptyPattern(t *testing.T) {
	// Test empty pattern should return true
	matched, err := regexMatch("", "test")
	if err != nil {
		t.Errorf("Expected no error for empty pattern, got: %v", err)
	}
	if !matched {
		t.Error("Expected match for empty pattern")
	}
}

func TestRouterEntrypointsMatch(t *testing.T) {
	tests := []struct {
		name      string
		routerEPs []string
		filterEPs []string
		expected  bool
	}{
		{
			name:      "empty filter entrypoints",
			routerEPs: []string{"web", "websecure"},
			filterEPs: []string{},
			expected:  true,
		},
		{
			name:      "router contains all filter entrypoints",
			routerEPs: []string{"web", "websecure", "api"},
			filterEPs: []string{"web", "websecure"},
			expected:  true,
		},
		{
			name:      "router missing some filter entrypoints",
			routerEPs: []string{"web"},
			filterEPs: []string{"web", "websecure"},
			expected:  false,
		},
		{
			name:      "empty router entrypoints",
			routerEPs: []string{},
			filterEPs: []string{"web"},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := routerEntrypointsMatch(tt.routerEPs, tt.filterEPs)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnmarshalRouter(t *testing.T) {
	tests := []struct {
		name        string
		routerMap   map[string]interface{}
		expectError bool
	}{
		{
			name: "valid router",
			routerMap: map[string]interface{}{
				"rule":        "Host(`example.com`)",
				"entryPoints": []string{"web"},
				"service":     "my-service",
			},
			expectError: false,
		},
		{
			name:        "invalid router data",
			routerMap:   map[string]interface{}{"invalid": make(chan int)},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var router dynamic.Router
			err := unmarshalRouter(tt.routerMap, &router)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHTTPRoutersWithEntrypoints(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"web-router": {
			Rule:        "Host(`example.com`)",
			EntryPoints: []string{"web"},
			Service:     "my-service",
		},
		"api-router": {
			Rule:        "PathPrefix(`/api`)",
			EntryPoints: []string{"api"},
			Service:     "api-service",
		},
	}

	// Test entrypoint filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:        ".*",
			Entrypoints: []string{"web"},
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with web entrypoint, got %d", len(result))
	}

	if _, found := result["web-router"]; !found {
		t.Error("Expected web-router to be found")
	}
}

func TestHTTPRoutersWithRuleFilter(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"host-router": {
			Rule:    "Host(`example.com`)",
			Service: "host-service",
		},
		"path-router": {
			Rule:    "PathPrefix(`/api`)",
			Service: "path-service",
		},
	}

	// Test rule filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name: ".*",
			Rule: "Host\\(.*\\)",
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with Host rule, got %d", len(result))
	}

	if _, found := result["host-router"]; !found {
		t.Error("Expected host-router to be found")
	}
}

func TestHTTPRoutersWithServiceFilter(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"service1-router": {
			Rule:    "Host(`example.com`)",
			Service: "my-service",
		},
		"service2-router": {
			Rule:    "PathPrefix(`/api`)",
			Service: "api-service",
		},
	}

	// Test service filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:    ".*",
			Service: "^my-.*",
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 router with my-service, got %d", len(result))
	}

	if _, found := result["service1-router"]; !found {
		t.Error("Expected service1-router to be found")
	}
}

func TestHTTPRoutersWithInvalidData(t *testing.T) {
	// Test with valid router data
	routers := map[string]*dynamic.Router{
		"valid-router": {
			Rule:    "Host(`example.com`)",
			Service: "my-service",
		},
	}

	result := HTTPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{Name: ".*"},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 valid router, got %d", len(result))
	}

	if _, found := result["valid-router"]; !found {
		t.Error("Expected valid-router to be found")
	}
}

func TestTCPRoutersWithEntrypoints(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"tcp-web": {
			Rule:        "HostSNI(`example.com`)",
			EntryPoints: []string{"tcp-web", "tcp-secure"},
			Service:     "tcp-service",
		},
		"tcp-api": {
			Rule:        "HostSNI(`api.example.com`)",
			EntryPoints: []string{"tcp-api"},
			Service:     "tcp-api-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name:        ".*",
			Entrypoints: []string{"tcp-web"},
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 TCP router with tcp-web entrypoint, got %d", len(result))
	}

	if _, found := result["tcp-web"]; !found {
		t.Error("Expected tcp-web router to be found")
	}
}

func TestTCPRoutersWithRuleFilter(t *testing.T) {
	routers := map[string]*dynamic.TCPRouter{
		"sni-router": {
			Rule:    "HostSNI(`example.com`)",
			Service: "sni-service",
		},
		"catch-all": {
			Rule:    "HostSNI(`*`)",
			Service: "catch-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filter: config.RouterFilter{
			Name: ".*",
			Rule: "HostSNI\\(`\\*`\\)",
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 TCP router with catch-all rule, got %d", len(result))
	}

	if _, found := result["catch-all"]; !found {
		t.Error("Expected catch-all router to be found")
	}
}

func TestUDPRoutersWithServiceFilter(t *testing.T) {
	routers := map[string]*dynamic.UDPRouter{
		"dns-router": {
			Service: "dns-service",
		},
		"syslog-router": {
			Service: "syslog-service",
		},
	}

	result := UDPRouters(routers, &config.UDPRoutersConfig{
		Filter: config.UDPRouterFilter{
			Name:    ".*",
			Service: "^dns-.*",
		},
	})

	if len(result) != 1 {
		t.Errorf("Expected 1 UDP router with dns service, got %d", len(result))
	}

	if _, found := result["dns-router"]; !found {
		t.Error("Expected dns-router to be found")
	}
}

func TestInvalidInputHandling(t *testing.T) {
	// Test with empty maps - the functions now expect typed maps, not interface{}

	// HTTP functions
	httpResult := HTTPRouters(map[string]*dynamic.Router{}, &config.RoutersConfig{})
	if len(httpResult) != 0 {
		t.Error("Expected empty result for empty HTTP routers input")
	}

	httpServicesResult := HTTPServices(map[string]*dynamic.Service{}, &config.ServicesConfig{})
	if len(httpServicesResult) != 0 {
		t.Error("Expected empty result for empty HTTP services input")
	}

	httpMiddlewaresResult := HTTPMiddlewares(map[string]*dynamic.Middleware{}, &config.MiddlewaresConfig{})
	if len(httpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty HTTP middlewares input")
	}

	// TCP functions
	tcpResult := TCPRouters(map[string]*dynamic.TCPRouter{}, &config.RoutersConfig{})
	if len(tcpResult) != 0 {
		t.Error("Expected empty result for empty TCP routers input")
	}

	tcpServicesResult := TCPServices(map[string]*dynamic.TCPService{}, &config.ServicesConfig{})
	if len(tcpServicesResult) != 0 {
		t.Error("Expected empty result for empty TCP services input")
	}

	tcpMiddlewaresResult := TCPMiddlewares(map[string]*dynamic.TCPMiddleware{}, &config.MiddlewaresConfig{})
	if len(tcpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for empty TCP middlewares input")
	}

	// UDP functions
	udpResult := UDPRouters(map[string]*dynamic.UDPRouter{}, &config.UDPRoutersConfig{})
	if len(udpResult) != 0 {
		t.Error("Expected empty result for empty UDP routers input")
	}

	udpServicesResult := UDPServices(map[string]*dynamic.UDPService{}, &config.UDPServicesConfig{})
	if len(udpServicesResult) != 0 {
		t.Error("Expected empty result for empty UDP services input")
	}
}
