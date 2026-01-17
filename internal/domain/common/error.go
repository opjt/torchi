package common

// DomainError는 서비스 레이어에서 발생하여 핸들러로 전달되는 커스텀 에러입니다.
type DomainError struct {
	Code       string
	Message    string
	HttpStatus int
}

func (e *DomainError) Error() string {
	return e.Message
}

// 헬퍼 함수: 에러 생성을 간편하게 함
func NewError(status int, code string, message string) *DomainError {
	return &DomainError{
		HttpStatus: status,
		Code:       code,
		Message:    message,
	}
}

// 프론트와 한 약속.
var (
	ErrBadRequest          = NewError(400, "BAD_REQUEST", "잘못된 요청입니다.")
	ErrUnauthorized        = NewError(401, "UNAUTHORIZED", "인증이 필요합니다.")
	ErrInternalServerError = NewError(500, "INTERNAL_SERVER_ERROR", "서버 내부 오류가 발생했습니다.")

	ErrUserNotFound     = NewError(404, "USER_NOT_FOUND", "사용자를 찾을 수 없습니다.")
	ErrPushTokenExpired = NewError(400, "PUSH_TOKEN_EXPIRED", "푸시 토큰이 만료되었습니다.")
)
