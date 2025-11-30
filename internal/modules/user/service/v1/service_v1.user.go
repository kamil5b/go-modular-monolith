package v1

import (
	"context"
	"go-modular-monolith/internal/domain/user"
	"time"
)

type ServiceV1 struct {
	repo user.UserRepository
}

func NewServiceV1(r user.UserRepository) *ServiceV1 { return &ServiceV1{repo: r} }

func (s *ServiceV1) Create(ctx context.Context, req *user.CreateUserRequest, createdBy string) (*user.User, error) {
	ctx = s.repo.StartContext(ctx)
	var u user.User
	u.Name = req.Name
	u.Email = req.Email
	u.CreatedAt = time.Now().UTC()
	u.CreatedBy = createdBy
	if err := s.repo.Create(ctx, &u); err != nil {
		s.repo.DeferErrorContext(ctx, err)
		return nil, err
	}
	return &u, nil
}

func (s *ServiceV1) Get(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ServiceV1) List(ctx context.Context) ([]user.User, error) {
	return s.repo.List(ctx)
}

func (s *ServiceV1) Update(ctx context.Context, req *user.UpdateUserRequest, updatedBy string) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	now := time.Now().UTC()
	u.UpdatedAt = &now
	u.UpdatedBy = &updatedBy
	if err := s.repo.Update(ctx, u); err != nil {
		s.repo.DeferErrorContext(ctx, err)
		return nil, err
	}
	return u, nil
}

func (s *ServiceV1) Delete(ctx context.Context, id, by string) error {
	return s.repo.SoftDelete(ctx, id, by)
}
