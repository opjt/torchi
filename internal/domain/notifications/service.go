package notifications

import (
	"context"
	"time"
	"torchi/internal/domain/sse"

	"github.com/google/uuid"
)

type NotiService struct {
	repo      NotiRepository
	sseBroker *sse.Broker
}

func NewNotiService(
	repo NotiRepository,
	sseBroker *sse.Broker,
) *NotiService {
	return &NotiService{
		repo:      repo,
		sseBroker: sseBroker,
	}
}

func (s *NotiService) MarkDelete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.MarkDelete(ctx, userID, id)
}
func (s *NotiService) MarkAllAsRead(ctx context.Context, userID uuid.UUID, lastID uuid.UUID, endpointID *uuid.UUID) error {
	return s.repo.MarkAsReadBefore(ctx, userID, lastID, endpointID)
}

func (s *NotiService) GetListWithCursor(ctx context.Context, userID uuid.UUID, lastID *uuid.UUID, limit int32, endpointID *uuid.UUID, query *string) ([]Noti, error) {

	return s.repo.GetWithCursor(ctx, userID, lastID, limit, endpointID, query)
}

type ReqRegister struct {
	EndpointID         uuid.UUID
	UserID             uuid.UUID
	Body               string
	NotificationEnable bool
	Actions            []string
}

func (s *NotiService) Register(ctx context.Context, req ReqRegister) (noti Noti, err error) {
	newNoti := Noti{
		EndpointID: &req.EndpointID,
		Body:       req.Body,
		UserID:     req.UserID,
		Actions:    req.Actions,
	}

	if !req.NotificationEnable {
		newNoti.Status = notiStatusMute
		now := time.Now()
		newNoti.ReadAt = &now
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

func (s *NotiService) UpdateStatusTimeout(ctx context.Context, reqID uuid.UUID) error {
	return s.updateStatus(ctx, ReqUpdateStatus{
		ID:     reqID,
		Status: notiStatusTimeoutReply,
	})
}

func (s *NotiService) SaveReaction(ctx context.Context, notiID uuid.UUID, reaction string) error {
	// TODO(low): 이미 타임아웃 난 거는 액션할필요 없음.
	return s.repo.SaveReaction(ctx, notiID, reaction)
}
