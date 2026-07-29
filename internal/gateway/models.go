package gateway

// Model maps a Desktop-facing alias to a real opencode model + picker label,
// with the model's real context and max-output token limits (from opencode's
// models.json) so Desktop's context meter and output caps are accurate.
// Aliases are typo'd on purpose to slip past Desktop's third-party-brand filter.
type YamlConfig struct {
	ExtraApiKeys map[string]string         `yaml:"extra-api-keys,omitempty"`
	Providers    map[string]ProviderConfig `yaml:"providers"`
}

// the Provider level (zenapi, agent-router, etc.)
type ProviderConfig struct {
	Enabled   bool                   `yaml:"enabled"`
	APIType   string                 `yaml:"api_type,omitempty"`
	APIKey    string                 `yaml:"env_var,omitempty"`
	BaseURL   string                 `yaml:"base_url"`
	UserAgent string                 `yaml:"user_agent"`
	Models    map[string]ModelConfig `yaml:"models"`
}

// / the Model level (big-pickle, kimi-k3, etc.)
type ModelConfig struct {
	Enabled bool   `yaml:"enabled"`
	Label   string `yaml:"label"`
	Alias   string `yaml:"alias"`
	Real    string `yaml:"real"`
	MaxIn   int    `yaml:"max_in"`
	MaxOut  int    `yaml:"max_out"`
	Vision  bool   `yaml:"vision"`
}
