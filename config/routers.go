package config

type RoutersConfig struct {
	Discover         bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	DiscoverPriority bool             `json:"discoverPriority,omitempty" yaml:"discoverPriority,omitempty"`
	Filters          RouterFilters    `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides        RouterOverrides  `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes      []interface{}    `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

type RouterFilters struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	NameRegex   string   `json:"nameRegex,omitempty" yaml:"nameRegex,omitempty"`
	Entrypoints []string `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Rule        string   `json:"rule,omitempty" yaml:"rule,omitempty"`
	RuleRegex   string   `json:"ruleRegex,omitempty" yaml:"ruleRegex,omitempty"`
	Service     string   `json:"service,omitempty" yaml:"service,omitempty"`
	ServiceRegex string  `json:"serviceRegex,omitempty" yaml:"serviceRegex,omitempty"`
}

type RouterOverrides struct {
	Rules       []OverrideRule       `json:"rules,omitempty" yaml:"rules,omitempty"`
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
	Middlewares []OverrideMiddleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

type OverrideRule struct {
	Mode    string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Values  string                 `json:"values,omitempty" yaml:"values,omitempty"`
	Operator string                `json:"operator,omitempty" yaml:"operator,omitempty"`
	Filters map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideEntrypoint struct {
	Mode    string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Values  []string               `json:"values,omitempty" yaml:"values,omitempty"`
	Filters map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideService struct {
	Mode    string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Values  string                 `json:"values,omitempty" yaml:"values,omitempty"`
	Filters map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideMiddleware struct {
	Mode    string                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Values  []string               `json:"values,omitempty" yaml:"values,omitempty"`
	Filters map[string]interface{}  `json:"filters,omitempty" yaml:"filters,omitempty"`
}
