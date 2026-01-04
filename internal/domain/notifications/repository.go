package notifications

import (
	"context"
	db "ohp/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type notiRepository struct {
	queries *db.Queries
}

type NotiRepository interface {
	Create(context.Context, Noti) (Noti, error)
	UpdateStatus(context.Context, Noti) error
	GetList(context.Context, uuid.UUID) ([]NotiWithEndpoint, error)
	GetWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32) ([]NotiWithEndpoint, error)
	MarkAsReadBefore(ctx context.Context, userID uuid.UUID, lastID uuid.UUID) error
}

func NewNotiRepository(queries *db.Queries) NotiRepository {
	return &notiRepository{
		queries: queries,
	}
}

func (r *notiRepository) MarkAsReadBefore(ctx context.Context, userID uuid.UUID, lastID uuid.UUID) error {
	return r.queries.MarkNotificationsAsReadBefore(ctx, db.MarkNotificationsAsReadBeforeParams{
		UserID: userID,
		ID:     lastID,
	})
}

func (r *notiRepository) GetWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32) ([]NotiWithEndpoint, error) {

	params := db.GetNotificationsWithCursorParams{
		UserID: userID,
		Limit:  limit,
		LastID: lastID,
	}

	rows, err := r.queries.GetNotificationsWithCursor(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []NotiWithEndpoint
	for _, row := range rows {
		var s notiStatus
		if row.Status != nil {
			s = notiStatus(*row.Status)
		}

		result = append(result, NotiWithEndpoint{
			Noti: Noti{
				ID:         row.ID,
				EndpointID: row.EndpointID,
				UserID:     row.UserID,
				Body:       row.Body,
				Status:     s,
				IsRead:     row.IsRead,
				ReadAt:     row.ReadAt,
				CreatedAt:  row.CreatedAt,
			},
			EndpointInfo: EndpointInfo{
				Name: row.EndpointName,
			},
		})
	}

	return result, nil
}

func (r *notiRepository) GetList(ctx context.Context, userID uuid.UUID) ([]NotiWithEndpoint, error) {
	notis, err := r.queries.FindNotificationByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []NotiWithEndpoint
	for _, noti := range notis {
		result = append(result, NotiWithEndpoint{
			Noti: Noti{
				ID:         noti.ID,
				EndpointID: noti.EndpointID,
				Body:       noti.Body,
				Status:     notiStatus(*noti.Status),
				IsRead:     noti.IsRead,
				CreatedAt:  noti.CreatedAt,
				ReadAt:     noti.ReadAt,
				IsDeleted:  noti.IsDeleted,
			},
			EndpointInfo: EndpointInfo{
				Name: noti.EndpointName,
			},
		})
	}
	return result, nil
}
func (r *notiRepository) Create(ctx context.Context, noti Noti) (Noti, error) {
	createdRow, err := r.queries.CreateNotification(ctx, db.CreateNotificationParams{
		EndpointID: noti.EndpointID,
		Body:       noti.Body,
		UserID:     noti.UserID,
	})
	entity := Noti{
		ID:         createdRow.ID,
		EndpointID: createdRow.EndpointID,
		Body:       createdRow.Body,
	}
	return entity, err
}

func (r *notiRepository) UpdateStatus(ctx context.Context, noti Noti) error {
	status := string(noti.Status)
	err := r.queries.UpdateStatusNotification(ctx, db.UpdateStatusNotificationParams{
		ID:     noti.ID,
		Status: &status,
	})
	return err
}
