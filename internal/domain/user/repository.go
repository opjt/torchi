package user

import (
	"context"
	db "ohp/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type userRepository struct {
	queries *db.Queries
}
type UserRepository interface {
	UpsertUserByEmail(context.Context, string) (*User, error)
	FindByID(context.Context, uuid.UUID) (*User, error)
	TermsAgree(context.Context, uuid.UUID) error
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}
func (r *userRepository) TermsAgree(ctx context.Context, userID uuid.UUID) error {

	return r.queries.UpdateUserTermsAgreed(ctx, userID)
}
func (r *userRepository) UpsertUserByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: github oauth에서 email도 변경될 수가 있음
	user, err := r.queries.UpsertUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	entity := &User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		TermsAgreed: user.TermsAgreed,
	}
	return entity, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := r.queries.FindUserById(ctx, id)
	if err != nil {
		if db.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}

	entity := &User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		TermsAgreed: user.TermsAgreed,
	}
	return entity, nil
}
