package service

import (
	"context"
)

type IHealthService interface {
	Check(ctx context.Context) (string, error)
}

type HealthService struct {
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (h *HealthService) Check(ctx context.Context) (string, error) {
	return "UP", nil
}
