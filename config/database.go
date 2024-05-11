package config

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const DatabaseTimeout = 5 * time.Second

type Query struct {
	queryStr    string
	queryParams []any
}

type Database struct {
	db *gorm.DB
}

func NewDatabase(url string) (*Database, error) {
	database, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(4)
	sqlDB.SetMaxOpenConns(10)

	return &Database{
		db: database,
	}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func NewQuery(query string, args ...any) Query {
	return Query{
		queryStr:    query,
		queryParams: args,
	}
}

func (d *Database) AutoMigrate(models ...any) error {
	return d.db.AutoMigrate(models...)
}

func (d *Database) WithTransaction(function func() error) error {
	callback := func(db *gorm.DB) error {
		return function()
	}

	tx := d.db.Begin()
	if err := callback(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Preload(query string, args ...interface{}) IDatabase {
	d.db.Preload(query, args...)
	return d
}

func (d *Database) Create(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.Create(doc).Error
}

func (d *Database) CreateInBatches(ctx context.Context, docs any, batchSize int) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.CreateInBatches(docs, batchSize).Error
}

func (d *Database) Update(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.Save(doc).Error
}

func (d *Database) Delete(ctx context.Context, value any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	return query.Delete(value).Error
}

func (d *Database) FindById(ctx context.Context, id string, result any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	if err := d.db.Where("id = ? ", id).First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) FindOne(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Find(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.Find(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Count(ctx context.Context, model any, total *int64, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.Model(model).Count(total).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) applyOptions(opts ...FindOption) *gorm.DB {
	query := d.db

	opt := getOptions(opts...)

	if len(opt.preloads) != 0 {
		for _, preload := range opt.preloads {
			query = query.Preload(preload)
		}
	}

	if opt.query != nil {
		for _, q := range opt.query {
			query = query.Where(q.queryStr, q.queryParams)
		}
	}

	if opt.order != "" {
		query = query.Order(opt.order)
	}

	if opt.offset != 0 {
		query = query.Offset(opt.offset)
	}

	if opt.limit != 0 {
		query = query.Limit(opt.limit)
	}

	return query
}
