package filters

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

func TestFilterMapByNameRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		pattern  string
		expected []string
	}{
		{
			name: "match all services",
			input: map[string]interface{}{
				"web-app":     struct{}{},
				"api-service": struct{}{},
				"database":    struct{}{},
			},
			pattern:  ".*",
			expected: []string{"api-service", "database", "web-app"},
		},
		{
			name: "match specific pattern",
			input: map[string]interface{}{
				"web-app":     struct{}{},
				"api-service": struct{}{},
				"database":    struct{}{},
			},
			pattern:  "web-.*",
			expected: []string{"web-app"},
		},
		{
			name: "exclude internal services",
			input: map[string]interface{}{
				"web-service":    struct{}{},
				"api-service":    struct{}{},
				"admin@internal": struct{}{},
				"debug@internal": struct{}{},
			},
			pattern:  ".*",
			expected: []string{"api-service", "web-service"},
		},
		{
			name: "no matches",
			input: map[string]interface{}{
				"web-app":     struct{}{},
				"api-service": struct{}{},
			},
			pattern:  "nonexistent",
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    map[string]interface{}{},
			pattern:  ".*",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterMapByNameRegex(tt.input, tt.pattern)

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
	input := map[string]interface{}{
		"service1": map[string]interface{}{"url": "http://example.com"},
	}

	// Invalid regex pattern should return empty result
	result := filterMapByNameRegex(input, "[")
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

func TestExtractRouterPriority(t *testing.T) {
	tests := []struct {
		name     string
		router   map[string]interface{}
		expected int
	}{
		{
			name:     "router with priority",
			router:   map[string]interface{}{"priority": 100.0},
			expected: 100,
		},
		{
			name:     "router without priority",
			router:   map[string]interface{}{"rule": "Host(`example.com`)"},
			expected: 0,
		},
		{
			name:     "router with invalid priority type",
			router:   map[string]interface{}{"priority": "high"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRouterPriority(tt.router, "test-router")
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
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
	routers := map[string]interface{}{
		"web-router": map[string]interface{}{
			"rule":        "Host(`example.com`)",
			"entryPoints": []string{"web", "websecure"},
			"service":     "my-service",
		},
		"api-router": map[string]interface{}{
			"rule":        "PathPrefix(`/api`)",
			"entryPoints": []string{"api"},
			"service":     "api-service",
		},
	}

	// Test entrypoint filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{
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
	routers := map[string]interface{}{
		"host-router": map[string]interface{}{
			"rule":    "Host(`example.com`)",
			"service": "host-service",
		},
		"path-router": map[string]interface{}{
			"rule":    "PathPrefix(`/api`)",
			"service": "path-service",
		},
	}

	// Test rule filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{
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
	routers := map[string]interface{}{
		"service1-router": map[string]interface{}{
			"rule":    "Host(`example.com`)",
			"service": "my-service",
		},
		"service2-router": map[string]interface{}{
			"rule":    "PathPrefix(`/api`)",
			"service": "api-service",
		},
	}

	// Test service filtering
	result := HTTPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{
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

func TestHTTPRoutersInvalidData(t *testing.T) {
	// Test with invalid router data
	routers := map[string]interface{}{
		"invalid-router": map[string]interface{}{
			"invalid": make(chan int), // This will cause JSON marshal error
		},
		"valid-router": map[string]interface{}{
			"rule":    "Host(`example.com`)",
			"service": "my-service",
		},
	}

	result := HTTPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{Name: ".*"},
	})

	// Should only get the valid router
	if len(result) != 1 {
		t.Errorf("Expected 1 valid router, got %d", len(result))
	}

	if _, found := result["valid-router"]; !found {
		t.Error("Expected valid-router to be found")
	}
}

func TestTCPRoutersWithEntrypoints(t *testing.T) {
	routers := map[string]interface{}{
		"tcp-web": map[string]interface{}{
			"rule":        "HostSNI(`example.com`)",
			"entryPoints": []string{"tcp-web", "tcp-secure"},
			"service":     "tcp-service",
		},
		"tcp-api": map[string]interface{}{
			"rule":        "HostSNI(`api.example.com`)",
			"entryPoints": []string{"tcp-api"},
			"service":     "tcp-api-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{
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
	routers := map[string]interface{}{
		"sni-router": map[string]interface{}{
			"rule":    "HostSNI(`example.com`)",
			"service": "sni-service",
		},
		"catch-all": map[string]interface{}{
			"rule":    "HostSNI(`*`)",
			"service": "catch-service",
		},
	}

	result := TCPRouters(routers, &config.RoutersConfig{
		Filters: config.RouterFilters{
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
	routers := map[string]interface{}{
		"dns-router": map[string]interface{}{
			"service": "dns-service",
		},
		"syslog-router": map[string]interface{}{
			"service": "syslog-service",
		},
	}

	result := UDPRouters(routers, &config.UDPRoutersConfig{
		Filters: config.UDPRouterFilters{
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

func TestInvalidInputTypes(t *testing.T) {
	// Test all functions with invalid input types

	// HTTP functions
	httpResult := HTTPRouters("invalid", &config.RoutersConfig{})
	if len(httpResult) != 0 {
		t.Error("Expected empty result for invalid HTTP routers input")
	}

	httpServicesResult := HTTPServices("invalid", &config.ServicesConfig{})
	if len(httpServicesResult) != 0 {
		t.Error("Expected empty result for invalid HTTP services input")
	}

	httpMiddlewaresResult := HTTPMiddlewares("invalid", &config.MiddlewaresConfig{})
	if len(httpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for invalid HTTP middlewares input")
	}

	httpTransportsResult := HTTPServerTransports("invalid", &config.ServerTransportsConfig{})
	if len(httpTransportsResult) != 0 {
		t.Error("Expected empty result for invalid HTTP transports input")
	}

	// TCP functions
	tcpResult := TCPRouters("invalid", &config.RoutersConfig{})
	if len(tcpResult) != 0 {
		t.Error("Expected empty result for invalid TCP routers input")
	}

	tcpServicesResult := TCPServices("invalid", &config.ServicesConfig{})
	if len(tcpServicesResult) != 0 {
		t.Error("Expected empty result for invalid TCP services input")
	}

	tcpMiddlewaresResult := TCPMiddlewares("invalid", &config.MiddlewaresConfig{})
	if len(tcpMiddlewaresResult) != 0 {
		t.Error("Expected empty result for invalid TCP middlewares input")
	}

	// UDP functions
	udpResult := UDPRouters("invalid", &config.UDPRoutersConfig{})
	if len(udpResult) != 0 {
		t.Error("Expected empty result for invalid UDP routers input")
	}

	udpServicesResult := UDPServices("invalid", &config.UDPServicesConfig{})
	if len(udpServicesResult) != 0 {
		t.Error("Expected empty result for invalid UDP services input")
	}

}
