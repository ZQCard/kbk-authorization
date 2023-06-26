package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/domain"
)

func (s *AuthorizationService) GetRoleAll(ctx context.Context, req *emptypb.Empty) (*v1.GetRoleAllRes, error) {
	list, err := s.roleUsecase.ListRoleAll(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.GetRoleAllRes{}
	for _, v := range list {
		res.List = append(res.List, toPbRole(v))
	}
	return res, nil
}

func (s *AuthorizationService) CreateRole(ctx context.Context, req *v1.CreateRoleReq) (*v1.Role, error) {
	res, err := s.roleUsecase.CreateRole(ctx, &domain.Role{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return toPbRole(res), nil
}

func (s *AuthorizationService) UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (*v1.CheckResponse, error) {
	err := s.roleUsecase.UpdateRole(ctx, &domain.Role{
		Id:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{Success: true}, nil
}

func (s *AuthorizationService) DeleteRole(ctx context.Context, req *v1.DeleteRoleReq) (*v1.CheckResponse, error) {
	err := s.roleUsecase.DeleteRole(ctx, &domain.Role{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{Success: true}, nil
}

func toPbRole(role *domain.Role) *v1.Role {
	return &v1.Role{
		Id:        role.Id,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
