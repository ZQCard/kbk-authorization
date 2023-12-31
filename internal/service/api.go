package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/domain"
)

func (s *AuthorizationService) GetApiListAll(ctx context.Context, req *emptypb.Empty) (*v1.GetApiListAllRes, error) {
	list, err := s.apiUsecase.ListApiAll(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.GetApiListAllRes{}
	for _, v := range list {
		res.List = append(res.List, toPbApi(v))
	}
	return res, nil
}

func (s *AuthorizationService) GetApiList(ctx context.Context, req *v1.GetApiListReq) (*v1.GetApiListPageRes, error) {
	params := make(map[string]interface{})
	params["name"] = req.Name
	params["method"] = req.Method
	params["path"] = req.Path
	params["group"] = req.Group

	list, count, err := s.apiUsecase.ListApi(ctx, req.Page, req.PageSize, params)
	if err != nil {
		return nil, err
	}
	res := &v1.GetApiListPageRes{}
	res.Total = int64(count)
	for _, v := range list {
		res.List = append(res.List, toPbApi(v))
	}
	return res, nil
}

func (s *AuthorizationService) CreateApi(ctx context.Context, req *v1.CreateApiReq) (*v1.Api, error) {
	res, err := s.apiUsecase.CreateApi(ctx, &domain.Api{
		Name:   req.Name,
		Group:  req.Group,
		Method: req.Method,
		Path:   req.Path,
	})
	if err != nil {
		return nil, err
	}
	return toPbApi(res), nil
}

func (s *AuthorizationService) UpdateApi(ctx context.Context, req *v1.UpdateApiReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.apiUsecase.UpdateApi(ctx, &domain.Api{
		Id:     req.Id,
		Name:   req.Name,
		Group:  req.Group,
		Method: req.Method,
		Path:   req.Path,
	})
}

func (s *AuthorizationService) DeleteApi(ctx context.Context, req *v1.DeleteApiReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.apiUsecase.DeleteApi(ctx, &domain.Api{
		Id: req.Id,
	})
}

func toPbApi(api *domain.Api) *v1.Api {
	return &v1.Api{
		Id:        api.Id,
		Domain:    api.Domain,
		Name:      api.Name,
		Group:     api.Group,
		Method:    api.Method,
		Path:      api.Path,
		CreatedAt: api.CreatedAt,
		UpdatedAt: api.UpdatedAt,
	}
}
