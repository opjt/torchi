package endpoint

import (
	"context"
	"errors"
	"time"
	"torchi/internal/pkg/token"
)

const endpointLength = 11

type EndpointService struct {
	repo EndpointRepository
}

func NewEndpointService(
	repo EndpointRepository,
) *EndpointService {
	return &EndpointService{
		repo: repo,
	}
}

func (s *EndpointService) FindByToken(ctx context.Context, token string) (*Endpoint, error) {

	return s.repo.FindByToken(ctx, token)
}
func (s *EndpointService) List(ctx context.Context) ([]Endpoint, error) {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByUserID(ctx, userClaim.UserID)
}

func (s *EndpointService) UpdateMute(ctx context.Context, token string, notiEnable bool) error {
	if notiEnable {
		return s.repo.UpdateUnmute(ctx, token)
	} else {
		disabledTime := time.Now()
		return s.repo.UpdateMute(ctx, token, &disabledTime)
	}
}

func (s *EndpointService) Add(ctx context.Context, serviceName string) error {

	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return err
	}

	const maxRetry = 5

	for i := 0; i < maxRetry; i++ {
		endpoint, err := s.genEndpoint()
		if err != nil {
			return err
		}

		err = s.repo.Add(ctx, insertEndpointParams{
			userID:      userClaim.UserID,
			serviceName: serviceName,
			endpoint:    endpoint,
		})

		if err != nil {
			if err == ErrDuplicateToken {
				continue
			}
			return err
		}
		return nil

	}

	return errors.New("endpoint generate fail")
}

func (s *EndpointService) Remove(ctx context.Context, endpointToken string) error {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.repo.RemoveByToken(ctx, endpointToken, userClaim.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (s *EndpointService) genEndpoint() (string, error) {
	endpoint, err := token.GenerateEndpointToken(endpointLength)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}
