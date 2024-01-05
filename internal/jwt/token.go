package jwt

import (
	"fmt"

	gojwt "github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var signKey = []byte("secret")

type UserInfo struct {
	UserID   string
	UserName string
}

// ParseToken validates and parses the given JWT returning the claims
func ParseToken(tokenString string) (claims gojwt.MapClaims, err error) {

	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signKey, nil
	})
	if err != nil {
		return gojwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return gojwt.MapClaims{}, err
	}

	return claims, nil
}

// GetUserInfoFromToken Validates JWT and returns the userInfo struct
func GetUserInfoFromToken(tokenString string) (userInfo UserInfo, err error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		logrus.Errorf("error parsing token\nError: %s", err.Error())
		return userInfo, err
	}

	if claims["id"] == nil || claims["username"] == nil {
		return userInfo, fmt.Errorf("invalid token")
	}

	userInfo = UserInfo{
		UserID:   claims["id"].(string),
		UserName: claims["username"].(string),
	}

	return userInfo, nil
}

// GenerateToken generates a JWT token for the given user
func GenerateToken(userID, username string) (tokenString string, err error) {

	claims := gojwt.MapClaims{
		"id":       userID,
		"username": username,
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(signKey)
	if err != nil {
		logrus.Errorf("error generating token\nError: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
