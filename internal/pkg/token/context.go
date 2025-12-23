package token

import (
	"context"
	"errors"
)

type contextKey string

const UserContextKey contextKey = "user"

var ErrUnauthenticated = errors.New("unauthenticated")

func UserFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, ErrUnauthenticated
	}
	return claims, nil
}

func ContextWith(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}
