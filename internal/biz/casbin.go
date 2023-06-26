package biz

import (
	"context"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type CasbinRepo interface {
	SetRolesForUser(ctx context.Context, username string, roles []string) (bool, error)
	GetRolesForUser(ctx context.Context, username string) ([]string, error)
	GetUsersForRole(ctx context.Context, user string) ([]string, error)
	DeleteRoleForUser(ctx context.Context, username string, role string) (bool, error)
	DeleteRolesForUser(ctx context.Context, username string) (bool, error)
	GetPolicies(ctx context.Context, role string) ([]*PolicyRules, error)
	UpdatePolicies(ctx context.Context, username string, rules []PolicyRules) (bool, error)
	CheckAuthorization(ctx context.Context, sub, obj, act string) (success bool, err error)
}

type CasbinUsecase struct {
	repo   CasbinRepo
	logger *log.Helper
}

func NewCasbinUsecase(repo CasbinRepo, logger log.Logger) *CasbinUsecase {
	return &CasbinUsecase{repo: repo, logger: log.NewHelper(log.With(logger, "module", "usecase/casbin"))}
}

func (c *CasbinUsecase) SetRolesForUser(ctx context.Context, username string, roles []string) (bool, error) {
	return c.repo.SetRolesForUser(ctx, username, roles)
}

func (c *CasbinUsecase) GetRolesForUser(ctx context.Context, username string) ([]string, error) {
	return c.repo.GetRolesForUser(ctx, username)
}

func (c *CasbinUsecase) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	if role == "" {
		return []string{}, kerrors.BadRequest("BAD REQUEST", "角色不得为空")
	}
	return c.repo.GetUsersForRole(ctx, role)
}

func (c *CasbinUsecase) DeleteRoleForUser(ctx context.Context, username string, role string) (bool, error) {
	if role == "" || username == "" {
		return false, kerrors.BadRequest("BAD REQUEST", "用户与角色不得为空")
	}
	return c.repo.DeleteRoleForUser(ctx, username, role)
}

func (c *CasbinUsecase) DeleteRolesForUser(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, kerrors.BadRequest("BAD REQUEST", "用户不得为空")
	}
	return c.repo.DeleteRolesForUser(ctx, username)
}

type PolicyRules struct {
	Path   string
	Method string
}

func (c *CasbinUsecase) GetPolicies(ctx context.Context, role string) ([]*PolicyRules, error) {
	return c.repo.GetPolicies(ctx, role)
}

func (c *CasbinUsecase) UpdatePolicies(ctx context.Context, role string, rules []PolicyRules) (bool, error) {
	if role == "" {
		return false, kerrors.BadRequest("BAD REQUEST", "角色不得为空")
	}
	return c.repo.UpdatePolicies(ctx, role, rules)
}

func (c *CasbinUsecase) CheckAuthorization(ctx context.Context, sub, obj, act string) (bool, error) {
	return c.repo.CheckAuthorization(ctx, sub, obj, act)
}
