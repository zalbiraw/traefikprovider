package filters

import (
	"regexp"
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

// filterMapByNameRegex filters a map of typed objects by name (regex if set) and provider postfix (regex if set).
// The provider is extracted from the resource name as the substring after the last '@'.
func filterMapByNameRegex[T, R any](input map[string]*T, name, provider string) map[string]R {
	result := make(map[string]R)
	for k, v := range input {
		if name != "" {
			matched, err := regexMatch(name, k)
			if err != nil || !matched {
				continue
			}
		}
		if provider != "" {
			prov := extractProviderFromName(k)
			matched, err := regexMatch(provider, prov)
			if err != nil || !matched {
				continue
			}
		}
		if rv, ok := any(v).(R); ok {
			result[k] = rv
		}
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

// extractProviderFromName retrieves the postfix after the last '@' in a resource name.
// If no '@' is present, it returns an empty string.
func extractProviderFromName(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '@' {
			if i+1 < len(name) {
				return name[i+1:]
			}
			return ""
		}
	}
	return ""
}
