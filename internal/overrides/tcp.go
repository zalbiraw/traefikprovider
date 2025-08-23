package overrides

import (
	"github.com/traefik/genconf/dynamic"
)

// OverrideTCPRouters will later apply overrides to filtered TCP routers.
import "github.com/zalbiraw/traefik-provider/config"

func OverrideTCPRouters(filtered map[string]*dynamic.TCPRouter, overrides config.RouterOverrides) {
	// TODO: Implement override logic
}

// OverrideTCPServices will later apply overrides to filtered TCP services.
func OverrideTCPServices(filtered map[string]*dynamic.TCPService, overrides config.ServiceOverrides) {
	// TODO: Implement override logic
}

// OverrideTCPMiddlewares will later apply overrides to filtered TCP middlewares.
func OverrideTCPMiddlewares(filtered map[string]*dynamic.TCPMiddleware, overrides []interface{}) {
	// TODO: Implement override logic
}
