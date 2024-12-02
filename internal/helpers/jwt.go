package helpers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

const (
	userAccessTokenSecretKey       = "USER_ACCESS_TOKEN_SECRET"
	accessTokenLifetimeMinutesKey  = "ACCESS_TOKEN_LIFETIME_MINUTES"
	refreshTokenSecretKey          = "USER_REFRESH_TOKEN_SECRET"
	refreshTokenLifetimeMinutesKey = "REFRESH_TOKEN_LIFETIME_MINUTES"
)

type Jwt interface {
	CreateAccessToken(claims UserAccessTokenClaims) (*Token, error)
	Verify(accessToken string) (*tokenClaims, error)
}

type jwtStructure struct {
	userAccessTokenSecret       string
	accessTokenLiftimeMinutes   int
}

type Token struct {
	AccessToken           string
	AccessTokenExpiresAt  int64
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

type UserAccessTokenClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

func NewJwt() (*jwtStructure, error) {
	accessTokenLifetimeMinutes, err := strconv.Atoi(os.Getenv(accessTokenLifetimeMinutesKey))
	if err != nil {
		logrus.Fatalf("Error parsing access token lifetime minutes: %s", err)
		return nil, err
	}

	return &jwtStructure{
		userAccessTokenSecret:       os.Getenv(userAccessTokenSecretKey),
		accessTokenLiftimeMinutes:   accessTokenLifetimeMinutes,
	}, nil
}

func (j *jwtStructure) CreateAccessToken(claims UserAccessTokenClaims) (*Token, error) {
	var res Token
	expirationTime := time.Now().Add(time.Duration(j.accessTokenLiftimeMinutes) * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserId,
		"role":    claims.Role,
		"exp":     expirationTime,
	})
	tokenString, err := token.SignedString([]byte(j.userAccessTokenSecret))
	if err != nil {
		logrus.Errorf("Error signing: %s", err)
		return nil, err
	}
	res.AccessToken = tokenString
	res.AccessTokenExpiresAt = expirationTime

	return &res, nil
}

func (j *jwtStructure) Verify(accessToken string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(j.userAccessTokenSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}