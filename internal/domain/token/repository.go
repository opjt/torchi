package token

import (
	"context"
	db "ohp/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type TokenRepository interface {
	UpsertToken(ctx context.Context, token Token) (uuid.UUID, error)
}

type tokenRepository struct {
	queries *db.Queries
}

func NewTokenRepository(queries *db.Queries) TokenRepository {
	return tokenRepository{
		queries: queries,
	}
}

func (r tokenRepository) UpsertToken(ctx context.Context, token Token) (uuid.UUID, error) {

	param := db.UpsertTokenParams{
		UserID:    token.UserID,
		P256dhKey: token.P256dh,
		AuthKey:   token.Auth,
	}
	return r.queries.UpsertToken(ctx, param)
}
