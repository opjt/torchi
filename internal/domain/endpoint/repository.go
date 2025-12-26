package endpoint

import (
	"context"
	"errors"
	db "ohp/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

var ErrDuplicateToken = errors.New("duplicate endpoint token")

type EndpointRepository interface {
	Add(ctx context.Context, params insertEndpointParams) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Endpoint, error)
}

type endpointRepository struct {
	queries *db.Queries
}

type insertEndpointParams struct {
	userID      uuid.UUID
	serviceName string
	endpoint    string
}

func NewEndpointRepository(queries *db.Queries) EndpointRepository {
	return endpointRepository{
		queries: queries,
	}
}
func (r endpointRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Endpoint, error) {
	endpoints, err := r.queries.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []Endpoint
	for _, endpoint := range endpoints {
		result = append(result, Endpoint{
			ID:        endpoint.ID,
			Name:      endpoint.Name,
			Endpoint:  endpoint.Endpoint,
			CreatedAt: endpoint.CreatedAt,
		})
	}
	return result, nil
}
func (r endpointRepository) Add(ctx context.Context, params insertEndpointParams) error {
	_, err := r.queries.CreateEndpoint(ctx, db.CreateEndpointParams{
		UserID:   params.userID,
		Name:     params.serviceName,
		Endpoint: params.endpoint,
	})

	if err != nil {
		if db.IsUniqueViolation(err) {
			return ErrDuplicateToken
		}
		return err

	}
	return nil
}
