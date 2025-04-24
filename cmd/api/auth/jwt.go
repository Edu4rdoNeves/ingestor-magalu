package auth

import (
	"errors"
	"fmt"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"strings"
	"time"
)

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() *jwtService {
	return &jwtService{
		secretKey: env.JwtSecret,
		issuer:    env.JwtIssuer,
	}
}

type Claim struct {
	Sum uint `json:"sum"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(id uint) (*dto.LoginAuth, error) {
	expiresAt := time.Now().Add(time.Hour * 4).Unix()

	claim := &Claim{
		id,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	loginAuth := dto.LoginAuth{
		Token:     t,
		ExpiresAt: expiresAt,
	}

	return &loginAuth, nil
}

func (s *jwtService) ValidateToken(token string) (*Claim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claim{}, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %v", token)
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claim)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *jwtService) ExtractUserIDFromJwtClaims(context *gin.Context) (*uint, error) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == constants.Empity {
		return nil, errors.New("authorization header missing")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := s.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	userID := claims.Sum

	return &userID, nil
}
