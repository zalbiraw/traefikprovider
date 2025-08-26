package config

// RoutersConfig holds discovery, matcher, and override settings for routers.
type RoutersConfig struct {
	Discover             bool            `json:"discover,omitempty" yaml:"discover,omitempty"`
	DiscoverPriority     bool            `json:"discoverPriority,omitempty" yaml:"discoverPriority,omitempty"`
	Matcher              string          `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	StripServiceProvider bool            `json:"stripServiceProvider,omitempty" yaml:"stripServiceProvider,omitempty"`
	Overrides            RouterOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes          []interface{}   `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

// RouterOverrides defines override rules applied to matched routers.
type RouterOverrides struct {
	Name        string               `json:"name,omitempty" yaml:"name,omitempty"`
	Rules       []OverrideRule       `json:"rules,omitempty" yaml:"rules,omitempty"`
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
	Middlewares []OverrideMiddleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

// OverrideRule applies a rule value to matching routers.
type OverrideRule struct {
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
	Matcher string `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}

// OverrideEntrypoint applies entrypoint values to matching routers.
type OverrideEntrypoint struct {
	Value   interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Matcher string      `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}

// OverrideService applies a service value to matching routers.
type OverrideService struct {
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
	Matcher string `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}

// OverrideMiddleware applies middleware values to matching routers.
type OverrideMiddleware struct {
	Value   interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Matcher string      `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}
