package notifications

import (
	"context"
	db "torchi/internal/infrastructure/db/postgresql"

	"github.com/google/uuid"
)

type notiRepository struct {
	queries *db.Queries
}

type NotiRepository interface {
	Create(context.Context, Noti) (Noti, error)
	InsertMute(context.Context, Noti) (Noti, error)
	UpdateStatus(context.Context, Noti) error
	GetWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32, endpointID *uuid.UUID, query *string) ([]Noti, error)
	MarkAsReadBefore(ctx context.Context, userID uuid.UUID, lastID uuid.UUID, endpointID *uuid.UUID) error
	MarkDelete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	SaveReaction(ctx context.Context, notiID uuid.UUID, reaction string) error
}

func NewNotiRepository(queries *db.Queries) NotiRepository {
	return &notiRepository{
		queries: queries,
	}
}

func (r *notiRepository) InsertMute(ctx context.Context, noti Noti) (Noti, error) {
	statusStr := string(noti.Status)
	createdRow, err := r.queries.CreateMuteNotification(ctx, db.CreateMuteNotificationParams{
		UserID:  noti.UserID,
		Body:    noti.Body,
		Status:  &statusStr,
		ID:      *noti.EndpointID,
		Actions: noti.Actions,
	})
	entity := Noti{
		ID:         createdRow.ID,
		EndpointID: createdRow.EndpointID,
		Body:       createdRow.Body,
	}
	return entity, err

}

func (r *notiRepository) MarkDelete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {

	return r.queries.MarkDeleteNotificationByID(ctx, db.MarkDeleteNotificationByIDParams{
		UserID: userID,
		ID:     id,
	})
}
func (r *notiRepository) MarkAsReadBefore(ctx context.Context, userID uuid.UUID, lastID uuid.UUID, endpointID *uuid.UUID) error {
	return r.queries.MarkNotificationsAsReadBefore(ctx, db.MarkNotificationsAsReadBeforeParams{
		UserID:     userID,
		ID:         lastID,
		EndpointID: endpointID,
	})
}

func (r *notiRepository) GetWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32, endpointID *uuid.UUID, query *string) ([]Noti, error) {

	params := db.GetNotificationsWithCursorParams{
		UserID:     userID,
		Limit:      limit,
		LastID:     lastID,
		EndpointID: endpointID,
		Query:      query,
	}

	rows, err := r.queries.GetNotificationsWithCursor(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []Noti
	for _, row := range rows {
		var s notiStatus
		if row.Status != nil {
			s = notiStatus(*row.Status)
		}

		result = append(result, Noti{
			ID:           row.ID,
			EndpointID:   row.EndpointID,
			EndpointName: row.EndpointName,
			UserID:       row.UserID,
			Body:         row.Body,
			Status:       s,
			ReadAt:       row.ReadAt,
			CreatedAt:    row.CreatedAt,
			Actions:      row.Actions,
			Reaction:     row.Reaction,
			ReactionAt:   row.ReactionAt,
		})
	}

	return result, nil
}

func (r *notiRepository) Create(ctx context.Context, noti Noti) (Noti, error) {
	createdRow, err := r.queries.CreateNotification(ctx, db.CreateNotificationParams{
		ID:      *noti.EndpointID,
		Body:    noti.Body,
		UserID:  noti.UserID,
		Actions: noti.Actions,
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

func (r *notiRepository) SaveReaction(ctx context.Context, notiID uuid.UUID, reaction string) error {
	return r.queries.SaveReaction(ctx, db.SaveReactionParams{
		ID:       notiID,
		Reaction: &reaction,
	})
}
