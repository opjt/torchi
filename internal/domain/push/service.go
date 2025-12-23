package push

import (
	"context"
	"encoding/json"
	"fmt"
	"ohp/internal/domain/token"
	"ohp/internal/pkg/config"

	"github.com/SherClockHolmes/webpush-go"
)

type PushService struct {
	repo     SubscriptionRepository
	vapidKey config.Vapid

	tokenService *token.TokenService
}

func NewPushService(repo SubscriptionRepository, env config.Env, tokenService *token.TokenService) *PushService {
	return &PushService{
		repo:         repo,
		vapidKey:     env.Vapid,
		tokenService: tokenService,
	}
}

func (s *PushService) Subscribe(ctx context.Context, sub Subscription) error {

	if err := s.tokenService.Register(ctx, token.Token{
		P256dh: sub.P256dh,
		Auth:   sub.Auth,
		UserID: sub.UserID,
	}); err != nil {
		return err
	}

	subs := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}
	s.repo.Save(sub.Endpoint, subs)
	return nil
}
func (s *PushService) Push(ctx context.Context, sub Subscription) error {

	subs := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	options := &webpush.Options{
		VAPIDPublicKey:  s.vapidKey.PublicKey,
		VAPIDPrivateKey: s.vapidKey.PrivateKey,
		TTL:             300,
		Subscriber:      "jtpark1957@gmail.com",
	}
	payload := map[string]interface{}{
		"title": "OhP test notification",
		"body":  "PWA 푸시 알림 테스트 ",
		// "icon":  "/icon-192x192.png",
		// "badge": "/badge-72x72.png",
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
