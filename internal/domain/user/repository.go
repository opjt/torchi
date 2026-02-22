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
	UpsertUserByEmail(context.Context, string) (*User, error)
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
	//cascade 로 관련 테이블 데이터 같이 제거됨.
	// Endpoints, Push_tokens, Notifications
	return r.queries.DeleteUser(ctx, userID)
}
func (r *userRepository) TermsAgree(ctx context.Context, userID uuid.UUID) error {

	return r.queries.UpdateUserTermsAgreed(ctx, userID)
}
func (r *userRepository) UpsertUserByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: github oauth에서 email도 변경될 수가 있음
	user, err := r.queries.UpsertUserByEmail(ctx, &email)
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

func (r *userRepository) InsertGuestUser(ctx context.Context) (*User, error) {
	user, err := r.queries.InsertGuestUser(ctx)
	if err != nil {
		return nil, err
	}

	entity := &User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		TermsAgreed: user.TermsAgreed,
		IsGuest:     user.Guest,
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
		IsGuest:     user.Guest,
	}
	return entity, nil
}
