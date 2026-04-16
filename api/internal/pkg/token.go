package pkg

import (
	"errors"
	"planet/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	tempTokenSecret    []byte
)

type TokenClaims struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// temp_token에 담을 claims
type TempTokenClaims struct {
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
	jwt.RegisteredClaims
}

func InitToken(cfg *config.Config) {
	accessTokenSecret = []byte(cfg.Token.AccessTokenSecret)
	refreshTokenSecret = []byte(cfg.Token.RefreshTokenSecret)
	tempTokenSecret = []byte(cfg.Token.TempTokenSecret)
}

// Access Token 생성 (1시간)
func GenerateAccessToken(username, nickname string) (string, error) {
	claims := TokenClaims{
		Username: username,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(accessTokenSecret)
}

// Refresh Token 생성 (7일)
func GenerateRefreshToken(username, nickname string) (string, error) {
	claims := TokenClaims{
		Username: username,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(refreshTokenSecret)
}

// temp_token 생성
func GenerateTempToken(provider, providerID string) (string, error) {
	claims := TempTokenClaims{
		Provider:   provider,
		ProviderID: providerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)), // 10분
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tempTokenSecret)
}

// parse 공통 함수
func parseToken(tokenString string, secret []byte, claims jwt.Claims) error {
	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return err
	}
	if !t.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// Access Token 파싱 → 미들웨어에서 사용
func ParseAccessToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	if err := parseToken(tokenString, accessTokenSecret, claims); err != nil {
		return nil, err
	}
	return claims, nil
}

// Refresh Token 파싱 → 토큰 재발급 시 사용
func ParseRefreshToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	if err := parseToken(tokenString, refreshTokenSecret, claims); err != nil {
		return nil, err
	}
	return claims, nil
}

// Temp Token 파싱 → OAuth 회원가입 시 사용
func ParseTempToken(tokenString string) (*TempTokenClaims, error) {
	claims := &TempTokenClaims{}
	if err := parseToken(tokenString, tempTokenSecret, claims); err != nil {
		return nil, err
	}
	return claims, nil
}
