package token

import "context"

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

// user token 삭제
