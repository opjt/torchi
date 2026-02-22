package handler

import (
	"encoding/json"
	"net/http"
	"torchi/internal/api/wrapper"
	"torchi/internal/domain/auth"
	"torchi/internal/domain/common"
	"torchi/internal/domain/token"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	log          *log.Logger
	frontUrl     string
	service      *auth.AuthService
	tokenService *token.TokenService
	env          config.Env
}

const (
	AccessCookieKey  = "access_token"
	RefreshCookieKey = "refresh_token"
)

func NewAuthHandler(
	log *log.Logger,
	env config.Env,
	service *auth.AuthService,
	tokenService *token.TokenService,
) *AuthHandler {

	return &AuthHandler{
		log:          log,
		frontUrl:     env.FrontUrl,
		env:          env,
		service:      service,
		tokenService: tokenService,
	}
}
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/github/callback", h.OauthGithubCallback)
	r.Post("/logout", h.Logout)
	r.Post("/refresh", h.Refresh)
	r.Post("/guest", h.GuestLogin)

	if h.env.Stage == config.StageDev {
		r.Get("/fake/login", h.FakeLogin)
	}

	return r
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// 쿠키에서 refresh_token 추출
	cookie, err := r.Cookie(RefreshCookieKey)
	if err != nil {
		h.log.Error("refresh cookie missing", "error", err)
		http.Error(w, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	// 서비스 레이어 호출 (Stateless 검증 및 새 토큰 생성)
	at, rt, err := h.service.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		h.log.Error("failed to refresh token", "error", err)
		// 토큰이 만료되었거나 변조된 경우 401을 내려주어 프론트에서 재로그인 유도
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// 새로운 Access/Refresh 토큰을 쿠키에 설정 (Sliding Window)
	h.setAuthCookies(w, at, rt)

	// 성공 응답
	wrapper.RespondJSON(w, http.StatusOK, nil)
}

// 쿠키 설정을 위한 헬퍼 함수
func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessCookieKey,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !config.IsDev(h.env.Stage),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 3, // 3시간
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshCookieKey,
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   !config.IsDev(h.env.Stage),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 30, // 30일
	})
}

func (h *AuthHandler) FakeLogin(w http.ResponseWriter, r *http.Request) {
	at, rt, err := h.service.TestLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, at, rt)
	wrapper.RespondJSON(w, http.StatusOK, nil)
}

func (h *AuthHandler) GuestLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID *string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		wrapper.RespondJSON(w, http.StatusBadRequest, err)
		return
	}

	var userID *uuid.UUID
	if req.UserID != nil {
		id, err := uuid.Parse(*req.UserID)
		if err != nil {
			wrapper.RespondJSON(w, http.StatusBadRequest, common.ErrBadRequest)
			return
		}
		userID = &id
	}

	loginResult, err := h.service.GuestLogin(r.Context(), userID)
	if err != nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.setAuthCookies(w, loginResult.AT, loginResult.RT)
	wrapper.RespondJSON(w, http.StatusOK, map[string]string{
		"user_id": loginResult.UserID.String(),
	})
}

type LogoutRequest struct {
	Endpoint string `json:"endpoint"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.clearAuthCookies(w)
	h.tryDeactivateToken(r)
	wrapper.RespondJSON(w, http.StatusOK, nil)
}

func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	cookies := map[string]string{
		AccessCookieKey:  "/",
		RefreshCookieKey: "/auth/refresh",
	}
	for key, path := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     key,
			Value:    "",
			Path:     path,
			HttpOnly: true,
			MaxAge:   -1,
		})
	}
}

func (h *AuthHandler) tryDeactivateToken(r *http.Request) {
	if r.Body == http.NoBody {
		return
	}

	var reqBody LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Endpoint == "" {
		return
	}

	h.log.Info("deactive token", "endpoint", reqBody.Endpoint)
	if err := h.tokenService.DeactiveToken(r.Context(), reqBody.Endpoint); err != nil {
		h.log.Error("failed to deactive token", "error", err)
	}
}
func (h *AuthHandler) OauthGithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.frontUrl+"/login", http.StatusFound)
		return
	}

	//  GitHub API로 사용자 정보 가져오기
	at, rt, err := h.service.OauthGithubFlow(r.Context(), code)
	if err != nil {
		h.log.Error("failed to get user profile", "error", err)
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, at, rt)
	http.Redirect(w, r, h.frontUrl+"/app", http.StatusFound)
}
