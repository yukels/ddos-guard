package request

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/yukels/util/context"
)

// TokenFromRequest returns Authorization header value
func TokenFromRequest(ctx context.Context, r *http.Request) (string, error) {
	return TokenFromString(ctx, r.Header.Get("Authorization"))
}

// TokenFromString extract the OAuth2 token from string
func TokenFromString(ctx context.Context, t string) (string, error) {
	auth := strings.SplitN(t, " ", 2)
	if len(auth) != 2 || auth[0] != "Bearer" {
		return "", errors.New("Received invalid Token: No Authorization Header with Bearer token")
	}
	return auth[1], nil
}
