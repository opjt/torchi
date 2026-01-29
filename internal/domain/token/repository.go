package token

import (
	"context"
	db "torchi/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type TokenRepository interface {
	UpsertToken(ctx context.Context, token Token) (uuid.UUID, error)
	RemoveToken(ctx context.Context, token Token) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Token, error)
	FindByEndpoint(ctx context.Context, endpoint string) (*Token, error)
	DeactivateToken(ctx context.Context, endpoint string) error
}

type tokenRepository struct {
	queries *db.Queries
}

func NewTokenRepository(queries *db.Queries) TokenRepository {
	return &tokenRepository{
		queries: queries,
	}
}

func (r *tokenRepository) FindByEndpoint(ctx context.Context, endpoint string) (*Token, error) {

	token, err := r.queries.FindTokenByEndpoint(ctx, endpoint)
	if err != nil {
		if db.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &Token{
		ID:       token.ID,
		UserID:   token.UserID,
		P256dh:   token.P256dhKey,
		Auth:     token.AuthKey,
		EndPoint: token.Endpoint,
	}, nil
}
func (r *tokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Token, error) {

	tokens, err := r.queries.FindTokenByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []Token
	for _, token := range tokens {
		result = append(result, Token{
			ID:       token.ID,
			UserID:   token.UserID,
			P256dh:   token.P256dhKey,
			Auth:     token.AuthKey,
			EndPoint: token.Endpoint,
		})
	}
	return result, nil
}
func (r *tokenRepository) UpsertToken(ctx context.Context, token Token) (uuid.UUID, error) {

	param := db.UpsertTokenParams{
		UserID:    token.UserID,
		P256dhKey: token.P256dh,
		AuthKey:   token.Auth,
		Endpoint:  token.EndPoint,
	}
	return r.queries.UpsertToken(ctx, param)
}

func (r *tokenRepository) RemoveToken(ctx context.Context, token Token) error {
	param := db.DeleteTokenParams{
		Endpoint:  token.EndPoint,
		P256dhKey: token.P256dh,
		AuthKey:   token.Auth,
	}
	return r.queries.DeleteToken(ctx, param)
}

func (r *tokenRepository) DeactivateToken(ctx context.Context, endpoint string) error {
	return r.queries.DeactivatePushToken(ctx, endpoint)
}
