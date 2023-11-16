package service

import (
	"context"

	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthorizationService) GetMenuAll(ctx context.Context, req *emptypb.Empty) (*v1.GetMenuTreeRes, error) {

	menu, err := s.menuUsecase.GetMenuAll(ctx)
	if err != nil {
		return nil, err
	}
	list := []*v1.Menu{}
	for _, v := range menu {
		list = append(list, toPbMenu(v))
	}

	return &v1.GetMenuTreeRes{
		List: list,
	}, nil
}

func (s *AuthorizationService) GetMenuTree(ctx context.Context, req *emptypb.Empty) (*v1.GetMenuTreeRes, error) {

	menu, err := s.menuUsecase.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}
	list := []*v1.Menu{}
	for k, v := range menu {
		children := findChildrenMenu(menu[k])
		btns := []*v1.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, &v1.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   btn.CreatedAt,
				UpdatedAt:   btn.UpdatedAt,
			})
		}

		res := &v1.Menu{
			Id:        menu[k].Id,
			ParentId:  menu[k].ParentId,
			ParentIds: menu[k].ParentIds,
			Path:      menu[k].Path,
			Name:      menu[k].Name,
			Hidden:    menu[k].Hidden,
			Component: menu[k].Component,
			Sort:      menu[k].Sort,
			Title:     menu[k].Title,
			Icon:      menu[k].Icon,
			CreatedAt: menu[k].CreatedAt,
			UpdatedAt: menu[k].UpdatedAt,
			Children:  children,
			MenuBtns:  btns,
		}
		list = append(list, res)
	}

	return &v1.GetMenuTreeRes{
		List: list,
	}, nil

}

func findChildrenMenu(menu *domain.Menu) []*v1.Menu {
	children := []*v1.Menu{}
	if len(menu.Children) != 0 {
		for k := range menu.Children {
			btns := []*v1.MenuBtn{}
			for _, btn := range menu.Children[k].MenuBtns {
				btns = append(btns, &v1.MenuBtn{
					Id:          btn.Id,
					MenuId:      btn.MenuId,
					Name:        btn.Name,
					Description: btn.Description,
					Identifier:  btn.Identifier,
					CreatedAt:   btn.CreatedAt,
					UpdatedAt:   btn.UpdatedAt,
				})
			}

			children = append(children, &v1.Menu{
				Id:        menu.Children[k].Id,
				Name:      menu.Children[k].Name,
				Path:      menu.Children[k].Path,
				ParentId:  menu.Children[k].ParentId,
				ParentIds: menu.Children[k].ParentIds,
				Hidden:    menu.Children[k].Hidden,
				Component: menu.Children[k].Component,
				Sort:      menu.Children[k].Sort,
				Title:     menu.Children[k].Title,
				Icon:      menu.Children[k].Icon,
				CreatedAt: menu.Children[k].CreatedAt,
				UpdatedAt: menu.Children[k].UpdatedAt,
				MenuBtns:  btns,
				Children:  findChildrenMenu(&menu.Children[k]),
			})
		}
	}

	return children
}

func (s *AuthorizationService) CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (*v1.Menu, error) {
	btns := []domain.MenuBtn{}
	for _, v := range req.MenuBtns {
		btns = append(btns, domain.MenuBtn{
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
		})
	}
	bc := &domain.Menu{
		Name:      req.Name,
		Path:      req.Path,
		ParentId:  req.ParentId,
		ParentIds: req.ParentIds,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Title:     req.Title,
		Icon:      req.Icon,
		MenuBtns:  btns,
	}

	menu, err := s.menuUsecase.CreateMenu(ctx, bc)
	if err != nil {
		return nil, err
	}

	return toPbMenu(menu), nil
}

func (s *AuthorizationService) UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (*emptypb.Empty, error) {
	btns := []domain.MenuBtn{}
	for _, v := range req.MenuBtns {
		btns = append(btns, domain.MenuBtn{
			Id:          v.Id,
			MenuId:      v.MenuId,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
		})
	}
	bc := &domain.Menu{
		Id:        req.Id,
		Name:      req.Name,
		Path:      req.Path,
		ParentId:  req.ParentId,
		ParentIds: req.ParentIds,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Title:     req.Title,
		Icon:      req.Icon,
		MenuBtns:  btns,
	}
	return &emptypb.Empty{}, s.menuUsecase.UpdateMenu(ctx, bc)
}

func (s *AuthorizationService) DeleteMenu(ctx context.Context, req *v1.IdReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.menuUsecase.DeleteMenu(ctx, req.Id)
}

func toPBMenuBtn(btn domain.MenuBtn) *v1.MenuBtn {
	return &v1.MenuBtn{
		Id:          btn.Id,
		MenuId:      btn.MenuId,
		Name:        btn.Name,
		Description: btn.Description,
		Identifier:  btn.Identifier,
		CreatedAt:   btn.CreatedAt,
		UpdatedAt:   btn.UpdatedAt,
	}
}

func toPbMenu(menu *domain.Menu) *v1.Menu {
	var btns []*v1.MenuBtn
	if len(menu.MenuBtns) != 0 {
		for _, v := range menu.MenuBtns {
			btns = append(btns, toPBMenuBtn(v))
		}
	}

	return &v1.Menu{
		Id:        menu.Id,
		ParentId:  menu.ParentId,
		ParentIds: menu.ParentIds,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Title:     menu.Title,
		Icon:      menu.Icon,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
		MenuBtns:  btns,
	}
}

func (s *AuthorizationService) GetRoleMenuTree(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {

	menu, err := s.menuUsecase.GetRoleMenuTree(ctx, req.Role)
	if err != nil {
		return nil, err
	}
	list := []*v1.Menu{}
	for k := range menu {
		children := findChildrenMenu(menu[k])
		res := &v1.Menu{
			Id:        menu[k].Id,
			ParentId:  menu[k].ParentId,
			Path:      menu[k].Path,
			Name:      menu[k].Name,
			Hidden:    menu[k].Hidden,
			Component: menu[k].Component,
			Sort:      menu[k].Sort,
			Title:     menu[k].Title,
			Icon:      menu[k].Icon,
			CreatedAt: menu[k].CreatedAt,
			UpdatedAt: menu[k].UpdatedAt,
			Children:  children,
		}
		list = append(list, res)
	}

	return &v1.GetMenuTreeRes{
		List: list,
	}, nil

}

func (s *AuthorizationService) GetRoleMenu(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {
	menu, err := s.menuUsecase.GetRoleMenu(ctx, req.Role)
	if err != nil {
		return nil, err
	}
	list := []*v1.Menu{}
	for _, v := range menu {
		list = append(list, toPbMenu(v))
	}

	return &v1.GetMenuTreeRes{
		List: list,
	}, nil
}
func (s *AuthorizationService) SaveRoleMenu(ctx context.Context, req *v1.SaveRoleMenuReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.menuUsecase.SaveRoleMenu(ctx, req.RoleId, req.MenuIds)
}

func (s *AuthorizationService) GetRoleMenuBtn(ctx context.Context, req *v1.GetRoleMenuBtnReq) (*v1.GetRoleMenuBtnRes, error) {
	list, err := s.menuUsecase.GetRoleMenuBtn(ctx, req.RoleId, req.RoleName, req.MenuId)
	if err != nil {
		return nil, err
	}
	resList := []*v1.MenuBtn{}
	for _, v := range list {
		resList = append(resList, toPBMenuBtn(*v))
	}

	return &v1.GetRoleMenuBtnRes{
		List: resList,
	}, nil

}

func (s *AuthorizationService) SaveRoleMenuBtn(ctx context.Context, req *v1.SaveRoleMenuBtnReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.menuUsecase.SaveRoleMenuBtn(ctx, req.RoleId, req.MenuId, req.MenuBtnIds)
}
