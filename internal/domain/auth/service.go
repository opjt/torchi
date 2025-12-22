package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ohp/internal/domain/user"
	"ohp/internal/pkg/config"
	"ohp/internal/pkg/token"
	"time"
)

type AuthService struct {
	githubConfig  config.Github
	userService   *user.UserService
	tokenProvider *token.TokenProvider
}

func NewAuthService(
	env config.Env,
	userService *user.UserService,
) *AuthService {

	tokenProvider := token.NewTokenProvider(env.JWTSecret, "ohp-api", 24*time.Hour) // TODO: env 전환 필요
	return &AuthService{
		githubConfig:  env.Github,
		userService:   userService,
		tokenProvider: tokenProvider,
	}
}

// TODO: 플랫폼별 분리가 필요함(Github, ...)
type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *AuthService) OauthGithubFlow(ctx context.Context, code string) (token string, err error) {
	accessToken, err := s.githubGetAccessToken(code)
	if err != nil {
		return token, err
	}
	githubUser, err := s.getGithubUserProfile(accessToken)
	if err != nil {
		return token, err
	}
	dbUser, err := s.userService.UpsertUserByEmail(ctx, githubUser.Email)
	if err != nil {
		return token, err
	}

	token, err = s.tokenProvider.Create(dbUser.ID.String(), dbUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate service token: %w", err)
	}

	return token, nil
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
