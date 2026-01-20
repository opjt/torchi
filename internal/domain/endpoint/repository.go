package endpoint

import (
	"context"
	"errors"
	"time"
	db "torchi/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

var ErrDuplicateToken = errors.New("duplicate endpoint token")

type EndpointRepository interface {
	Add(ctx context.Context, params insertEndpointParams) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Endpoint, error)
	RemoveByToken(ctx context.Context, token string, userID uuid.UUID) error
	FindByToken(ctx context.Context, token string) (*Endpoint, error)
	UpdateMute(ctx context.Context, token string, disabledAt *time.Time) error
	UpdateUnmute(ctx context.Context, token string) error
}

type endpointRepository struct {
	queries *db.Queries
}

func NewEndpointRepository(queries *db.Queries) EndpointRepository {
	return &endpointRepository{
		queries: queries,
	}
}

func (r *endpointRepository) UpdateMute(ctx context.Context, token string, disabledAt *time.Time) error {
	return r.queries.UpdateEndpointMute(ctx, db.UpdateEndpointMuteParams{
		Token:                  token,
		NotificationDisabledAt: disabledAt,
	})
}
func (r *endpointRepository) UpdateUnmute(ctx context.Context, token string) error {
	return r.queries.UpdateEndpointUnmute(ctx, token)
}
func (r *endpointRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Endpoint, error) {
	endpoints, err := r.queries.FindEndpointByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []Endpoint
	for _, endpoint := range endpoints {
		result = append(result, Endpoint{
			ID:                 endpoint.ID,
			Name:               endpoint.Name,
			Token:              endpoint.Token,
			CreatedAt:          endpoint.CreatedAt,
			UserID:             endpoint.UserID,
			NotificationEnable: endpoint.NotificationEnabled,
		})
	}
	return result, nil
}
func (r *endpointRepository) FindByToken(ctx context.Context, token string) (*Endpoint, error) {
	rowData, err := r.queries.FindEndpointByToken(ctx, token)

	if err != nil {
		if db.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}

	return &Endpoint{
		ID:                 rowData.ID,
		Name:               rowData.Name,
		Token:              rowData.Token,
		CreatedAt:          rowData.CreatedAt,
		UserID:             rowData.UserID,
		NotificationEnable: rowData.NotificationEnabled,
	}, nil
}

type insertEndpointParams struct {
	userID      uuid.UUID
	serviceName string
	endpoint    string
}

func (r *endpointRepository) Add(ctx context.Context, params insertEndpointParams) error {
	_, err := r.queries.CreateEndpoint(ctx, db.CreateEndpointParams{
		UserID: params.userID,
		Name:   params.serviceName,
		Token:  params.endpoint,
	})

	if err != nil {
		if db.IsUniqueViolation(err) {
			return ErrDuplicateToken
		}
		return err

	}
	return nil
}

func (r *endpointRepository) RemoveByToken(ctx context.Context, token string, userID uuid.UUID) error {
	err := r.queries.DeleteEndpointByToken(ctx, db.DeleteEndpointByTokenParams{
		Token:  token,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	return nil
}
