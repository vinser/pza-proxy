package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vinser/pza-proxy/internal/config"
	"github.com/vinser/pza-proxy/internal/polza"
)

// ChatHandler handles /chat/completions and /v1/chat/completions.
type ChatHandler struct{}

func NewChatHandler() http.Handler {
	return &ChatHandler{}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Debug logging of incoming request from Cline
	log.Println("---- Incoming request ----")
	log.Printf("Path: %s", r.URL.Path)
	log.Println("Headers:")
	for k, v := range r.Header {
		log.Printf("  %s: %v", k, v)
	}

	// Read and log body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	log.Printf("Body: %s", string(body))
	defer r.Body.Close()

	// Extract API key from incoming request
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Parse incoming OpenAI-compatible request
	var chatReq polza.ChatRequest
	if err := json.Unmarshal(body, &chatReq); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Resolve alias → real model config
	var modelCfg config.ModelConfig
	if cfg, ok := config.ResolveModel(chatReq.Model); ok {
		modelCfg = cfg
		chatReq.Model = cfg.ID
	} else {
		modelCfg = config.ModelConfig{}
	}

	// Build Polza request
	pReq := polza.PolzaRequest{
		Model:       chatReq.Model,
		Messages:    chatReq.Messages,
		Stream:      chatReq.Stream,
		Temperature: chatReq.Temperature,
		Provider: polza.ProviderLimits{
			Order:          modelCfg.Provider.Order,
			Allow:          modelCfg.Provider.Allow,
			Deny:           modelCfg.Provider.Deny,
			AllowFallbacks: modelCfg.Provider.AllowFallbacks,
			MaxPrice: func() *polza.MaxPrice {
				if modelCfg.Provider.MaxPrice == nil {
					return nil
				}
				return &polza.MaxPrice{
					Prompt:     modelCfg.Provider.MaxPrice.Prompt,
					Completion: modelCfg.Provider.MaxPrice.Completion,
				}
			}(),
		},
	}

	pBody, err := json.Marshal(pReq)
	if err != nil {
		http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
		return
	}

	// Create Polza client using Authorization header from Cline
	client := polza.NewClient(auth)

	// Send request to Polza
	upstreamResp, err := client.Do(pBody, chatReq.Stream)
	if err != nil {
		log.Printf("upstream error: %v", err)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer upstreamResp.Body.Close()

	// Copy upstream headers
	for k, v := range upstreamResp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// Propagate status code
	w.WriteHeader(upstreamResp.StatusCode)

	// Transparent streaming (SSE or plain JSON)
	if flusher, ok := w.(http.Flusher); ok {
		buf := make([]byte, 32*1024)
		for {
			n, readErr := upstreamResp.Body.Read(buf)
			if n > 0 {
				if _, writeErr := w.Write(buf[:n]); writeErr != nil {
					return
				}
				flusher.Flush()
			}
			if readErr != nil {
				break
			}
		}
	} else {
		io.Copy(w, upstreamResp.Body)
	}
}
