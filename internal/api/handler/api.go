package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"torchi/internal/api/wrapper"
	"torchi/internal/domain/push"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ApiHandler struct {
	log     *log.Logger
	service *push.PushService

	frontUrl string
}

func NewApiHandler(
	log *log.Logger,
	env config.Env,

	service *push.PushService,
) *ApiHandler {
	return &ApiHandler{
		log:     log,
		service: service,

		frontUrl: env.FrontUrl,
	}
}
func (h *ApiHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/push/{token}", h.Push)
	r.Post("/demo", h.Demo)
	r.Post("/push/{token}/ask", h.Ask)
	r.Post("/react/{notiID}", wrapper.WrapJson(h.React, h.log.Error, wrapper.RespondJSON))
	r.Post("/push-test", wrapper.WrapJson(h.TestPush, h.log.Error, wrapper.RespondJSON))
	r.Post("/push-demo", wrapper.WrapJson(h.DemoPush, h.log.Error, wrapper.RespondJSON))

	return r
}

type reqDemoPush struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	P256dh   string `json:"p256dh"`
	Message  string `json:"message"`
}

func (h *ApiHandler) DemoPush(ctx context.Context, req reqDemoPush) (interface{}, error) {

	h.service.DemoPush(ctx, push.DemoPushParams{
		Endpoint: req.Endpoint,
		Auth:     req.Auth,
		P256dh:   req.P256dh,
	}, req.Message)
	return nil, nil
}

type reqTestPush struct {
	Endpoint string `json:"endpoint"`
}

func (h *ApiHandler) TestPush(ctx context.Context, req reqTestPush) (interface{}, error) {

	h.log.Info("req", "sub", req)

	if err := h.service.PushByEndpoint(ctx, req.Endpoint, fmt.Sprintf(
		"Hello World %s",
		time.Now().Format("2006-01-02 15:04:05"),
	)); err != nil {
		return nil, err
	}
	return nil, nil

}

func (h *ApiHandler) Demo(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{

		"message": "이 엔드포인트는 데모용입니다. CLI 호출은 지원하지 않으며, 웹에서 로그인 후 푸시 구독을 완료해야 사용할 수 있습니다.",
		"action":  h.frontUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden) // 403

	_ = json.NewEncoder(w).Encode(response)

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
func (h *ApiHandler) Ask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParamFromCtx(ctx, "token")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	msg := r.FormValue("msg")

	var actions []string
	if raw := r.FormValue("actions"); raw != "" {
		for _, a := range strings.Split(raw, ",") {
			if t := strings.TrimSpace(a); t != "" {
				actions = append(actions, t)
			}
		}
	}

	timeout := 300
	if t := r.FormValue("timeout"); t != "" {
		if v, err := strconv.Atoi(t); err == nil {
			timeout = v
		}
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	h.log.Debug("...", "token", token, "msg", msg, "actions", actions, "timeout", timeout)
	reaction, err := h.service.PushAndWait(timeoutCtx, token, msg, actions)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			wrapper.RespondJSON(w, http.StatusRequestTimeout, map[string]string{
				"status": "timeout",
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(reaction)); err != nil {
		h.log.Error("failed to write reaction response", "err", err)
	}
}

type reqReact struct {
	Reaction string `json:"reaction"`
}

func (h *ApiHandler) React(ctx context.Context, req reqReact) (interface{}, error) {

	notiID := chi.URLParamFromCtx(ctx, "notiID")
	notiUUID, err := uuid.Parse(notiID)
	if err != nil {
		return nil, err
	}

	if err := h.service.React(ctx, notiUUID, req.Reaction); err != nil {
		return nil, err
	}

	return nil, nil
}
