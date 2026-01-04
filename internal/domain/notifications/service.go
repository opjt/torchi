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
func (s *NotiService) MarkAllAsRead(ctx context.Context, userID uuid.UUID, lastID uuid.UUID) error {
	return s.repo.MarkAsReadBefore(ctx, userID, lastID)
}

func (s *NotiService) GetList(ctx context.Context, userID uuid.UUID) ([]NotiWithEndpoint, error) {
	return s.repo.GetList(ctx, userID)
}
func (s *NotiService) GetListWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32) ([]NotiWithEndpoint, error) {

	return s.repo.GetWithCursor(ctx, userID, lastID, limit)
}

type ReqRegister struct {
	EndpointID uuid.UUID
	UserID     uuid.UUID
	Body       string
}

func (s *NotiService) Register(ctx context.Context, req ReqRegister) (Noti, error) {
	noti, err := s.repo.Create(ctx, Noti{
		EndpointID: req.EndpointID,
		Body:       req.Body,
		UserID:     req.UserID,
	})
	if err != nil {
		return Noti{}, err
	}
	return noti, nil
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
