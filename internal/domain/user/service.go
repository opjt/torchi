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

func (s *UserService) UpsertUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.UpsertUserByEmail(ctx, email)
}

func (s *UserService) FindByEmail(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
