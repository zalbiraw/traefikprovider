package config

type RoutersConfig struct {
	Discover             bool            `json:"discover,omitempty" yaml:"discover,omitempty"`
	DiscoverPriority     bool            `json:"discoverPriority,omitempty" yaml:"discoverPriority,omitempty"`
	Filter               RouterFilter    `json:"filter,omitempty" yaml:"filter,omitempty"`
	StripServiceProvider bool            `json:"stripServiceProvider,omitempty" yaml:"stripServiceProvider,omitempty"`
	Overrides            RouterOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes          []interface{}   `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

type RouterFilter struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Provider    string   `json:"provider,omitempty" yaml:"provider,omitempty"`
	Rule        string   `json:"rule,omitempty" yaml:"rule,omitempty"`
	Entrypoints []string `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Service     string   `json:"service,omitempty" yaml:"service,omitempty"`
}

type RouterOverrides struct {
	Name        string               `json:"name,omitempty" yaml:"name,omitempty"`
	Rules       []OverrideRule       `json:"rules,omitempty" yaml:"rules,omitempty"`
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
	Middlewares []OverrideMiddleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

type OverrideRule struct {
	Value  string       `json:"value,omitempty" yaml:"value,omitempty"`
	Filter RouterFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}

type OverrideEntrypoint struct {
	Value  interface{}  `json:"value,omitempty" yaml:"value,omitempty"`
	Filter RouterFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}

type OverrideService struct {
	Value  string       `json:"value,omitempty" yaml:"value,omitempty"`
	Filter RouterFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}

type OverrideMiddleware struct {
	Value  interface{}  `json:"value,omitempty" yaml:"value,omitempty"`
	Filter RouterFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
}
