package common

// DomainError는 서비스 레이어에서 발생하여 핸들러로 전달되는 커스텀 에러입니다.
type DomainError struct {
	Code       string
	HttpStatus int
}

func (e *DomainError) Error() string {
	return e.Code
}

// 헬퍼 함수: 에러 생성을 간편하게 함
func NewError(status int, code string) *DomainError {
	return &DomainError{
		HttpStatus: status,
		Code:       code,
	}
}

// 프론트와 한 약속.
var (
	ErrBadRequest   = NewError(400, "BAD_REQUEST")
	ErrUnauthorized = NewError(401, "UNAUTHORIZED")
	ErrInvalidParam = NewError(400, "INVALID_PARAMETER")

	ErrUserNotFound     = NewError(404, "USER_NOT_FOUND")
	ErrPushTokenExpired = NewError(400, "PUSH_TOKEN_EXPIRED")
)
