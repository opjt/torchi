package wrapper

import (
	"context"
	"encoding/json"
	"net/http"
	"torchi/internal/domain/common"
)

type ErrorDetail struct {
	// common.DomainError로 대체할지 고민..
	Code    string `json:"code"`
	Message string `json:"message"`
}
type APIResponse struct {
	Code    int          `json:"code"`
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Code:    status,
		Success: status >= 200 && status < 300,
	}

	if payload != nil {
		switch v := payload.(type) {
		case *common.DomainError:
			resp.Error = &ErrorDetail{
				Code:    v.Code,
				Message: v.Message,
			}
		case error: // 일반적인 Go 에러일 때 (예상치 못한 에러)
			resp.Error = &ErrorDetail{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: v.Error(),
			}
		default: // 성공 데이터일 때
			resp.Data = payload
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}
func WrapJson[T any](
	handler func(context.Context, T) (interface{}, error),
	logger func(msg string, keyvals ...any),
	respond func(w http.ResponseWriter, status int, data any),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto T
		if r.Body != http.NoBody {
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				logger("parse json error", "err", err)
				respond(w, http.StatusBadRequest, err)
				return
			}
		}
		res, err := handler(r.Context(), dto)
		if err != nil {
			logger("handler error", "err", err)
			respond(w, http.StatusBadRequest, err)
			return
		}
		respond(w, http.StatusOK, res)
	}
}
