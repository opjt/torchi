package handler

import (
	"io"
	"net/http"
	"ohp/internal/api/wrapper"
	"ohp/internal/domain/push"
	"ohp/internal/pkg/config"
	"ohp/internal/pkg/log"

	"github.com/go-chi/chi/v5"
)

type ApiHandler struct {
	log     *log.Logger
	service *push.PushService
}

func NewApiHandler(
	log *log.Logger,
	env config.Env,

	service *push.PushService,
) *ApiHandler {
	return &ApiHandler{
		log:     log,
		service: service,
	}
}
func (h *ApiHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/push/{token}", h.Push)

	return r
}

type resPush struct {
	Sent uint64 `json:"sent"`
}

func (h *ApiHandler) Push(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParamFromCtx(ctx, "token")
	h.log.Info("...", "token", token)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	message := string(bodyBytes)
	if message == "" {
		http.Error(w, "Message body is empty", http.StatusBadRequest)
		return
	}

	count, err := h.service.Push(ctx, token, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wrapper.RespondJSON(w, http.StatusOK, resPush{
		Sent: count,
	})
}
