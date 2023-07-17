package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"github.com/yukels/util/context"
)

const (
	userFromClaim = "sub"
)

func UsernameFromUnverifiedToken(ctx context.Context, tokenString string) string {
	jwtParser := &jwt.Parser{}

	if unverifiedToken, _, err := jwtParser.ParseUnverified(tokenString, jwt.MapClaims{}); err == nil {
		username, err := GetUsername(ctx, unverifiedToken)

		if err == nil {
			return username
		}
	}
	return ""
}

func getClaims(ctx context.Context, token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("No valid claims")
}

func GetUsername(ctx context.Context, token *jwt.Token) (string, error) {
	claims, err := getClaims(ctx, token)

	if err != nil {
		return "", err
	}
	if username, ok := claims[userFromClaim].(string); ok {
		return username, nil
	}
	return "", errors.New("No valid username claim")
}
