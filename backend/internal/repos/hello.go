package repos

import (
	"context"
	"gorm.io/gorm"
)

type Hello struct {
	db gorm.DB
}

func NewHello(db *gorm.DB) *Hello {
	return &Hello{*db}
}

func (h *Hello) Hello(ctx context.Context) (string, error) {
	return "hello, world", nil
}
