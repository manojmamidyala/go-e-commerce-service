package repository

import (
	"context"
	"mami/e-commerce/config"
	"mami/e-commerce/user/model"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
	db config.IDatabase
}

func NewUserRepository(db config.IDatabase) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(ctx, user)
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.Update(ctx, user)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := config.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, config.WithQuery(query)); err != nil {
		return nil, err
	}

	return &user, nil
}
