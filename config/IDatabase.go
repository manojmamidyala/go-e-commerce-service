package config

import (
	"context"

	"gorm.io/gorm"
)

type IDatabase interface {
	GetDB() *gorm.DB
	AutoMigrate(models ...any) error
	WithTransaction(function func() error) error
	Create(ctx context.Context, doc any) error
	CreateInBatches(ctx context.Context, docs any, batchSize int) error
	Update(ctx context.Context, doc any) error
	Delete(ctx context.Context, value any, opts ...FindOption) error
	FindById(ctx context.Context, id string, result any) error
	FindOne(ctx context.Context, result any, opts ...FindOption) error
	Find(ctx context.Context, result any, opts ...FindOption) error
	Count(ctx context.Context, model any, total *int64, opts ...FindOption) error
}
