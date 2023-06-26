package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/domain"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/pkg/utils/timeHelper"
)

type RoleEntity struct {
	BaseFields
	Domain string `gorm:"type:varchar(255);not null;comment:所在域"`
	Name   string `gorm:"type:varchar(255);not null;comment:名称"`
}

func (RoleEntity) TableName() string {
	return "roles"
}

type RoleRepo struct {
	data *Data
	log  *log.Helper
}

func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	repo := &RoleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/role")),
	}
	return repo
}

// searchRoleParam 搜索条件
func (r RoleRepo) searchRoleParam(ctx context.Context, params map[string]interface{}) *gorm.DB {
	conn := getDbWithDomain(ctx, r.data.db).Model(&RoleEntity{})
	if Id, ok := params["id"]; ok && Id.(int64) != 0 {
		conn = conn.Where("id = ?", Id)
	}
	if Id, ok := params["neq_id"]; ok && Id.(int64) != 0 {
		conn = conn.Where("id != ?", Id)
	}
	if name, ok := params["name"]; ok && name.(string) != "" {
		conn = conn.Where("name = ?", name)
	}
	// 开始时间
	if start, ok := params["created_at_start"]; ok && start.(string) != "" {
		conn = conn.Where("created_at >= ?", start.(string)+" 00:00:00")
	}
	// 结束时间
	if end, ok := params["created_at_end"]; ok && end.(string) != "" {
		conn = conn.Where("created_at <= ?", end.(string)+" 23:59:59")
	}

	return conn
}

func (r RoleRepo) GetRoleByParams(ctx context.Context, params map[string]interface{}) (record *RoleEntity, err error) {
	if len(params) == 0 {
		return &RoleEntity{}, v1.ErrorBadRequest("缺少搜索条件")
	}
	conn := r.searchRoleParam(ctx, params)
	if err = conn.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &RoleEntity{}, v1.ErrorRecordNotFound("数据不存在")
		}
		return record, v1.ErrorSystemError("GetRoleByParams First Error : %s", err.Error())
	}
	return record, nil
}

func (r RoleRepo) ListRoleAll(ctx context.Context) ([]*domain.Role, error) {
	list := []*RoleEntity{}
	conn := r.searchRoleParam(ctx, map[string]interface{}{})
	err := conn.Find(&list).Error
	if err != nil {
		return nil, v1.ErrorSystemError("ListRoleAll Find Error : %s", err.Error())
	}

	rv := make([]*domain.Role, 0, len(list))
	for _, record := range list {
		role := toDomainRole(record)
		rv = append(rv, role)
	}
	return rv, nil
}

func (r RoleRepo) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// 名称不得重复
	count := int64(0)
	r.searchRoleParam(ctx, map[string]interface{}{
		"name": role.Name,
	}).Count(&count)
	if count > 0 {
		return nil, v1.ErrorBadRequest("名称不得重复")
	}
	entity := &RoleEntity{}
	entity.Id = role.Id
	entity.Name = role.Name
	entity.Domain = role.Domain
	entity.Domain = getDomain(ctx)
	if err := getDbWithDomain(ctx, r.data.db).Create(entity).Error; err != nil {
		return nil, v1.ErrorSystemError("CreateRole Create Error : %s", err.Error())
	}
	response := toDomainRole(entity)
	return response, nil
}

func (r RoleRepo) UpdateRole(ctx context.Context, role *domain.Role) error {
	// 名称不得重复
	count := int64(0)
	r.searchRoleParam(ctx, map[string]interface{}{
		"name":   role.Name,
		"neq_id": role.Id,
	}).Count(&count)
	if count > 0 {
		return v1.ErrorBadRequest("名称不得重复")
	}
	// 根据Id查找记录
	record, err := r.GetRoleByParams(ctx, map[string]interface{}{
		"id": role.Id,
	})
	if err != nil {
		return err
	}
	// 更新字段
	record.Name = role.Name
	record.Domain = role.Domain
	record.Domain = getDomain(ctx)
	if err := r.data.db.Where("id = ?", record.Id).Save(record).Error; err != nil {
		return v1.ErrorSystemError("UpdateRole Save Error : %s", err.Error())
	}

	return nil
}

func (r RoleRepo) GetRole(ctx context.Context, params map[string]interface{}) (*domain.Role, error) {
	// 根据Id查找记录
	record, err := r.GetRoleByParams(ctx, params)
	if err != nil {
		return nil, err
	}
	// 返回数据
	response := toDomainRole(record)
	return response, nil
}

func (r RoleRepo) DeleteRole(ctx context.Context, role *domain.Role) error {
	// 根据Id查找记录
	record, err := r.GetRoleByParams(ctx, map[string]interface{}{
		"id": role.Id,
	})
	if err != nil {
		return err
	}
	if role.Id != record.Id {
		return v1.ErrorBadRequest("缺少搜索条件")
	}
	if err := getDbWithDomain(ctx, r.data.db).Unscoped().Where("Id = ?", role.Id).Delete(&RoleEntity{}).Error; err != nil {
		return v1.ErrorSystemError("DeleteRole Delete Error : %s", err.Error())
	}
	return nil
}

func toDomainRole(role *RoleEntity) *domain.Role {
	return &domain.Role{
		Id:        role.Id,
		Name:      role.Name,
		CreatedAt: timeHelper.FormatYMDHIS(&role.CreatedAt),
		UpdatedAt: timeHelper.FormatYMDHIS(&role.UpdatedAt),
	}
}
