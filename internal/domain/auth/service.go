package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"torchi/internal/domain/user"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/token"

	"github.com/google/uuid"
)

type AuthService struct {
	githubConfig  config.Github
	userService   *user.UserService
	tokenProvider *token.TokenProvider
}

func NewAuthService(
	env config.Env,
	userService *user.UserService,
	tokenProvider *token.TokenProvider,
) *AuthService {

	return &AuthService{
		githubConfig:  env.Github,
		userService:   userService,
		tokenProvider: tokenProvider,
	}
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	// 서명 검증 및 클레임 추출 (DB 조회 없음)
	claims, err := s.tokenProvider.Validate(refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 추출한 정보로 새 토큰 세트 생성 (Sliding Window)
	newAccess, err := s.tokenProvider.CreateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := s.tokenProvider.CreateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *AuthService) TestLogin(ctx context.Context) (string, string, error) {

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	email := "tester@gmail.com"

	// Access Token (짧은 수명)
	accessToken, err := s.tokenProvider.CreateAccessToken(userID, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh Token (긴 수명)
	refreshToken, err := s.tokenProvider.CreateRefreshToken(userID, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return accessToken, refreshToken, nil

}

type LoginResult struct {
	UserID uuid.UUID
	AT     string
	RT     string
}

func (s *AuthService) GuestLogin(ctx context.Context, guestID *uuid.UUID) (LoginResult, error) {
	dbUser, err := s.userService.UpsertGuestUser(ctx, guestID)
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to upsert guest user: %w", err)
	}

	at, rt, err := s.tokenProvider.CreatePairToken(dbUser.ID, "guest")
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		UserID: dbUser.ID,
		AT:     at,
		RT:     rt,
	}, nil
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *AuthService) OauthGithubFlow(ctx context.Context, code string) (at string, rt string, err error) {
	githubAccessToken, err := s.githubGetAccessToken(code)
	if err != nil {
		return at, rt, err
	}
	githubUser, err := s.getGithubUserProfile(githubAccessToken)
	if err != nil {
		return at, rt, err
	}
	dbUser, err := s.userService.UpsertUserByEmail(ctx, githubUser.Email) // TODO: 이메일이 없을 수도 있음
	if err != nil {
		return at, rt, err
	}

	var email string
	if dbUser.Email != nil {
		email = *dbUser.Email
	}

	// 페어 토큰 생성
	at, err = s.tokenProvider.CreateAccessToken(dbUser.ID, email)
	if err != nil {
		return at, rt, err
	}
	rt, err = s.tokenProvider.CreateRefreshToken(dbUser.ID, email)
	if err != nil {
		return at, rt, err
	}

	return at, rt, nil
}
func (s *AuthService) githubGetAccessToken(code string) (string, error) {
	// 요청 데이터 준비
	reqBody := map[string]string{
		"client_id":     s.githubConfig.ClientID,
		"client_secret": s.githubConfig.ClientSecret,
		"code":          code,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// GitHub Token 엔드포인트에 POST 요청
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonBody))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

type githubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"` // 유저 아이디
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (s *AuthService) getGithubUserProfile(accessToken string) (*githubUser, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var user githubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
