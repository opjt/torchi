package handler

import (
	"context"
	"net/http"
	"ohp/internal/api/wrapper"
	"ohp/internal/domain/user"
	"ohp/internal/pkg/config"
	"ohp/internal/pkg/log"
	"ohp/internal/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	log      *log.Logger
	frontUrl string
	service  *user.UserService
}

func NewUserHandler(log *log.Logger, env config.Env, service *user.UserService) *UserHandler {
	return &UserHandler{
		log:      log,
		frontUrl: env.FrontUrl,
		service:  service,
	}
}
func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/whoami", h.Whoami)
	r.Post("/terms-agree", wrapper.WrapJson(h.TermsAgree, h.log.Error, wrapper.RespondJSON))

	return r
}
func (h *UserHandler) TermsAgree(ctx context.Context, _ interface{}) (interface{}, error) {

	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.service.TermsAgree(ctx, userClaim.UserID)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

type resWhoami struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	TermsAgreed bool      `json:"terms_agreed"`
}

func (h *UserHandler) Whoami(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userClaim, err := token.UserFromContext(ctx)
	if err != nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.service.FindByEmail(ctx, userClaim.UserID)
	if err != nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, "user not found")
		return
	}
	h.log.Debug("...", "user", user)

	resp := resWhoami{
		UserID:      user.ID,
		Email:       user.Email,
		TermsAgreed: user.TermsAgreed,
	}
	wrapper.RespondJSON(w, http.StatusOK, resp)
}
