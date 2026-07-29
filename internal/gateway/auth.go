package gateway

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func readKeyFile(p string) string {
	if len(p) == 0 {
		return ""
	}

	if b, err := os.ReadFile(p); err == nil {
		return strings.TrimSpace(string(b))
	}

	return ""
}

// takes in the env var map and returns the loaded keys, not present if not found
func loadEnvKeys(providerEnvVars map[string]string) map[string]string {
	out := map[string]string{}

	for keyName, envVar := range providerEnvVars {
		if token := strings.TrimSpace(os.Getenv(envVar)); token != "" {
			out[keyName] = token
		}
	}

	if k := readKeyFile(os.Getenv("OPENCODE_KEY_FILE")); out["opencode-go"] == "" && k != "" {
		out["opencode-go"] = k
	}

	return out
}

func LoadAPIKeys(providerEnvVars map[string]string) map[string]string {
	// load from the opencode location
	// lets make it one of the options to have keys there
	keys := loadAuthFromFile()

	// now look at the env vars in the yaml file
	envKeys := loadEnvKeys(providerEnvVars)

	for k, v := range envKeys {
		if _, ok := keys[k]; !ok {
			keys[k] = v
		}
	}

	return keys
}

// loadAuthFromFile reads the api key opencode itself stores in auth.json,
// shaped as { "<provider>": { "type":"api", "key":"..." }, ... }.
// Prefers the "opencode-go" provider, then "opencode".
func loadAuthFromFile() map[string]string {
	var paths []string

	p := os.Getenv("OPENCODE_AUTH_FILE")
	paths = append(paths, p)

	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		p := filepath.Join(x, "opencode", "auth.json")
		paths = append(paths, p)
	}

	if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, ".local", "share", "opencode", "auth.json") // real path on Win/Linux/mac today
		paths = append(paths, p)

		switch runtime.GOOS {
		case "windows":
			if ad := os.Getenv("APPDATA"); ad != "" {
				p := filepath.Join(ad, "opencode", "auth.json")
				paths = append(paths, p)
			}
		case "darwin":
			p := filepath.Join(home, "Library", "Application Support", "opencode", "auth.json")
			paths = append(paths, p)
		}
	}

	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}

		var m map[string]struct {
			Key string `json:"key"`
		}

		if json.Unmarshal(b, &m) != nil {
			continue
		}

		keys := map[string]string{}
		for name, k := range m {
			keys[name] = k.Key
		}

		return keys
	}

	return nil
}
