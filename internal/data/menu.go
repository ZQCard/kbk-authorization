package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/domain"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/pkg/utils/timeHelper"
)

type MenuEntity struct {
	BaseFields
	Domain    string           `gorm:"type:varchar(255);not null;comment:所在域"`
	ParentId  int64            `gorm:"type:int;comment:父级id"`
	ParentIds string           `gorm:"type:int;comment:父级id字符串 英文逗号分割"`
	Name      string           `gorm:"type:varchar(255);not null;comment:菜单名"`
	Path      string           `gorm:"type:varchar(255);not null;comment:前端路径"`
	Hidden    bool             `gorm:"not null;comment:是否隐藏 0否1是"`
	Component string           `gorm:"type:varchar(255);not null;comment:前端文件路径"`
	Sort      int64            `gorm:"type:int;comment:排序"`
	Title     string           `gorm:"type:varchar(255);not null;comment:页面名称"`
	Icon      string           `gorm:"type:varchar(255);not null;comment:菜单图标"`
	MenuBtns  []*MenuBtnEntity `gorm:"foreignKey:menu_id;"`
	Children  []MenuEntity     `gorm:"-"`
}

func (MenuEntity) TableName() string {
	return "menus"
}

type MenuBtnEntity struct {
	Id          int64     `gorm:"primarykey;type:int;comment:主键id"`
	Domain      string    `gorm:"type:varchar(255);not null;comment:所在域"`
	MenuId      int64     `gorm:"type:int;comment:菜单id"`
	Name        string    `gorm:"type:varchar(255);not null;comment:按钮名称"`
	Description string    `gorm:"type:varchar(255);not null;comment:描述"`
	Identifier  string    `gorm:"type:varchar(255);not null;comment:英文标识"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;comment:更新时间"`
}

func (MenuBtnEntity) TableName() string {
	return "menu_btns"
}

type RoleMenuEntity struct {
	Id     int64  `gorm:"primarykey;type:int;comment:主键id"`
	Domain string `gorm:"type:varchar(255);not null;comment:所在域"`
	RoleId int64  `gorm:"type:int;comment:角色id"`
	MenuId int64  `gorm:"type:int;comment:菜单id"`
}

type RoleMenuBtnEntity struct {
	Id     int64  `gorm:"primarykey;type:int;comment:主键id"`
	Domain string `gorm:"type:varchar(255);not null;comment:所在域"`
	RoleId int64  `gorm:"type:int;comment:角色id"`
	MenuId int64  `gorm:"type:int;comment:菜单id"`
	BtnId  int64  `gorm:"type:int;comment:按钮id"`
}

func (RoleMenuBtnEntity) TableName() string {
	return "role_menu_btn"
}

func (RoleMenuEntity) TableName() string {
	return "role_menu"
}

type MenuRepo struct {
	data *Data
	log  *log.Helper
}

func NewMenuRepo(data *Data, logger log.Logger) biz.MenuRepo {
	repo := &MenuRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/menu")),
	}
	return repo
}

const childModuleMenu = "Menu"
const childModuleRoleMenu = "RoleMenu"

func (r MenuRepo) GetMenuAll(ctx context.Context) ([]*domain.Menu, error) {
	var response []*domain.Menu
	var menus []MenuEntity
	// 获取所有根菜单
	err := getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Preload("MenuBtns").Order("sort ASC").Find(&menus).Error
	if err != nil {
		return response, v1.ErrorSystemError("GetMenuAll Find Error : %s", err.Error())
	}
	for _, v := range menus {
		btns := []domain.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, domain.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   timeHelper.FormatYMDHIS(&btn.CreatedAt),
				UpdatedAt:   timeHelper.FormatYMDHIS(&btn.UpdatedAt),
			})
		}

		response = append(response, &domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
			MenuBtns:  btns,
		})
	}
	return response, nil
}

func (r MenuRepo) GetMenuTree(ctx context.Context) ([]*domain.Menu, error) {

	var response []*domain.Menu

	var menus []MenuEntity
	// 获取所有根菜单
	err := getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("parent_id = 0").Preload("MenuBtns").Order("sort ASC").Find(&menus).Error
	if err != nil {
		return response, v1.ErrorSystemError("GetMenuTree Find Error : %s", err.Error())
	}
	for _, v := range menus {
		btns := []domain.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, domain.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   timeHelper.FormatYMDHIS(&btn.CreatedAt),
				UpdatedAt:   timeHelper.FormatYMDHIS(&btn.UpdatedAt),
			})
		}
		response = append(response, &domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Path:      v.Path,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
			MenuBtns:  btns,
		})
	}
	for k := range response {
		err := r.findChildrenMenu(ctx, response[k])
		if err != nil {
			return response, v1.ErrorSystemError("GetMenuTree Find Error : %s", err.Error())
		}
	}
	return response, nil
}

func (r MenuRepo) findChildrenMenu(ctx context.Context, menu *domain.Menu) (err error) {
	var tmp []MenuEntity
	err = getDbWithDomain(ctx, r.data.db).Where("parent_id = ? ", menu.Id).Preload("MenuBtns").Order("sort ASC").Find(&tmp).Error
	menu.Children = []domain.Menu{}
	for _, v := range tmp {
		btns := []domain.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, domain.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   timeHelper.FormatYMDHIS(&btn.CreatedAt),
				UpdatedAt:   timeHelper.FormatYMDHIS(&btn.UpdatedAt),
			})
		}

		menu.Children = append(menu.Children, domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			Path:      v.Path,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
			MenuBtns:  btns,
		})
	}
	if len(menu.Children) > 0 {
		for k := range menu.Children {
			err = r.findChildrenMenu(ctx, &menu.Children[k])
		}
	}
	return err
}

func (r MenuRepo) CreateMenu(ctx context.Context, reqData *domain.Menu) (*domain.Menu, error) {
	domainStr := getDomain(ctx)
	btns := []*MenuBtnEntity{}
	for _, v := range reqData.MenuBtns {
		btns = append(btns, &MenuBtnEntity{
			Domain:      domainStr,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
		})
	}

	var menu MenuEntity
	menu = MenuEntity{
		Domain:    domainStr,
		Name:      reqData.Name,
		ParentId:  reqData.ParentId,
		ParentIds: reqData.ParentIds,
		Path:      reqData.Path,
		Hidden:    reqData.Hidden,
		Component: reqData.Component,
		Sort:      reqData.Sort,
		Title:     reqData.Title,
		Icon:      reqData.Icon,
		MenuBtns:  btns,
	}
	err := getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Create(&menu).Error
	if err != nil {
		return nil, v1.ErrorSystemError("CreateMenu Create Error : %s", err.Error())
	}
	btns2 := []domain.MenuBtn{}
	for _, v := range menu.MenuBtns {
		btns2 = append(btns2, domain.MenuBtn{
			Id:          v.Id,
			MenuId:      v.MenuId,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
		})
	}
	res := &domain.Menu{
		Id:        menu.Id,
		Name:      menu.Name,
		Path:      menu.Path,
		ParentId:  menu.ParentId,
		ParentIds: menu.ParentIds,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Title:     menu.Title,
		Icon:      menu.Icon,
		CreatedAt: timeHelper.FormatYMDHIS(&menu.CreatedAt),
		UpdatedAt: timeHelper.FormatYMDHIS(&menu.UpdatedAt),
		MenuBtns:  btns2,
	}
	return res, nil
}

func (r MenuRepo) UpdateMenu(ctx context.Context, reqData *domain.Menu) error {
	domainStr := getDomain(ctx)
	btns := []*MenuBtnEntity{}
	for _, v := range reqData.MenuBtns {
		btns = append(btns, &MenuBtnEntity{
			Id:          v.Id,
			Domain:      domainStr,
			MenuId:      reqData.Id,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
		})
	}
	var menu MenuEntity
	getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("id = ?", reqData.Id).First(&menu)
	menu.Id = reqData.Id
	menu.Domain = domainStr
	menu.Name = reqData.Name
	menu.ParentId = reqData.ParentId
	menu.ParentIds = reqData.ParentIds
	menu.Path = reqData.Path
	menu.Hidden = reqData.Hidden
	menu.Component = reqData.Component
	menu.Sort = reqData.Sort
	menu.Title = reqData.Title
	menu.Icon = reqData.Icon
	menu.MenuBtns = btns
	// 关联数据更新
	tx := getDbWithDomain(ctx, r.data.db).Begin()
	err := tx.Model(MenuEntity{}).Where("id = ?", menu.Id).Session(&gorm.Session{FullSaveAssociations: true}).Save(&menu).Error
	if err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("UpdateMenu Save Error : %s", err.Error())
	}
	// 先删除,后添加
	tx.Where("menu_id  = ?", menu.Id).Unscoped().Delete(&MenuBtnEntity{})
	// 保存按钮
	for _, v := range menu.MenuBtns {
		if err = tx.Model(MenuBtnEntity{}).Where("id = ?", v.Id).Create(&v).Error; err != nil {
			tx.Rollback()
			return v1.ErrorSystemError("UpdateMenu MenuBtnEntityCreate Error : %s", err.Error())
		}
	}
	tx.Commit()
	return nil
}

func (r MenuRepo) DeleteMenu(ctx context.Context, id int64) error {
	var menu MenuEntity
	// 查看菜单是否存在
	err := getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("id = ?", id).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrorRecordNotFound("数据不存在")
	}
	if err != nil {
		return v1.ErrorSystemError("DeleteMenu First Error : %s", err.Error())
	}

	// 如果有角色使用菜单无法删除
	roleMenu := RoleMenuEntity{}
	getDbWithDomain(ctx, r.data.db).Model(RoleMenuEntity{}).Where("menu_id = ?", id).First(&roleMenu)
	if roleMenu.Id != 0 {
		return v1.ErrorBadRequest("菜单已被使用,无法删除")
	}

	tx := getDbWithDomain(ctx, r.data.db).Begin()
	// 删除菜单与按钮的关联关系
	err = tx.Unscoped().Model(MenuBtnEntity{}).Where("menu_id = ?", id).Delete(&RoleMenuEntity{}).Error
	if err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("DeleteMenu RoleMenuEntityDelete Error : %s", err.Error())
	}
	// 删除菜单
	err = tx.Unscoped().Model(MenuEntity{}).Where("id = ?", id).Delete(&menu).Error
	if err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("DeleteMenu RoleMenuEntityDelete Error : %s", err.Error())
	}

	tx.Commit()
	return nil
}

func (r MenuRepo) GetRoleMenuTree(ctx context.Context, role string) ([]*domain.Menu, error) {
	// 查询角色拥有哪些菜单
	var response []*domain.Menu

	var menus []MenuEntity
	tmpRole := RoleEntity{}
	err := getDbWithDomain(ctx, r.data.db).Model(RoleEntity{}).Where("name = ?", role).First(&tmpRole).Error
	if tmpRole.Id == 0 {
		return response, nil
	}
	// 查看角色拥有菜单id
	menuIds := r.getMenuIdsByRoleId(ctx, tmpRole.Id)
	if len(menuIds) == 0 {
		return response, nil
	}

	// 获取所有根菜单
	err = getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("parent_id = 0 AND id IN (?)", menuIds).Preload("MenuBtns").Order("sort ASC").Find(&menus).Error
	if err != nil {
		return response, v1.ErrorSystemError("GetRoleMenuTree Find Error : %s", err.Error())
	}
	for _, v := range menus {
		response = append(response, &domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			ParentId:  v.ParentId,
			Path:      v.Path,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
		})
	}
	for k := range response {
		err := r.findChildrenRoleMenu(ctx, response[k], menuIds)
		if err != nil {
			return response, v1.ErrorSystemError("GetRoleMenuTree Find Error : %s", err.Error())
		}
	}
	return response, nil
}

func (r MenuRepo) findChildrenRoleMenu(ctx context.Context, menu *domain.Menu, menuIds []int64) (err error) {
	var tmp []MenuEntity
	err = getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("parent_id = ? AND id IN (?) ", menu.Id, menuIds).Preload("MenuBtns").Find(&tmp).Error
	menu.Children = []domain.Menu{}
	for _, v := range tmp {
		btns := []domain.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, domain.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				CreatedAt:   timeHelper.FormatYMDHIS(&btn.CreatedAt),
				UpdatedAt:   timeHelper.FormatYMDHIS(&btn.UpdatedAt),
			})
		}
		menu.Children = append(menu.Children, domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			Path:      v.Path,
			ParentId:  v.ParentId,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
			MenuBtns:  btns,
		})
	}
	if len(menu.Children) > 0 {
		for k := range menu.Children {
			err = r.findChildrenRoleMenu(ctx, &menu.Children[k], menuIds)
		}
	}
	return err
}

func (r MenuRepo) GetRoleMenu(ctx context.Context, role string) ([]*domain.Menu, error) {
	// 查询角色拥有哪些菜单
	var response []*domain.Menu

	var menus []MenuEntity
	tmpRole := RoleEntity{}
	err := getDbWithDomain(ctx, r.data.db).Model(RoleEntity{}).Where("name = ?", role).First(&tmpRole).Error
	if tmpRole.Id == 0 {
		return response, nil
	}

	// 查看角色拥有菜单id
	menuIds := r.getMenuIdsByRoleId(ctx, tmpRole.Id)
	if len(menuIds) == 0 {
		return response, nil
	}

	// 获取所有根菜单
	err = getDbWithDomain(ctx, r.data.db).Model(MenuEntity{}).Where("id IN (?)", menuIds).Preload("MenuBtns").Find(&menus).Error
	if err != nil {
		return response, v1.ErrorSystemError("GetRoleMenu Find Error : %s", err.Error())
	}
	for _, v := range menus {
		btns := []domain.MenuBtn{}
		for _, btn := range v.MenuBtns {
			btns = append(btns, domain.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				CreatedAt:   timeHelper.FormatYMDHIS(&btn.CreatedAt),
				UpdatedAt:   timeHelper.FormatYMDHIS(&btn.UpdatedAt),
			})
		}
		response = append(response, &domain.Menu{
			Id:        v.Id,
			Name:      v.Name,
			ParentId:  v.ParentId,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			CreatedAt: timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt: timeHelper.FormatYMDHIS(&v.UpdatedAt),
			MenuBtns:  btns,
		})
	}
	return response, nil
}

func (r MenuRepo) SaveRoleMenu(ctx context.Context, roleId int64, menuIds []int64) error {
	domainStr := getDomain(ctx)
	tx := getDbWithDomain(ctx, r.data.db).Begin()
	// 先删除数据
	err := tx.Where("role_id = ?", roleId).Delete(&RoleMenuEntity{}).Error
	if err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("SaveRoleMenu Delete Error : %s", err.Error())
	}
	if len(menuIds) == 0 {
		tx.Commit()
		return nil
	}
	// 批量插入数据
	roleMenu := []RoleMenuEntity{}
	for _, v := range menuIds {
		roleMenu = append(roleMenu, RoleMenuEntity{
			Domain: domainStr,
			RoleId: roleId,
			MenuId: v,
		})
	}
	if err := tx.Create(&roleMenu).Error; err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("SaveRoleMenu Create Error : %s", err.Error())
	}
	tx.Commit()
	return nil
}

func (r MenuRepo) getMenuIdsByRoleId(ctx context.Context, roleId int64) (menuIds []int64) {
	// 查询角色拥有哪些菜单
	var roleMenu []RoleMenuEntity

	getDbWithDomain(ctx, r.data.db).Model(RoleMenuEntity{}).Where("role_id = ?", roleId).Find(&roleMenu)
	// 查看角色拥有菜单id
	for _, v := range roleMenu {
		menuIds = append(menuIds, v.MenuId)
	}
	return menuIds
}

func (r MenuRepo) GetRoleMenuBtn(ctx context.Context, roleId int64, roleName string, menuId int64) (response []*domain.MenuBtn, err error) {
	// 如果角色名称不为空， 则根据名称查找角色id
	if roleName != "" {
		tmpRole := RoleEntity{}
		err := getDbWithDomain(ctx, r.data.db).Model(RoleEntity{}).Where("name = ?", roleName).First(&tmpRole).Error
		if tmpRole.Id == 0 {
			return response, nil
		}
		if err != nil {
			return nil, err
		}

		if roleId != 0 && tmpRole.Id != roleId {
			return nil, v1.ErrorBadRequest("角色参数错误")
		}
		roleId = tmpRole.Id
	}

	// 查询角色拥有哪些菜单按钮
	var roleMenuBtn []RoleMenuBtnEntity
	conn := getDbWithDomain(ctx, r.data.db).Model(RoleMenuBtnEntity{})
	if menuId != 0 {
		conn = conn.Where("menu_id = ?", menuId)
	}
	if roleId != 0 {
		conn = conn.Where("role_id = ?", roleId)
	}
	err = conn.Find(&roleMenuBtn).Error
	var btnIds []int64
	// 查看角色拥有菜单id
	for _, v := range roleMenuBtn {
		btnIds = append(btnIds, v.BtnId)
	}

	btnResponse := []*MenuBtnEntity{}
	if len(btnIds) == 0 {
		return nil, nil
	}
	getDbWithDomain(ctx, r.data.db).Model(&MenuBtnEntity{}).Where("id IN (?)", btnIds).Find(&btnResponse)
	for _, v := range btnResponse {
		response = append(response, &domain.MenuBtn{
			Id:          v.Id,
			MenuId:      v.MenuId,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
			CreatedAt:   timeHelper.FormatYMDHIS(&v.CreatedAt),
			UpdatedAt:   timeHelper.FormatYMDHIS(&v.UpdatedAt),
		})
	}
	return response, err
}

func (r MenuRepo) SaveRoleMenuBtn(ctx context.Context, roleId int64, menuId int64, btnIds []int64) error {
	domain := getDomain(ctx)
	tx := getDbWithDomain(ctx, r.data.db).Begin()
	// 先删除数据
	err := tx.Where("role_id = ? AND menu_id = ?", roleId, menuId).Delete(&RoleMenuBtnEntity{}).Error
	if err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("SaveRoleMenuBtn Delete Error : %s", err.Error())
	}
	if len(btnIds) == 0 {
		tx.Commit()
		return nil
	}
	// 批量插入数据
	roleMenuBtn := []RoleMenuBtnEntity{}
	for _, v := range btnIds {
		roleMenuBtn = append(roleMenuBtn, RoleMenuBtnEntity{
			Domain: domain,
			RoleId: roleId,
			MenuId: menuId,
			BtnId:  v,
		})
	}
	if err := tx.Create(&roleMenuBtn).Error; err != nil {
		tx.Rollback()
		return v1.ErrorSystemError("SaveRoleMenuBtn Create Error : %s", err.Error())
	}
	tx.Commit()
	return nil
}

func (r MenuRepo) getRoleMenuRedisKeyPre(domain string) string {
	return r.data.rdbPrefix + ":" + domain + ":" + childModuleRoleMenu + ":"
}

func (r MenuRepo) getMenuRedisKeyPre(domain string) string {
	return r.data.rdbPrefix + ":" + domain + ":" + ":" + childModuleMenu + ":"
}
