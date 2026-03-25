package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (s *UserService) Withdraw(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByID(ctx, userID)
}
func (s *UserService) TermsAgree(ctx context.Context, userID uuid.UUID) error {
	return s.repo.TermsAgree(ctx, userID)
}

func (s *UserService) UpsertByProvider(ctx context.Context, provider, providerID string, email *string) (*User, error) {
	return s.repo.UpsertUserByProvider(ctx, provider, providerID, email)
}

func (s *UserService) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) UpsertGuestUser(ctx context.Context, guestID *uuid.UUID) (*User, error) {
	if guestID != nil {
		user, err := s.repo.FindByID(ctx, *guestID)
		if err == nil && user != nil {
			return user, nil
		}
	}
	return s.repo.InsertGuestUser(ctx)
}
