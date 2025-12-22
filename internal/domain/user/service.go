package user

import "context"

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) UpsertUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.UpsertUserByEmail(ctx, email)
}
