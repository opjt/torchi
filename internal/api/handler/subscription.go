package handler

import (
	"context"
	"torchi/internal/api/wrapper"
	"torchi/internal/domain/push"
	"torchi/internal/pkg/log"
	"torchi/internal/pkg/token"

	"github.com/go-chi/chi/v5"
)

type SubscriptionHandler struct {
	log     *log.Logger
	service *push.PushService
}

func NewSubscriptionHandler(log *log.Logger, service *push.PushService) *SubscriptionHandler {
	return &SubscriptionHandler{
		log:     log,
		service: service,
	}
}
func (h *SubscriptionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", wrapper.WrapJson(h.Subscribe, h.log.Error, wrapper.RespondJSON))
	r.Post("/unsubscribe", wrapper.WrapJson(h.Unsubscribe, h.log.Error, wrapper.RespondJSON))
	return r
}

type reqSubscribe struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// Subscribe, push notiification
func (h *SubscriptionHandler) Subscribe(ctx context.Context, req reqSubscribe) (interface{}, error) {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		h.log.Info("...", "error", err)
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

// Unsubscribe
func (h *SubscriptionHandler) Unsubscribe(ctx context.Context, req reqSubscribe) (interface{}, error) {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	h.log.Info("...", "user", userClaim.UserID)
	h.log.Info("req", "sub", req)
	if err := h.service.Unsubscribe(ctx, push.Subscription{
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
