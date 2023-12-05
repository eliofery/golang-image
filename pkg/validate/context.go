package validate

import (
	"context"
	"github.com/go-playground/validator/v10"
)

type key string

var validationKey key = "validation"

func WithValidation(ctx context.Context, v *validator.Validate) context.Context {
	return context.WithValue(ctx, validationKey, v)
}

func Validation(ctx context.Context) *validator.Validate {
	v, ok := ctx.Value(validationKey).(*validator.Validate)
	if !ok {
		return nil
	}

	return v
}
