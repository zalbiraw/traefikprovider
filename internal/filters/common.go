package filters

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"

	"github.com/traefik/genconf/dynamic"
)

// regexMatch is a helper to match a string with a regex pattern.
func regexMatch(pattern, value string) (bool, error) {
	if pattern == "" {
		return true, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(value), nil
}

// filterMapByNameRegex filters a map[string]interface{} by name as regex if set, otherwise exact match.
func filterMapByNameRegex(
	input map[string]interface{},
	name string,
) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		// Skip internal services by default
		if k != "" && len(k) > 9 && k[len(k)-9:] == "@internal" {
			continue
		}
		
		if name != "" {
			matched, err := regexMatch(name, k)
			if err != nil || !matched {
				continue
			}
		}
		result[k] = v
	}
	return result
}

// routerEntrypointsMatch checks if router entrypoints contain all filter entrypoints.
func routerEntrypointsMatch(routerEPs, filterEPs []string) bool {
	if len(filterEPs) == 0 {
		return true
	}
	set := make(map[string]struct{})
	for _, ep := range routerEPs {
		set[ep] = struct{}{}
	}
	for _, ep := range filterEPs {
		if _, ok := set[ep]; !ok {
			return false
		}
	}
	return true
}

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

func unmarshalRouter(routerMap map[string]interface{}, router *dynamic.Router) error {
	b, err := json.Marshal(routerMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, router)
}
