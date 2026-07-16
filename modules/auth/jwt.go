package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	Secret []byte
	Issuer string
	AccessTokenTTL time.Duration
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	RefreshTokenHash []byte
}

func NewJWTManager(secret string) (*JWTManager, error) {

	if secret == "" {
		return nil, errors.New("JWT-SECRET is empty")
	}

	return &JWTManager{
		Secret: []byte(secret),
		Issuer: JwT_ISSUER,
		AccessTokenTTL: JWT_ACCESS_TOKEN_TTL,
	}, nil
}

// GenerateAccessToken creates a signed JWT access token.
func (j *JWTManager) GenerateAccessToken(userID uuid.UUID) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.Secret)
}

// ParseAccessToken validates and parses the JWT.
func (j *JWTManager) ParseAccessToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Signing Method is not matching")
			}
			return j.Secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

// GenerateRefreshToken creates a cryptographically secure random refresh token.
func (j *JWTManager) GenerateRefreshToken() (string, error) {

	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// HashRefreshToken hashes the refresh token before storing it.
func (j *JWTManager) HashRefreshToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

// GenerateTokenPair generates an access token, refresh token,
// and the hashed refresh token ready for database storage.
func (j *JWTManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {

	accessToken, err := j.GenerateAccessToken(userID)
	if err != nil {
		return nil,err
	}

	refreshToken, err := j.GenerateRefreshToken()
	if err != nil {
		return nil,err
	}

	refreshTokenHash := j.HashRefreshToken(refreshToken)

	return &TokenPair{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		RefreshTokenHash: refreshTokenHash,
	}, nil
}
