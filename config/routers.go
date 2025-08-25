package config

type RoutersConfig struct {
	Discover         bool            `json:"discover,omitempty" yaml:"discover,omitempty"`
	DiscoverPriority bool            `json:"discoverPriority,omitempty" yaml:"discoverPriority,omitempty"`
	Filters          RouterFilters   `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides        RouterOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes      []interface{}   `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

type RouterFilters struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
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
	Value   string        `json:"value,omitempty" yaml:"value,omitempty"`
	Filters RouterFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideEntrypoint struct {
	Value   interface{}   `json:"value,omitempty" yaml:"value,omitempty"`
	Filters RouterFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideService struct {
	Value   string        `json:"value,omitempty" yaml:"value,omitempty"`
	Filters RouterFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type OverrideMiddleware struct {
	Value   interface{}   `json:"value,omitempty" yaml:"value,omitempty"`
	Filters RouterFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
}
