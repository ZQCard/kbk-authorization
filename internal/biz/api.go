package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/domain"
)

type ApiRepo interface {
	ListApiAll(ctx context.Context) ([]*domain.Api, error)
	ListApi(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Api, int64, error)
	CreateApi(ctx context.Context, Api *domain.Api) (*domain.Api, error)
	UpdateApi(ctx context.Context, Api *domain.Api) error
	DeleteApi(ctx context.Context, Api *domain.Api) error
}

type ApiUsecase struct {
	repo   ApiRepo
	logger *log.Helper
}

func NewApiUsecase(repo ApiRepo, logger log.Logger) *ApiUsecase {
	return &ApiUsecase{repo: repo, logger: log.NewHelper(log.With(logger, "module", "usecase/api"))}
}

func (c *ApiUsecase) ListApiAll(ctx context.Context) ([]*domain.Api, error) {
	return c.repo.ListApiAll(ctx)
}

func (c *ApiUsecase) ListApi(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Api, int64, error) {
	return c.repo.ListApi(ctx, page, pageSize, params)
}

func (c *ApiUsecase) CreateApi(ctx context.Context, api *domain.Api) (*domain.Api, error) {
	return c.repo.CreateApi(ctx, api)
}

func (c *ApiUsecase) UpdateApi(ctx context.Context, api *domain.Api) error {
	return c.repo.UpdateApi(ctx, api)
}

func (c *ApiUsecase) DeleteApi(ctx context.Context, api *domain.Api) error {
	return c.repo.DeleteApi(ctx, api)
}
