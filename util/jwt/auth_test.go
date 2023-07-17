package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yukels/util/context"
)

func TestUsernameFromUnverifiedToken(t *testing.T) {
	tests := []struct {
		testCase string
		token    string
		expected string
	}{
		{"valid and not expired", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwiRGlzcGxheU5hbWUiOiJ2YWxpZCBhbmQgbm90IGV4cGlyZWQiLCJyb2xlIjpbIkFkbWluIiwidmQiXSwianRpIjoiMTIzNDU2IiwiZXhwIjo5NjY2NjE5NDY1fQ.VvRK63-wYFewML6VpKRitgqyuCoUpwoEB_I6xCqdgsY", "test"},
		{"valid and expired", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwiRGlzcGxheU5hbWUiOiJ2YWxpZCBhbmQgbm90IGV4cGlyZWQiLCJyb2xlIjpbIkFkbWluIiwidmQiXSwianRpIjoiMTIzNDU2IiwiZXhwIjoxNTY2NjE5NDY1fQ.tECKs1mT6FlezFDeLljGo_Jmt4yCLbDeZXoyK8Y-Y94", "test"},
		{"not valid", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwiRGlzcGxheU5hbWUiOiJ2YWxpZCBhbmQgbm90IGV4cGlyZWQiLCJyb2xlIjpbIkFkbWluIiwidmQiXSwianRpIjoiMTIzNDU2IiwiZXhwIjoxNTY2NjE5NDY1fQ.cMdEEvxjx0f4mVwHBkXZX2AsU3CeQ0UGLPTd-QoUKMg", "test"},
		{"not valid and not expired", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwiRGlzcGxheU5hbWUiOiJ2YWxpZCBhbmQgbm90IGV4cGlyZWQiLCJyb2xlIjpbIkFkbWluIiwidmQiXSwianRpIjoiMTIzNDU2IiwiZXhwIjo5NjY2NjE5NDY1fQ.3hQQUEkRB4QHKAHXHPw9RIHJqRSAucx2lzJzYyoOODs", "test"},
		{"not a token", "123456", ""},
	}

	ctx := context.Background()

	for idx, tst := range tests {
		actual := UsernameFromUnverifiedToken(ctx, tst.token)

		assert.Equalf(t, tst.expected, actual, "[%d] unexpected result on %s", idx, tst.testCase)
	}
}
