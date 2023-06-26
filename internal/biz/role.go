package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/ZQCard/kbk-authorization/internal/domain"
)

type RoleRepo interface {
	ListRoleAll(ctx context.Context) ([]*domain.Role, error)
	CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	GetRole(ctx context.Context, params map[string]interface{}) (*domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, role *domain.Role) error
}

type RoleUsecase struct {
	repo   RoleRepo
	logger *log.Helper
}

func NewRoleUsecase(repo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{repo: repo, logger: log.NewHelper(log.With(logger, "module", "usecase/role"))}
}

func (c *RoleUsecase) ListRoleAll(ctx context.Context) ([]*domain.Role, error) {
	return c.repo.ListRoleAll(ctx)
}

func (c *RoleUsecase) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	return c.repo.CreateRole(ctx, role)
}

func (c *RoleUsecase) UpdateRole(ctx context.Context, role *domain.Role) error {
	return c.repo.UpdateRole(ctx, role)
}

func (c *RoleUsecase) DeleteRole(ctx context.Context, role *domain.Role) error {
	return c.repo.DeleteRole(ctx, role)
}
