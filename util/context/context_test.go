package context

import (
	"context"
	"testing"
)

const testKey key = 999

func TestWrap(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey, "123")
	myCtx := New(ctx)

	if myCtx.Context != ctx {
		t.Error("myCtx.Context should be the wrapped context")
	}
	if myCtx.Value(testKey) != "123" {
		t.Error("Values of wrapped context should be accessible")
	}
}

func TestNew(t *testing.T) {
	ctx := New(context.Background())
	ctx2 := New(ctx)

	if ctx2 != ctx {
		t.Error("Wrapped ak context should be the same context")
	}
}

func TestUser(t *testing.T) {
	expected := "user123"
	ctx := Background().WithUser(expected)

	user := ctx.User()
	if *user != expected {
		t.Errorf("User expected to be %s but was %s", expected, *user)
	}
}

func TestUserNull(t *testing.T) {
	ctx := New(context.Background())
	if ctx.User() != nil {
		t.Errorf("User expected to be null")
	}
}

func TestToken(t *testing.T) {
	expected := "my-token-insersion"
	ctx := Background().WithToken(expected)
	token := ctx.Token()
	if *token != expected {
		t.Errorf("token was not as expected")
	}
}

func TestTokenNull(t *testing.T) {
	ctx := Background()
	if ctx.Token() != nil {
		t.Errorf("Token expected to be null")
	}
}
