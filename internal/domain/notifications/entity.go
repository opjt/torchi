package notifications

import (
	"time"

	"github.com/google/uuid"
)

/*
notifications

알림이력을 위한 도메인.
- [push]에서 실제 기기에 push를 보낼 때 서비스는 어떤 알림을 보냈는지 알아야함.
*/

type notiStatus string

const (
	notiStatusPending notiStatus = "pending" // default
	notiStatusSent    notiStatus = "sent"    // push 전송 완료
	notiStatusFailed  notiStatus = "failed"  // push 전송 실패
	notiStatusMute    notiStatus = "mute"    // endpoint의 알림 끈 경우
)

type Noti struct {
	ID           uuid.UUID
	EndpointID   *uuid.UUID
	EndpointName string
	UserID       uuid.UUID
	Body         string
	Status       notiStatus
	CreatedAt    time.Time
	ReadAt       *time.Time
	IsDeleted    bool
	Actions      []string
	Reaction     *string
	ReactionAt   *time.Time
}

func (n *Noti) IsMute() bool {
	return n.Status == notiStatusMute
}
