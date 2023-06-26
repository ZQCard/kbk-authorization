package service

import (
	"context"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
)

func (s *AuthorizationService) SetRolesForUser(ctx context.Context, req *v1.SetRolesForUserReq) (*v1.CheckResponse, error) {
	success, err := s.casbinUsecase.SetRolesForUser(ctx, req.Username, req.Roles)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{
		Success: success,
	}, nil
}

func (s *AuthorizationService) GetRolesForUser(ctx context.Context, req *v1.GetRolesForUserReq) (*v1.GetRolesForUserRes, error) {
	roles, err := s.casbinUsecase.GetRolesForUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &v1.GetRolesForUserRes{
		Roles: roles,
	}, nil
}

func (s *AuthorizationService) GetUsersForRole(ctx context.Context, req *v1.RoleNameReq) (*v1.GetUsersForRoleRes, error) {
	users, err := s.casbinUsecase.GetUsersForRole(ctx, req.Role)
	if err != nil {
		return nil, err
	}
	return &v1.GetUsersForRoleRes{
		Users: users,
	}, nil
}

func (s *AuthorizationService) DeleteRoleForUser(ctx context.Context, req *v1.DeleteRoleForUserReq) (*v1.CheckResponse, error) {
	success, err := s.casbinUsecase.DeleteRoleForUser(ctx, req.Username, req.Role)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{
		Success: success,
	}, nil
}

func (s *AuthorizationService) DeleteRolesForUser(ctx context.Context, req *v1.DeleteRolesForUserReq) (*v1.CheckResponse, error) {
	success, err := s.casbinUsecase.DeleteRolesForUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{
		Success: success,
	}, nil
}

func (s *AuthorizationService) GetPolicies(ctx context.Context, req *v1.RoleNameReq) (*v1.GetPoliciesRes, error) {

	rules, err := s.casbinUsecase.GetPolicies(ctx, req.Role)
	if err != nil {
		return nil, err
	}
	Res := []*v1.PolicyRules{}
	for _, v := range rules {
		Res = append(Res, &v1.PolicyRules{
			Path:   v.Path,
			Method: v.Method,
		})
	}
	return &v1.GetPoliciesRes{
		PolicyRules: Res,
	}, nil
}

func (s *AuthorizationService) UpdatePolicies(ctx context.Context, req *v1.UpdatePoliciesReq) (*v1.CheckResponse, error) {
	rules := []biz.PolicyRules{}
	for _, v := range req.PolicyRules {
		rules = append(rules, biz.PolicyRules{
			Path:   v.Path,
			Method: v.Method,
		})
	}
	success, err := s.casbinUsecase.UpdatePolicies(ctx, req.Role, rules)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{
		Success: success,
	}, nil
}
