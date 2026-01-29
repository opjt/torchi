package token

import (
	"context"

	"github.com/google/uuid"
)

type TokenService struct {
	repo TokenRepository
}

func NewTokenService(repo TokenRepository) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

// user token 추가
func (s *TokenService) Register(ctx context.Context, token Token) error {
	_, err := s.repo.UpsertToken(ctx, token)

	return err
}

// token delete
func (s *TokenService) Unregister(ctx context.Context, token Token) error {
	return s.repo.RemoveToken(ctx, token)
}

func (s *TokenService) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Token, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *TokenService) FindByEndpoint(ctx context.Context, endpoint string) (*Token, error) {
	return s.repo.FindByEndpoint(ctx, endpoint)
}

func (s *TokenService) DeactiveToken(ctx context.Context, endpoint string) error {
	return s.repo.DeactivateToken(ctx, endpoint)
}
