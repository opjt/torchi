package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"torchi/internal/domain/endpoint"
	"torchi/internal/domain/notifications"
	"torchi/internal/domain/token"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"github.com/SherClockHolmes/webpush-go"
)

type PushService struct {
	vapidKey config.Vapid
	log      *log.Logger

	tokenService    *token.TokenService
	endpointService *endpoint.EndpointService
	notiService     *notifications.NotiService
}

func NewPushService(
	env config.Env,
	log *log.Logger,

	tokenService *token.TokenService,
	endpointService *endpoint.EndpointService,
	notiService *notifications.NotiService,
) *PushService {
	return &PushService{
		log:             log,
		vapidKey:        env.Vapid,
		tokenService:    tokenService,
		endpointService: endpointService,
		notiService:     notiService,
	}
}

func (s *PushService) Subscribe(ctx context.Context, sub Subscription) error {

	if err := s.tokenService.Register(ctx, token.Token{
		P256dh:   sub.P256dh,
		Auth:     sub.Auth,
		UserID:   sub.UserID,
		EndPoint: sub.Endpoint,
	}); err != nil {
		return err
	}

	return nil
}

func (s *PushService) Unsubscribe(ctx context.Context, sub Subscription) error {

	if err := s.tokenService.Unregister(ctx, token.Token{
		P256dh:   sub.P256dh,
		Auth:     sub.Auth,
		EndPoint: sub.Endpoint,
	}); err != nil {
		return err
	}

	return nil
}

func (s *PushService) PushByEndpoint(ctx context.Context, endpoint string, message string) error {
	token, err := s.tokenService.FindByEndpoint(ctx, endpoint)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("endpoint not found")
	}

	if err := s.pushNotification(*token, "TEST!", message); err != nil {
		return err
	}

	return nil

}

// Push notification using endpoint token
func (s *PushService) Push(ctx context.Context, endpointToken string, message string) (uint64, error) {

	endpoint, err := s.endpointService.FindByToken(ctx, endpointToken)

	if err != nil {
		return 0, err
	}
	if endpoint == nil {
		return 0, nil
	}
	userID := endpoint.UserID

	tokens, err := s.tokenService.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	noti, err := s.notiService.Register(ctx, notifications.ReqRegister{
		EndpointID:         endpoint.ID,
		UserID:             userID,
		Body:               message,
		NotificationEnable: endpoint.NotificationEnable,
	})

	if !endpoint.NotificationEnable {
		return 0, err
	}

	var count uint64

	for _, token := range tokens {
		if err := s.pushNotification(token, endpoint.Name, message); err != nil {
			// TODO: 에러 처리 개선 필요.
			return count, err
		}
		count = count + 1
	}
	if err = s.notiService.UpdateStatusSent(ctx, noti.ID); err != nil {
		return count, err
	}

	return count, nil
}

type DemoPushParams struct {
	Endpoint string
	Auth     string
	P256dh   string
}

func (s *PushService) DemoPush(ctx context.Context, req DemoPushParams, message string) (interface{}, error) {

	if err := s.pushNotification(token.Token{
		P256dh:   req.P256dh,
		Auth:     req.Auth,
		EndPoint: req.Endpoint,
	}, "Demo", message); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *PushService) pushNotification(token token.Token, title, body string) error {

	subs := &webpush.Subscription{
		Endpoint: token.EndPoint,
		Keys: webpush.Keys{
			P256dh: token.P256dh,
			Auth:   token.Auth,
		},
	}

	options := &webpush.Options{
		VAPIDPublicKey:  s.vapidKey.PublicKey,
		VAPIDPrivateKey: s.vapidKey.PrivateKey,
		TTL:             300,
		Subscriber:      "jtpark1957@gmail.com",
	}
	payload := map[string]interface{}{
		"title": title,
		"body":  body,
		"data": map[string]string{
			"url":       "/",
			"timestamp": fmt.Sprintf("%d", 1234567890),
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	resp, err := webpush.SendNotification(payloadBytes, subs, options)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil

}
