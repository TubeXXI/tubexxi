package utils

import (
	"fmt"
	"tubexxi/video-api/internal/entity"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(jwtSecret string, user *entity.User, expiry time.Duration) (string, error) {
	now := time.Now().UTC()
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}
	if user.ID == uuid.Nil {
		return "", fmt.Errorf("user ID cannot be empty")
	}

	claims := entity.JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		RoleID:    user.Role.ID.String(),
		RoleLevel: int(user.Role.Level),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "idtubexxi",
			Subject:   user.ID.String(),
			ID:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
func VerifyJWT(tokenString string, jwtSecret string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify jwt token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsedToken, nil
}
func GenerateCSRFToken(userId, email string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(time.Second * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
func VerifyCSRFToken(tokenString string, jwtSecret string) (*jwt.Token, error) {

	csrfToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected request method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return csrfToken, nil
}
