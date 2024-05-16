package service

import (
	"context"
	"errors"
	"mami/e-commerce/commons/logger"
	jwtToken "mami/e-commerce/pkg/jwtToken"
	"mami/e-commerce/pkg/utils"
	"mami/e-commerce/user/dto"
	"mami/e-commerce/user/model"
	"mami/e-commerce/user/repository"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Login(ctx context.Context, req *dto.LoginReq) (*model.User, string, string, error)
	Register(ctx context.Context, req *dto.RegisterReq) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	RefreshToken(ctx context.Context, userID string) (string, error)
	ChangePassword(ctx context.Context, id string, req *dto.ChangePasswordReq) error
}

type UserService struct {
	repo      repository.IUserRepository
	validator validator.Validate
}

func NewUserService(repo repository.IUserRepository, validator validator.Validate) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator,
	}
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginReq) (*model.User, string, string, error) {

	if err := s.validator.Struct(req); err != nil {
		return nil, "", "", err
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("Login.GetUserByEmail fail, email:  %s, error: %s", req.Email, err)
		return nil, "", "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", "", errors.New("wrong password")
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jwtToken.GenerateAccessToken(tokenData)
	refreshToken := jwtToken.GenerateRefeshToken(tokenData)

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterReq) (*model.User, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	var user model.User
	utils.Copy(&user, &req)
	err := s.repo.Create(ctx, &user)
	if err != nil {
		logger.Errorf("Register.Create fail, email %s, error: %s", req.Email, err)
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("GetUserByID.GetUserByID fail, id %s, error: %s", id, err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("RefreshToken.GetUserByID fail, id: %s, error: %s", userID, err)
		return "", err
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jwtToken.GenerateRefeshToken(tokenData)

	return accessToken, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id string, req *dto.ChangePasswordReq) error {
	if err := s.validator.Struct(req); err != nil {
		return err
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("ChangePassword.GetUserByID fail, id: %s, error: %s", id, err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.New("wrong password")
	}

	user.Password = utils.HashAndSalt([]byte(req.NewPassword))
	err = s.repo.Update(ctx, user)
	if err != nil {
		logger.Errorf("ChangePassword.Update fail, id: %s, error: %s", user.ID, err)
		return err
	}
	return nil
}
