package notifications

import (
	"context"

	"github.com/google/uuid"
)

type NotiService struct {
	repo NotiRepository
}

func NewNotiService(
	repo NotiRepository,
) *NotiService {
	return &NotiService{
		repo: repo,
	}
}

func (s *NotiService) MarkDelete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.MarkDelete(ctx, userID, id)
}
func (s *NotiService) MarkAllAsRead(ctx context.Context, userID uuid.UUID, lastID uuid.UUID) error {
	return s.repo.MarkAsReadBefore(ctx, userID, lastID)
}

func (s *NotiService) GetListWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32, endpointID *uuid.UUID) ([]Noti, error) {

	return s.repo.GetWithCursor(ctx, userID, lastID, limit, endpointID)
}

type ReqRegister struct {
	EndpointID         uuid.UUID
	UserID             uuid.UUID
	Body               string
	NotificationEnable bool
}

func (s *NotiService) Register(ctx context.Context, req ReqRegister) (noti Noti, err error) {
	newNoti := Noti{
		EndpointID: &req.EndpointID,
		Body:       req.Body,
		UserID:     req.UserID,
	}

	if req.NotificationEnable == false {
		newNoti.Status = notiStatusMute
		return s.repo.InsertMute(ctx, newNoti)
	}
	return s.repo.Create(ctx, newNoti)
}

type ReqUpdateStatus struct {
	ID     uuid.UUID
	Status notiStatus
}

func (s *NotiService) updateStatus(ctx context.Context, req ReqUpdateStatus) error {
	err := s.repo.UpdateStatus(ctx, Noti{
		ID:     req.ID,
		Status: req.Status,
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *NotiService) UpdateStatusSent(ctx context.Context, reqID uuid.UUID) error {
	return s.updateStatus(ctx, ReqUpdateStatus{
		ID:     reqID,
		Status: notiStatusSent,
	})
}
