package notifications

import (
	"context"
	db "ohp/internal/infrastructure/db/postgresql"
)

type notiRepository struct {
	queries *db.Queries
}

type NotiRepository interface {
	Create(context.Context, Noti) error
}

func NewNotiRepository(queries *db.Queries) NotiRepository {
	return &notiRepository{
		queries: queries,
	}
}

func (r *notiRepository) Create(ctx context.Context, noti Noti) error {
	_, err := r.queries.CreateNotification(ctx, db.CreateNotificationParams{
		ServiceID: noti.ServiceID,
		Body:      noti.Body,
	})
	return err
}
