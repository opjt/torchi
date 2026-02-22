package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenProvider struct {
	secret        []byte
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewTokenProvider(secret string, issuer string, accessExp time.Duration, refreshExp time.Duration) *TokenProvider {
	return &TokenProvider{
		secret:        []byte(secret),
		issuer:        issuer,
		accessExpiry:  accessExp,
		refreshExpiry: refreshExp,
	}
}

// CreateAccessToken: 짧은 수명의 액세스 토큰 생성
func (p *TokenProvider) CreateAccessToken(userID uuid.UUID, email string) (string, error) {
	return p.generateToken(userID, email, p.accessExpiry)
}

// CreateRefreshToken: 긴 수명의 리프레시 토큰 생성
func (p *TokenProvider) CreateRefreshToken(userID uuid.UUID, email string) (string, error) {
	return p.generateToken(userID, email, p.refreshExpiry)
}

func (p *TokenProvider) CreatePairToken(userID uuid.UUID, email string) (string, string, error) {
	access, err := p.CreateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}
	refresh, err := p.CreateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// 내부 공통 토큰 생성 로직
func (p *TokenProvider) generateToken(userID uuid.UUID, email string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    p.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

// Validate: 액세스/리프레시 구분 없이 서명과 유효기간 검증
func (p *TokenProvider) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
