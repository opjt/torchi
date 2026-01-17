package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ohp/internal/api/wrapper"
	"ohp/internal/domain/push"
	"ohp/internal/pkg/config"
	"ohp/internal/pkg/log"
	"time"

	"github.com/go-chi/chi/v5"
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
