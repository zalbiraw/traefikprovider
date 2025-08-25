package config

type MiddlewaresConfig struct {
	Discover         bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter           MiddlewareFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	ExtraMiddlewares []interface{}    `json:"extraMiddlewares,omitempty" yaml:"extraMiddlewares,omitempty"`
}

type MiddlewareFilter struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
