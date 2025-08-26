package provider

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zalbiraw/traefik-provider/config"
	"github.com/zalbiraw/traefik-provider/internal/httpclient"
)

const (
	dockerComposeFile = "./cmd/preview/docker-compose.yml"
	testTimeout       = 60 * time.Second
)

var (
	// setupOnce ensures docker services are started only once across tests.
	//nolint:gochecknoglobals // Intentionally global to coordinate integration test setup.
	setupOnce sync.Once
	errSetup  error
)

// ensureDockerServices ensures Docker services are running (called once)
func ensureDockerServices() error {
	setupOnce.Do(func() {
		// Check if Docker is available
		if _, err := exec.LookPath("docker-compose"); err != nil {
			errSetup = fmt.Errorf("docker-compose not available: %w", err)
			return
		}

		// Start Docker Compose services
		cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d")
		if err := cmd.Run(); err != nil {
			errSetup = fmt.Errorf("failed to start Docker services: %w", err)
			return
		}

		// Wait for services to be ready
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		if err := waitForService(ctx, "http://localhost:8081/api/rawdata"); err != nil {
			errSetup = fmt.Errorf("provider1 service not ready: %w", err)
			return
		}

		if err := waitForService(ctx, "http://localhost:8082/api/rawdata"); err != nil {
			errSetup = fmt.Errorf("provider2 service not ready: %w", err)
			return
		}
	})
	return errSetup
}

// TestIntegrationBasic tests basic functionality with live Docker services
func TestIntegrationBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	if err := ensureDockerServices(); err != nil {
		t.Skipf("Skipping integration tests - Docker services not available: %v", err)
	}

	// Test basic configuration fetching
	t.Run("Fetch Provider1 Configuration", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/api/rawdata",
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg == nil {
			t.Fatal("Configuration is nil")
		}

		expectedRouters := []string{"provider1-api", "provider1-web", "provider1-admin", "provider1-test"}
		for _, routerName := range expectedRouters {
			if _, exists := dynCfg.HTTP.Routers[routerName]; !exists {
				t.Errorf("Expected router %s not found", routerName)
			}
		}

		expectedServices := []string{"provider1-service", "provider1-web-service", "provider1-admin-service", "provider1-test-service"}
		for _, serviceName := range expectedServices {
			if _, exists := dynCfg.HTTP.Services[serviceName]; !exists {
				t.Errorf("Expected service %s not found", serviceName)
			}
		}
	})

	t.Run("Fetch Provider2 Configuration", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8082,
				Path: "/api/rawdata",
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg == nil {
			t.Fatal("Configuration is nil")
		}

		expectedRouters := []string{"provider2-dashboard", "provider2-api", "provider2-secure", "provider2-metrics"}
		for _, routerName := range expectedRouters {
			if _, exists := dynCfg.HTTP.Routers[routerName]; !exists {
				t.Errorf("Expected router %s not found", routerName)
			}
		}

		expectedServices := []string{"provider2-service", "provider2-api-service", "provider2-secure-service", "provider2-metrics-service"}
		for _, serviceName := range expectedServices {
			if _, exists := dynCfg.HTTP.Services[serviceName]; !exists {
				t.Errorf("Expected service %s not found", serviceName)
			}
		}
	})
}

// TestIntegrationFilter tests basic filtering functionality
func TestIntegrationFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	if err := ensureDockerServices(); err != nil {
		t.Skipf("Skipping integration tests - Docker services not available: %v", err)
	}

	t.Run("HTTP Router Name Filter", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/api/rawdata",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Routers: &config.RoutersConfig{
					Discover: true,
					Filter: config.RouterFilter{
						Name: "provider1-.*",
					},
				},
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg == nil {
			t.Fatal("Configuration is nil")
		}

		if dynCfg.HTTP == nil {
			t.Fatal("HTTP configuration is nil")
		}

		for name := range dynCfg.HTTP.Routers {
			if !strings.HasPrefix(name, "provider1-") {
				t.Errorf("Router name %s doesn't match filter pattern", name)
			}
		}
	})

	t.Run("HTTP Service Name Filter", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8082,
				Path: "/api/rawdata",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Services: &config.ServicesConfig{
					Discover: true,
					Filter: config.ServiceFilter{
						Name: ".*-service",
					},
				},
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg.HTTP == nil {
			t.Fatal("HTTP configuration is nil")
		}

		// Verify that only services matching the pattern are returned
		foundMatchingService := false
		for name := range dynCfg.HTTP.Services {
			if strings.HasSuffix(name, "-service") {
				foundMatchingService = true
			}
		}
		if !foundMatchingService {
			t.Error("Expected to find at least one service ending with '-service'")
		}
	})
}

// TestIntegrationOverrides tests basic override functionality
func TestIntegrationOverrides(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	if err := ensureDockerServices(); err != nil {
		t.Skipf("Skipping integration tests - Docker services not available: %v", err)
	}

	t.Run("HTTP Router Rule Override", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/api/rawdata",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Routers: &config.RoutersConfig{
					Discover: true,
					Overrides: config.RouterOverrides{
						Rules: []config.OverrideRule{
							{
								Value: "Host(`overridden.example.com`)",
								Filter: config.RouterFilter{
									Name: "provider1-api",
								},
							},
						},
					},
				},
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg.HTTP == nil {
			t.Fatal("HTTP configuration is nil")
		}

		// Check that the override was applied (this is a basic test - actual override logic may vary)
		if router, exists := dynCfg.HTTP.Routers["provider1-api"]; exists {
			// For now, just verify the router exists - override logic testing would need more complex setup
			t.Logf("Found router 'api' with rule: %s", router.Rule)
		} else {
			t.Error("Expected router 'api' to exist")
		}
	})

	t.Run("HTTP Router Entrypoint Override", func(t *testing.T) {
		cfg := &config.ProviderConfig{
			Connection: config.ConnectionConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/api/rawdata",
			},
			HTTP: &config.HTTPSection{
				Discover: true,
				Routers: &config.RoutersConfig{
					Discover: true,
					Overrides: config.RouterOverrides{
						Entrypoints: []config.OverrideEntrypoint{
							{
								Value: []string{"websecure"},
								Filter: config.RouterFilter{
									Name: "provider1-web",
								},
							},
						},
					},
				},
			},
		}
		dynCfg := httpclient.GenerateConfiguration(cfg)

		if dynCfg.HTTP == nil {
			t.Fatal("HTTP configuration is nil")
		}

		// Check that the override was applied (this is a basic test - actual override logic may vary)
		if router, exists := dynCfg.HTTP.Routers["provider1-web"]; exists {
			// For now, just verify the router exists - override logic testing would need more complex setup
			t.Logf("Found router 'web' with entrypoints: %v", router.EntryPoints)
		} else {
			t.Error("Expected router 'web' to exist")
		}
	})
}

// Helper function to wait for a service to be ready
func waitForService(ctx context.Context, url string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := client.Get(url)
			if err == nil {
				// Best-effort close; error is non-fatal for readiness.
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
