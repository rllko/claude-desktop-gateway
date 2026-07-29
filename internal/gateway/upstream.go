package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) callUpstream(body oaiReq, provider ProviderConfig) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", provider.BaseURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "Bearer "+provider.APIKey)
	req.Header.Set("content-type", "application/json")
	fmt.Printf("%s", body)
	// Cloudflare (error 1010) blocks the default Go/Python client UA.
	if provider.UserAgent == "" {
		provider.UserAgent = "Mozilla/5.0 (opencode-gateway)"
	}

	req.Header.Set("user-agent", provider.UserAgent)
	return s.client.Do(req)
}
