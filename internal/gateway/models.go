package gateway

// PS: THE ALIAS FOR THE MODEL NEEDS TO BE A TYPO OR CLAUDE DESKTOP IIS NOT GONNA SHOW IT
type YamlConfig struct {
	ExtraAPIKeys map[string]string         `yaml:"extra-api-keys,omitempty"`
	Providers    map[string]ProviderConfig `yaml:"providers"`
}

// ProviderConfig (zenapi, agent-router, etc.)
type ProviderConfig struct {
	ClaudeSystemPrompt bool                   `yaml:"claude_system_prompt"`
	Enabled            bool                   `yaml:"enabled"`
	APIType            string                 `yaml:"api_type,omitempty"`
	APIKey             string                 `yaml:"env_var,omitempty"`
	BaseURL            string                 `yaml:"base_url"`
	UserAgent          string                 `yaml:"user_agent"`
	Models             map[string]ModelConfig `yaml:"models"`
}

// ModelConfig (big-pickle, kimi-k3, etc.)
type ModelConfig struct {
	Enabled bool   `yaml:"enabled"`
	Label   string `yaml:"label"`
	Alias   string `yaml:"alias"`
	Real    string `yaml:"real"`
	MaxIn   int    `yaml:"max_in"`
	MaxOut  int    `yaml:"max_out"`
	Vision  bool   `yaml:"vision"`
}
