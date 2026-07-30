//go:build tray

// System-tray entrypoint (Windows): shows a taskbar icon with Pause/Resume/Quit.
// Build:  GOOS=windows GOARCH=amd64 go build -tags tray -o opencode-gateway.exe
// The HTTP server runs in a goroutine; Pause calls Shutdown, Resume restarts it.
package main

import (
	"context"
	_ "embed"
	"log"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/stretchr/testify/assert/yaml"

	"opencode-gateway/internal/gateway"
)

//go:embed icon.ico
var iconData []byte

var (
	mu  sync.Mutex
	srv *http.Server
	gw  *gateway.Server
	cfg gateway.Config
)

func startServer() {
	mu.Lock()
	defer mu.Unlock()
	if srv != nil {
		return
	}
	srv = &http.Server{Addr: cfg.Addr, Handler: gw.Handler()}
	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}(srv)
}

func stopServer() {
	mu.Lock()
	s := srv
	srv = nil
	mu.Unlock()
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.Shutdown(ctx)
	}
}

func main() {
	cfg = gateway.DefaultConfig()

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

	gw = gateway.New(cfg, processedKeys, providers)
	systray.Run(func() { onReady(len(processedKeys) > 0) }, onExit)
}

// onExit stops the HTTP server and closes the request log file.
// Must not run on Pause — only on real Quit / process exit.
func onExit() {
	stopServer()
	_ = gw.Close()
}

func onReady(hasKey bool) {
	systray.SetIcon(iconData)
	systray.SetTitle("opencode-gateway")
	systray.SetTooltip("opencode-gateway — running on " + cfg.Addr)

	mStatus := systray.AddMenuItem("Running on "+cfg.Addr, "")
	mStatus.Disable()
	if !hasKey {
		mStatus.SetTitle("⚠ No API key (add opencode-key.txt next to the exe)")
	}
	systray.AddSeparator()
	mToggle := systray.AddMenuItem("Pause", "Pause or resume the gateway")
	mQuit := systray.AddMenuItem("Quit", "Stop the gateway and exit")

	startServer()
	running := true

	for {
		select {
		case <-mToggle.ClickedCh:
			if running {
				stopServer()
				mToggle.SetTitle("Resume")
				mStatus.SetTitle("Paused")
				systray.SetTooltip("opencode-gateway — paused")
			} else {
				startServer()
				mToggle.SetTitle("Pause")
				mStatus.SetTitle("Running on " + cfg.Addr)
				systray.SetTooltip("opencode-gateway — running on " + cfg.Addr)
			}
			running = !running
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}
