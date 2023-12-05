package user

import (
	"context"
)

type key string

const keyUser key = "user"

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, keyUser, user)
}

func CtxUser(ctx context.Context) *User {
	val := ctx.Value(keyUser)

	u, ok := val.(*User)
	if !ok {
		return nil
	}

	return u
}
