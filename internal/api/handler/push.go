package handler

import (
	"context"
	"ohp/internal/api/wrapper"
	"ohp/internal/domain/push"
	"ohp/internal/pkg/log"
	"ohp/internal/pkg/token"

	"github.com/go-chi/chi/v5"
)

type PushHandler struct {
	log     *log.Logger
	service *push.PushService
}

func NewPushHandler(log *log.Logger, service *push.PushService) *PushHandler {
	return &PushHandler{
		log:     log,
		service: service,
	}
}
func (h *PushHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/subscribe", wrapper.WrapJson(h.Subscribe, h.log.Error, wrapper.RespondJSON))
	r.Post("/push", wrapper.WrapJson(h.Push, h.log.Error, wrapper.RespondJSON))
	return r
}

type reqSubscribe struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// Subscribe, push subscribe
func (h *PushHandler) Subscribe(ctx context.Context, req reqSubscribe) (interface{}, error) {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	h.log.Info("...", "user", userClaim.UserID)
	h.log.Info("req", "sub", req)
	if err := h.service.Subscribe(ctx, push.Subscription{
		UserID:   userClaim.UserID,
		Endpoint: req.Endpoint,
		P256dh:   req.Keys.P256dh,
		Auth:     req.Keys.Auth,
	}); err != nil {
		return nil, err
	}

	h.log.Info("push subscribe", "endpoint", req.Endpoint)

	return "success1", nil
}

// Push, notification
func (h *PushHandler) Push(ctx context.Context, req reqSubscribe) (interface{}, error) {

	if err := h.service.Push(ctx, push.Subscription{
		Endpoint: req.Endpoint,
		P256dh:   req.Keys.P256dh,
		Auth:     req.Keys.Auth,
	}); err != nil {
		return nil, err
	}
	h.log.Info("noti push", "endpoint", req.Endpoint)
	return "success", nil
}
func (h *PushHandler) Broadcast(ctx context.Context, req any) error {
	return nil
}
