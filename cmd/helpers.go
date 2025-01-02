package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(userInfo GetAuthInfo) (string, error) {
	claims := jwtClaims{
		UserID:   userInfo.ID,
		Username: userInfo.Username,
		IsAdmin:  userInfo.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "book-app",
			Subject:   "access",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func parseJWT(tokenString string) (GetAuthInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return GetAuthInfo{}, errors.New("invalid token please login again")
		}
		return GetAuthInfo{}, err
	} else if claims, ok := token.Claims.(*jwtClaims); ok {
		return GetAuthInfo{ID: claims.UserID, Username: claims.Username, IsAdmin: claims.IsAdmin}, nil
	} else {
		return GetAuthInfo{}, errors.New("something went wrong while verifying user")
	}
}
