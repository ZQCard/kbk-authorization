package biz

import (
	"context"

	"github.com/ZQCard/kbk-authorization/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
)

type MenuRepo interface {
	GetMenuTree(ctx context.Context) ([]*domain.Menu, error)
	GetMenuAll(ctx context.Context) ([]*domain.Menu, error)
	CreateMenu(ctx context.Context, menu *domain.Menu) (*domain.Menu, error)
	UpdateMenu(ctx context.Context, menu *domain.Menu) error
	DeleteMenu(ctx context.Context, id int64) error
	SaveRoleMenu(ctx context.Context, roleId int64, menuIds []int64) error
	GetRoleMenuBtn(ctx context.Context, roleId int64, roleName string, menuId int64) ([]*domain.MenuBtn, error)
	SaveRoleMenuBtn(ctx context.Context, roleId int64, menuId int64, btnIds []int64) error
	GetRoleMenu(ctx context.Context, role string) ([]*domain.Menu, error)
	GetRoleMenuTree(ctx context.Context, role string) ([]*domain.Menu, error)
}

type MenuUsecase struct {
	repo   MenuRepo
	logger *log.Helper
}

func NewMenuUsecase(repo MenuRepo, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{repo: repo, logger: log.NewHelper(log.With(logger, "module", "usecase/menu"))}
}

func (c *MenuUsecase) CreateMenu(ctx context.Context, data *domain.Menu) (*domain.Menu, error) {
	result, err := c.repo.CreateMenu(ctx, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MenuUsecase) DeleteMenu(ctx context.Context, id int64) error {
	return c.repo.DeleteMenu(ctx, id)
}

func (c *MenuUsecase) UpdateMenu(ctx context.Context, data *domain.Menu) error {
	return c.repo.UpdateMenu(ctx, data)
}

func (c *MenuUsecase) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {
	result, err := c.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MenuUsecase) GetMenuAll(ctx context.Context) ([]*domain.Menu, error) {
	result, err := c.repo.GetMenuAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MenuUsecase) SaveRoleMenu(ctx context.Context, roleId int64, menuIds []int64) error {
	return c.repo.SaveRoleMenu(ctx, roleId, menuIds)
}

func (c *MenuUsecase) GetRoleMenu(ctx context.Context, role string) ([]*domain.Menu, error) {
	result, err := c.repo.GetRoleMenu(ctx, role)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MenuUsecase) GetRoleMenuTree(ctx context.Context, role string) ([]*domain.Menu, error) {
	result, err := c.repo.GetRoleMenuTree(ctx, role)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MenuUsecase) GetRoleMenuBtn(ctx context.Context, roleId int64, roleName string, menuId int64) ([]*domain.MenuBtn, error) {
	return c.repo.GetRoleMenuBtn(ctx, roleId, roleName, menuId)
}

func (c *MenuUsecase) SaveRoleMenuBtn(ctx context.Context, roleId int64, menuId int64, btnIds []int64) error {
	return c.repo.SaveRoleMenuBtn(ctx, roleId, menuId, btnIds)
}
