package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sdelacru/yodai/internal/config"
)

type Handler struct {
	client *Client
	cfg    *config.Config
}

func NewHandler(client *Client, cfg *config.Config) *Handler {
	return &Handler{client: client, cfg: cfg}
}

type ChatRequest struct {
	Message string    `json:"message"`
	History []Message `json:"history"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func (h *Handler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	resp, err := h.client.Chat(r.Context(), req.Message, req.History, h.cfg.MaxTokens, h.cfg.Temperature, h.cfg.TopP)
	if err != nil {
		log.Printf("chat error: %v", err)
		http.Error(w, "inference error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: resp})
}

func (h *Handler) HandleStreamChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	tokens, errs := h.client.StreamChat(r.Context(), req.Message, req.History, h.cfg.MaxTokens, h.cfg.Temperature, h.cfg.TopP)

	for {
		select {
		case token, ok := <-tokens:
			if !ok {
				fmt.Fprintf(w, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}
			data, _ := json.Marshal(map[string]string{"token": token})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case err, ok := <-errs:
			if ok && err != nil {
				log.Printf("stream error: %v", err)
				data, _ := json.Marshal(map[string]string{"error": err.Error()})
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
			return
		case <-r.Context().Done():
			return
		}
	}
}
