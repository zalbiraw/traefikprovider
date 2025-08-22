package httpclient

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/traefik/genconf/dynamic"
)

// extractRouterPriority extracts and validates the priority from the router map.
func extractRouterPriority(routerMap map[string]interface{}, name string) int {
	prio, ok := routerMap["priority"]
	if !ok {
		return 0
	}
	prioFloat, ok := prio.(float64)
	if ok && prioFloat == float64(int(prioFloat)) && prioFloat <= float64(math.MaxInt64) && prioFloat >= float64(math.MinInt64) {
		return int(prioFloat)
	}
	fmt.Printf("[WARN] Invalid router priority for %s: %v\n", name, prio)
	return 0
}

// unmarshalRouter marshals and unmarshals the router map into the router struct.
func unmarshalRouter(routerMap map[string]interface{}, router *dynamic.Router) error {
	b, err := json.Marshal(routerMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, router)
}
