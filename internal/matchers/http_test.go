package matchers

import (
	"testing"

	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

func TestMatchHTTPRouters(t *testing.T) {
	tests := []struct {
		name     string
		routers  map[string]*dynamic.Router
		pattern  string
		expected []string
	}{
		{
			name: "match all routers",
			routers: map[string]*dynamic.Router{
				"web-router": {Rule: "Host(`web.example.com`)", Service: "web-service"},
				"api-router": {Rule: "Host(`api.example.com`)", Service: "api-service"},
			},
			pattern:  ".*",
			expected: []string{"api-router", "web-router"},
		},
		{
			name: "match specific pattern",
			routers: map[string]*dynamic.Router{
				"web-router": {Rule: "Host(`web.example.com`)", Service: "web-service"},
				"api-router": {Rule: "Host(`api.example.com`)", Service: "api-service"},
			},
			pattern:  "web-.*",
			expected: []string{"web-router"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPRouters(tt.routers, &config.RoutersConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

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

func TestHTTPMiddlewares_ProviderMatcher(t *testing.T) {
	mws := map[string]*dynamic.Middleware{
		"m1@p1": {},
		"m2@p2": {},
		"m3@p1": {},
	}
	got := HTTPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: ""}, "Provider(`p1`)")
	if len(got) != 2 {
		t.Fatalf("expected 2 middlewares from provider p1, got %d", len(got))
	}
	if _, ok := got["m1@p1"]; !ok {
		t.Fatalf("missing m1@p1")
	}
	if _, ok := got["m3@p1"]; !ok {
		t.Fatalf("missing m3@p1")
	}
	if _, ok := got["m2@p2"]; ok {
		t.Fatalf("unexpected m2@p2 included")
	}
}

func TestHTTPMiddlewares_ProviderMatcher_NoProviderSuffix(t *testing.T) {
	mws := map[string]*dynamic.Middleware{
		"m1": {},
		"m2": {},
	}
	got := HTTPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: ""}, "Provider(`p1`)")
	if len(got) != 0 {
		t.Fatalf("expected 0 middlewares when names lack provider suffix, got %d", len(got))
	}
}

func TestHTTPMiddlewares_NoMatchValidRule(t *testing.T) {
	m := map[string]*dynamic.Middleware{
		"auth": {},
	}
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Matcher: "Name(`does-not-exist`)"}, "")
	if len(out) != 0 {
		t.Fatalf("expected 0 middlewares for no-match valid rule, got %d", len(out))
	}
}

func TestHTTPMiddlewares_NoMatchInvalidRule(t *testing.T) {
	m := map[string]*dynamic.Middleware{
		"auth": {},
	}
	// Use an invalid expression to trigger compile error
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Matcher: "!"}, "")
	if len(out) != 0 {
		t.Fatalf("expected 0 middlewares for no-match invalid rule, got %d", len(out))
	}
}

func TestHTTPMiddlewares_CompileErrorReturnsEmpty(t *testing.T) {
	m := map[string]*dynamic.Middleware{"auth": {}}
	// Malformed expression (missing RPAREN) -> parser error -> compileRule error path
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Matcher: "NameRegexp(`abc`"}, "")
	if len(out) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(out))
	}
}

func TestHTTPMiddlewares_CompileError_FromProviderRule(t *testing.T) {
	m := map[string]*dynamic.Middleware{"auth": {}}
	// Malformed provider expression
	out := HTTPMiddlewares(m, &config.MiddlewaresConfig{Matcher: ""}, "Name(`abc`")
	if len(out) != 0 {
		t.Fatalf("expected empty result on provider compile error, got %d", len(out))
	}
}

func TestHTTPRouters_DiscoverPriorityTrueNoReset(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1": {Service: "s1", Priority: 7},
	}
	// DiscoverPriority true => do not reset, preserve pointer
	out := HTTPRouters(routers, &config.RoutersConfig{Matcher: "Name(`r1`)", DiscoverPriority: true}, "")
	if len(out) != 1 {
		t.Fatalf("expected 1 router, got %d", len(out))
	}
	if out["r1"].Priority != 7 {
		t.Fatalf("expected priority unchanged (7), got %d", out["r1"].Priority)
	}
	if out["r1"] != routers["r1"] { // pointer preserved
		t.Fatalf("expected original pointer to be preserved when DiscoverPriority is true")
	}
}

func TestHTTPRouters_EarlyReturnWhenNoMatcher(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1": {Service: "s1"},
	}
	// Empty provider and section matcher -> early return original map
	got := HTTPRouters(routers, &config.RoutersConfig{Matcher: ""}, "")
	if len(got) != 1 || got["r1"] != routers["r1"] {
		t.Fatalf("expected original routers map to be returned unchanged")
	}
}

func TestHTTPRouters_CompileErrorReturnsEmpty(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1": {Service: "s1"},
	}
	// Invalid regex -> compileRule error path
	got := HTTPRouters(routers, &config.RoutersConfig{Matcher: "NameRegexp(`[` )"}, "")
	if len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestHTTPRouters_DiscoverPriorityReset(t *testing.T) {
	routers := map[string]*dynamic.Router{
		"r1": {Service: "s1", Priority: 42},
	}
	got := HTTPRouters(routers, &config.RoutersConfig{Matcher: "NameRegexp(`r1`)", DiscoverPriority: false}, "")
	if len(got) != 1 {
		t.Fatalf("expected 1 router, got %d", len(got))
	}
	if got["r1"].Priority != 0 {
		t.Fatalf("expected priority reset to 0, got %d", got["r1"].Priority)
	}
}

func TestHTTPMiddlewares_EarlyReturnAndCompileError(t *testing.T) {
	mws := map[string]*dynamic.Middleware{"m1": {}}
	// Early return
	if got := HTTPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: ""}, ""); len(got) != 1 || got["m1"] != mws["m1"] {
		t.Fatalf("expected original middleware map to be returned")
	}

	// Compile error path
	if got := HTTPMiddlewares(mws, &config.MiddlewaresConfig{Matcher: "NameRegexp(`[` )"}, ""); len(got) != 0 {
		t.Fatalf("expected empty result on compile error, got %d", len(got))
	}
}

func TestMatchHTTPServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]*dynamic.Service
		pattern  string
		expected []string
	}{
		{
			name: "match all services",
			services: map[string]*dynamic.Service{
				"web-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://web:80"}},
					},
				},
				"api-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://api:80"}},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"api-service", "web-service"},
		},
		{
			name: "match specific pattern",
			services: map[string]*dynamic.Service{
				"web-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://web:80"}},
					},
				},
				"api-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{{URL: "http://api:80"}},
					},
				},
			},
			pattern:  "web-.*",
			expected: []string{"web-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPServices(tt.services, &config.ServicesConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

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

func TestMatchHTTPMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		middlewares map[string]*dynamic.Middleware
		pattern     string
		expected    []string
	}{
		{
			name: "match all middlewares",
			middlewares: map[string]*dynamic.Middleware{
				"auth-middleware": {
					BasicAuth: &dynamic.BasicAuth{
						Users: []string{"user:pass"},
					},
				},
				"cors-middleware": {
					Headers: &dynamic.Headers{
						AccessControlAllowOriginList: []string{"*"},
					},
				},
			},
			pattern:  ".*",
			expected: []string{"auth-middleware", "cors-middleware"},
		},
		{
			name: "match specific pattern",
			middlewares: map[string]*dynamic.Middleware{
				"auth-middleware": {
					BasicAuth: &dynamic.BasicAuth{
						Users: []string{"user:pass"},
					},
				},
				"cors-middleware": {
					Headers: &dynamic.Headers{
						AccessControlAllowOriginList: []string{"*"},
					},
				},
			},
			pattern:  "auth-.*",
			expected: []string{"auth-middleware"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPMiddlewares(tt.middlewares, &config.MiddlewaresConfig{Matcher: "NameRegexp(`" + tt.pattern + "`)"}, "")

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
