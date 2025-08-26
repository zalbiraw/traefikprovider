package matchers

import (
	"regexp"
	"strings"

	"github.com/zalbiraw/traefikprovider/internal/rules"
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

// combineRules concatenates provider-level and section-level rules with logical AND.
// If either side is empty, returns the other; if both empty, returns empty.
func combineRules(providerRule, sectionRule string) string {
	p := strings.TrimSpace(providerRule)
	s := strings.TrimSpace(sectionRule)
	switch {
	case p == "" && s == "":
		return ""
	case p == "":
		return s
	case s == "":
		return p
	default:
		return "(" + p + ") && (" + s + ")"
	}
}

// compileRule compiles a rule string into a Program.
func compileRule(rule string) (*rules.Program, error) {
	return rules.Compile(rule)
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
