package services

import (
	"errors"
	"time"

	"game/internal/domain/entities"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtProvider struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJwtProvider(secret string, accessTTL, refreshTTL time.Duration) *JwtProvider {
	return &JwtProvider{secretKey: []byte(secret), accessTokenTTL: accessTTL, refreshTokenTTL: refreshTTL}
}

type userClaims struct {
	UserID string `json:"uid"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

func (p *JwtProvider) GenerateAccessToken(user *entities.User) (string, error) {
	return p.generateToken(user.UUID, "access", p.accessTokenTTL)
}

func (p *JwtProvider) GenerateRefreshToken(user *entities.User) (string, error) {
	return p.generateToken(user.UUID, "refresh", p.refreshTokenTTL)
}

func (p *JwtProvider) generateToken(userID, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := userClaims{
		UserID: userID,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}

func (p *JwtProvider) ValidateAccessToken(tokenString string) (string, error) {
	uid, typ, err := p.parseAndValidate(tokenString)
	if err != nil {
		return "", err
	}
	if typ != "access" {
		return "", errors.New("invalid token type")
	}
	return uid, nil
}

func (p *JwtProvider) ValidateRefreshToken(tokenString string) (string, error) {
	uid, typ, err := p.parseAndValidate(tokenString)
	if err != nil {
		return "", err
	}
	if typ != "refresh" {
		return "", errors.New("invalid token type")
	}
	return uid, nil
}

func (p *JwtProvider) GetUserID(tokenString string) (string, error) {
	uid, _, err := p.parseAndValidate(tokenString)
	return uid, err
}

func (p *JwtProvider) parseAndValidate(tokenString string) (string, string, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return p.secretKey, nil
	})
	if err != nil {
		return "", "", err
	}
	if !parsed.Valid {
		return "", "", errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(*userClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}
	return claims.UserID, claims.Type, nil
}
