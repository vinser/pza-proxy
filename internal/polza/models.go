package polza

import "encoding/json"

// ChatRequest is the OpenAI-compatible request coming from Cline.
type ChatRequest struct {
	Model       string          `json:"model"`
	Messages    json.RawMessage `json:"messages"`
	Stream      bool            `json:"stream,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// MaxPrice is the Polza max_price structure.
type MaxPrice struct {
	Prompt     float64 `json:"prompt,omitempty"`
	Completion float64 `json:"completion,omitempty"`
}

// ProviderLimits describes Polza provider constraints.
type ProviderLimits struct {
	Order          []string  `json:"order,omitempty"`
	Allow          []string  `json:"allow,omitempty"`
	Deny           []string  `json:"deny,omitempty"`
	MaxPrice       *MaxPrice `json:"max_price,omitempty"`
	AllowFallbacks *bool     `json:"allow_fallbacks,omitempty"`
}

// PolzaRequest is the request we send to Polza API.
type PolzaRequest struct {
	Model       string          `json:"model"`
	Messages    json.RawMessage `json:"messages"`
	Stream      bool            `json:"stream,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Provider    ProviderLimits  `json:"provider"`
}
