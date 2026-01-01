package push

import (
	"context"
	"encoding/json"
	"fmt"
	"ohp/internal/domain/endpoint"
	"ohp/internal/domain/token"
	"ohp/internal/pkg/config"
	"ohp/internal/pkg/log"

	"github.com/SherClockHolmes/webpush-go"
)

type PushService struct {
	// repo     SubscriptionRepository
	vapidKey config.Vapid
	log      *log.Logger

	tokenService    *token.TokenService
	endpointService *endpoint.EndpointService
}

func NewPushService(
	// repo SubscriptionRepository,
	env config.Env,
	log *log.Logger,

	tokenService *token.TokenService,
	endpointService *endpoint.EndpointService,
) *PushService {
	return &PushService{
		// repo:            repo,
		log:             log,
		vapidKey:        env.Vapid,
		tokenService:    tokenService,
		endpointService: endpointService,
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

// Push notification using endpoint token
func (s *PushService) Push(ctx context.Context, endpointToken string, message string) (uint64, error) {
	var count uint64

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
		return count, err
	}

	for _, token := range tokens {
		if err := s.pushNotification(token, endpoint.Name, message); err != nil {
			return count, err
		}
		count = count + 1
	}

	return count, nil
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
