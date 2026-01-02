package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	UserID  string `json:"userId"`
	IsOwner bool   `json:"isOwner"`
	jwt.RegisteredClaims
}

type Auth struct {
	secretKey []byte
}

func New(secretKey string) (*Auth, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("secret key is too short; use at least 32 characters")
	}
	return &Auth{
		secretKey: []byte(secretKey),
	}, nil
}

func GenerateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateJoinToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateUserID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *Auth) GenerateToken(userID string, isOwner bool) (string, error) {
	now := time.Now() // Capture time once to ensure consistency

	claims := Claims{
		UserID:  userID,
		IsOwner: isOwner,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(365 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),

			// FIX: Allow token to be valid 5 seconds in the past to handle
			// immediate verification or slight clock drift.
			NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *Auth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
