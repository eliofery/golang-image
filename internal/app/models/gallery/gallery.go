package gallery

import (
	"context"
)

type Gallery struct {
	ID     uint   `validate:"omitempty"`
	UserID uint   `validate:"required,min=1"`
	Title  string `validate:"required,max=255"`
}

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}
