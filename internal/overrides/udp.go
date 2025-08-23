package overrides

import (
	"github.com/traefik/genconf/dynamic"
)

// OverrideUDPRouters will later apply overrides to filtered UDP routers.
import "github.com/zalbiraw/traefik-provider/config"

func OverrideUDPRouters(filtered map[string]*dynamic.UDPRouter, overrides config.UDPOverrides) {
	// TODO: Implement override logic
}

// OverrideUDPServices will later apply overrides to filtered UDP services.
func OverrideUDPServices(filtered map[string]*dynamic.UDPService, overrides config.ServiceOverrides) {
	// TODO: Implement override logic
}
