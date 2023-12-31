package service

import (
	"context"

	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthorizationService struct {
	v1.UnimplementedAuthorizationServiceServer
	casbinUsecase *biz.CasbinUsecase
	roleUsecase   *biz.RoleUsecase
	apiUsecase    *biz.ApiUsecase
	menuUsecase   *biz.MenuUsecase
	log           *log.Helper
}

func NewAuthorizationService(
	menuUsecase *biz.MenuUsecase,
	casbinUsecase *biz.CasbinUsecase,
	roleUsecase *biz.RoleUsecase,
	apiUsecase *biz.ApiUsecase,
	logger log.Logger,
) *AuthorizationService {
	return &AuthorizationService{
		log:           log.NewHelper(log.With(logger, "module", "service/authorization")),
		roleUsecase:   roleUsecase,
		apiUsecase:    apiUsecase,
		casbinUsecase: casbinUsecase,
		menuUsecase:   menuUsecase,
	}
}

func (s *AuthorizationService) CheckAuthorization(ctx context.Context, req *v1.CheckAuthorizationReq) (*emptypb.Empty, error) {
	success, err := s.casbinUsecase.CheckAuthorization(ctx, req.Sub, req.Obj, req.Act)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, errors.Forbidden("Unauthorized", "暂无权限")
	}
	return &emptypb.Empty{}, nil
}
