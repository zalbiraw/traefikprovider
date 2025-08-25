package filters

import (
	"encoding/json"
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

// filterMapByNameRegex filter a map of typed objects by name as regex if set, otherwise exact match.
func filterMapByNameRegex[T any, R any](
	input map[string]*T,
	name string,
) map[string]R {
	result := make(map[string]R)
	for k, v := range input {
		if name != "" {
			matched, err := regexMatch(name, k)
			if err != nil || !matched {
				continue
			}
		}
		result[k] = any(v).(R)
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

func unmarshalRouter(routerMap map[string]interface{}, router *dynamic.Router) error {
	b, err := json.Marshal(routerMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, router)
}
