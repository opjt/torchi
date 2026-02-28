package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"torchi/internal/domain/sse"
	"torchi/internal/pkg/log"
	"torchi/internal/pkg/token"

	"github.com/go-chi/chi/v5"
)

type SSEHandler struct {
	log    *log.Logger
	broker *sse.Broker
}

func NewSSEHandler(log *log.Logger, broker *sse.Broker) *SSEHandler {
	return &SSEHandler{log: log, broker: broker}
}

func (h *SSEHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/notifications", h.Stream)
	return r
}
func (h *SSEHandler) Stream(w http.ResponseWriter, r *http.Request) {
	userClaim, err := token.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
	ch := h.broker.Subscribe(userClaim.UserID)
	defer func() {
		h.broker.Unsubscribe(userClaim.UserID, ch)
	}()

	// send 헬퍼 - 실패 시 false 반환
	send := func(format string, args ...any) bool {
		_, err := fmt.Fprintf(w, format, args...)
		if err != nil {
			h.log.Info("client disconnected while writing", "userID", userClaim.UserID)
			return false
		}
		flusher.Flush()
		return true
	}

	if !send("event: connected\ndata: {}\n\n") {
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			h.log.Info("context done")
			return
		case <-h.broker.Done():
			return
		case <-ticker.C:
			if !send(": heartbeat\n\n") {
				return
			}
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(event.Data)
			if err != nil {
				h.log.Error("sse marshal error", "err", err)
				continue
			}
			if !send("event: %s\ndata: %s\n\n", event.Event, string(data)) {
				return
			}
		}
	}
}
