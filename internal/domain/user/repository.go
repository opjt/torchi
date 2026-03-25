package user

import (
	"context"
	db "torchi/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type userRepository struct {
	queries *db.Queries
}
type UserRepository interface {
	UpsertUserByProvider(ctx context.Context, provider, providerID string, email *string) (*User, error)
	InsertGuestUser(context.Context) (*User, error)
	FindByID(context.Context, uuid.UUID) (*User, error)
	TermsAgree(context.Context, uuid.UUID) error
	DeleteByID(context.Context, uuid.UUID) error
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}
func (r *userRepository) DeleteByID(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteUser(ctx, userID)
}
func (r *userRepository) TermsAgree(ctx context.Context, userID uuid.UUID) error {
	return r.queries.UpdateUserTermsAgreed(ctx, userID)
}

func (r *userRepository) UpsertUserByProvider(ctx context.Context, provider, providerID string, email *string) (*User, error) {
	user, err := r.queries.UpsertUserByProvider(ctx, db.UpsertUserByProviderParams{
		Provider:   &provider,
		ProviderID: &providerID,
		Email:      email,
	})
	if err != nil {
		return nil, err
	}
	return toEntity(user), nil
}

func (r *userRepository) InsertGuestUser(ctx context.Context) (*User, error) {
	user, err := r.queries.InsertGuestUser(ctx)
	if err != nil {
		return nil, err
	}
	return toEntity(user), nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := r.queries.FindUserById(ctx, id)
	if err != nil {
		if db.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}

	return toEntity(user), nil
}

func toEntity(user db.User) *User {
	return &User{
		ID:          user.ID,
		Email:       user.Email,
		Provider:    user.Provider,
		ProviderID:  user.ProviderID,
		CreatedAt:   user.CreatedAt,
		TermsAgreed: user.TermsAgreed,
		IsGuest:     user.Guest,
	}
}
