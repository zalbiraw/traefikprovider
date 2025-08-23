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

// filterMapByNameAndRegex filters a map[string]interface{} by name and nameRegex fields.
func filterMapByNameAndRegex(
	input map[string]interface{},
	name, nameRegex string,
) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	for k, v := range input {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if name != "" && k != name {
			continue
		}
		if nameRegex != "" {
			matched, err := regexMatch(nameRegex, k)
			if err != nil || !matched {
				continue
			}
		}
		result[k] = item
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
