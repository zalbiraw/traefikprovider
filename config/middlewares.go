package config

type MiddlewaresConfig struct {
	Discover         bool              `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters          MiddlewareFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
	ExtraMiddlewares []interface{}     `json:"extraMiddlewares,omitempty" yaml:"extraMiddlewares,omitempty"`
}

type MiddlewareFilters struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
