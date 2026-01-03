package notifications

import "context"

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

func (s *NotiService) Register(ctx context.Context, noti Noti) error {
	return nil
}
