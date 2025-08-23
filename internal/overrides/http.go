package overrides

import (
	"github.com/traefik/genconf/dynamic"
)

// OverrideHTTPRouters will later apply overrides to filtered HTTP routers.
// For now, it only takes the filtered routers and a list of overrides.
import "github.com/zalbiraw/traefik-provider/config"

func OverrideHTTPRouters(filtered map[string]*dynamic.Router, overrides config.RouterOverrides) {
	// TODO: Implement override logic
}

// OverrideHTTPServices will later apply overrides to filtered HTTP services.
func OverrideHTTPServices(filtered map[string]*dynamic.Service, overrides config.ServiceOverrides) {
	// TODO: Implement override logic
}

// OverrideHTTPMiddlewares will later apply overrides to filtered HTTP middlewares.
func OverrideHTTPMiddlewares(filtered map[string]*dynamic.Middleware, overrides []interface{}) {
	// TODO: Implement override logic
}

// OverrideHTTPServerTransports will later apply overrides to filtered HTTP server transports.
func OverrideHTTPServerTransports(filtered map[string]*dynamic.ServersTransport, overrides []interface{}) {
	// TODO: Implement override logic
}
