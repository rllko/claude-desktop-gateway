//go:build !tray

// Plain headless entrypoint (default build): run the gateway and block.
// Used for the Linux/WSL build and testing. The Windows tray app is built
// with `-tags tray` (see tray.go).
package main

import (
	"log"
	"log/slog"
	"maps"
	"net/http"
	"os"

	"opencode-gateway/internal/gateway"

	"github.com/stretchr/testify/assert/yaml"
)

func main() {
	cfg := gateway.DefaultConfig()

	yamlFile, err := os.ReadFile("models.yaml")
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}

	var pConfig gateway.YamlConfig
	yaml.Unmarshal(yamlFile, &pConfig)

	p := make(map[string]string)
	providers := make([]gateway.ProviderConfig, 0)
	for _, config := range pConfig.Providers {
		if config.Enabled {
			if _, ok := p[config.APIType]; !ok {
				p[config.APIKey] = config.APIType
			}

			providers = append(providers, config)
		}
	}

	// get the keys, this way we can filter
	// name -> key ; example: openai -> KEY
	processedKeys := gateway.LoadAPIKeys(p)

	if pConfig.ExtraAPIKeys != nil {
		maps.Copy(processedKeys, pConfig.ExtraAPIKeys)
	}
	if len(processedKeys) == 0 {
		slog.Warn("no API key found — requests will 401 until one is set")
	}

	srv := gateway.New(cfg, processedKeys, providers)
	defer srv.Close()

	slog.Info("opencode-gateway starting",
		"addr", cfg.Addr, "models", srv.ModelCount())

	if err := http.ListenAndServe(cfg.Addr, srv.Handler()); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
