package endpoint

import (
	"context"
	"errors"
	"ohp/internal/pkg/token"
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

func (s *EndpointService) List(ctx context.Context) ([]Endpoint, error) {
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByUserID(ctx, userClaim.UserID)
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

func (s *EndpointService) genEndpoint() (string, error) {
	endpoint, err := token.GenerateEndpointToken(endpointLength)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}
