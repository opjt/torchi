package wrapper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"torchi/internal/domain/common"
)

type ErrorDetail struct {
	Code string `json:"code"`
}
type APIResponse struct {
	Code    int          `json:"code"`
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// RespondJSON 성공 응답 전용
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Code:    status,
		Success: true,
		Data:    data,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// RespondError 에러 응답 전용. DomainError면 정의된 status/code 사용, 아니면 500
func RespondError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusInternalServerError
	code := "INTERNAL_SERVER_ERROR"

	var domErr *common.DomainError
	if errors.As(err, &domErr) {
		status = domErr.HttpStatus
		code = domErr.Code
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Code:  status,
		Error: &ErrorDetail{Code: code},
	})
}

func WrapJson[T any](
	handler func(context.Context, T) (any, error),
	logger func(msg string, keyvals ...any),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto T
		if r.Body != http.NoBody {
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				logger("parse json error", "err", err)
				RespondError(w, common.ErrBadRequest)
				return
			}
		}
		res, err := handler(r.Context(), dto)
		if err != nil {
			logger("handler error", "err", err)
			RespondError(w, err)
			return
		}
		RespondJSON(w, http.StatusOK, res)
	}
}
