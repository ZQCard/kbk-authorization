package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/biz"
	"github.com/ZQCard/kbk-authorization/internal/domain"
	"github.com/ZQCard/kbk-authorization/pkg/middleware/requestInfo"
	"github.com/ZQCard/kbk-authorization/pkg/utils/timeHelper"
)

type ApiEntity struct {
	BaseFields
	Domain string `gorm:"type:varchar(255);not null;comment:所在域"`
	Name   string `gorm:"type:varchar(255);not null;comment:名称"`
	Group  string `gorm:"type:varchar(255);not null;comment:分组"`
	Method string `gorm:"type:varchar(255);not null;comment:请求方式"`
	Path   string `gorm:"type:varchar(255);not null;comment:请求路径"`
}

func (ApiEntity) TableName() string {
	return "apis"
}

type ApiRepo struct {
	data *Data
	log  *log.Helper
}

func NewApiRepo(data *Data, logger log.Logger) biz.ApiRepo {
	repo := &ApiRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/api")),
	}
	return repo
}

// searchApiParam 搜索条件
func (repo ApiRepo) searchApiParam(ctx context.Context, params map[string]interface{}) *gorm.DB {
	conn := getDbWithDomain(ctx, repo.data.db).Model(&ApiEntity{})
	if Id, ok := params["id"]; ok && Id.(int64) != 0 {
		conn = conn.Where("id = ?", Id)
	}
	if Id, ok := params["neq_id"]; ok && Id.(int64) != 0 {
		conn = conn.Where("id != ?", Id)
	}
	if v, ok := params["name"]; ok && v.(string) != "" {
		conn = conn.Where("name LIKE ?", "%"+v.(string)+"%")
	}
	if v, ok := params["group"]; ok && v.(string) != "" {
		conn = conn.Where("`group` LIKE ?", "%"+v.(string)+"%")
	}
	if v, ok := params["method"]; ok && v.(string) != "" {
		conn = conn.Where("method LIKE ?", "%"+v.(string)+"%")
	}
	if v, ok := params["path"]; ok && v.(string) != "" {
		conn = conn.Where("path LIKE ?", "%"+v.(string)+"%")
	}
	return conn
}

func (repo ApiRepo) GetApiByParams(ctx context.Context, params map[string]interface{}) (record *ApiEntity, err error) {
	if len(params) == 0 {
		return &ApiEntity{}, v1.ErrorBadRequest("缺少搜索条件")
	}
	conn := repo.searchApiParam(ctx, params)
	if err = conn.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ApiEntity{}, v1.ErrorRecordNotFound("数据不存在")
		}
		return record, v1.ErrorSystemError("GetApiByParams First Error : %s", err.Error())
	}
	return record, nil
}

func (repo ApiRepo) ListApi(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Api, int64, error) {
	list := []*ApiEntity{}
	conn := repo.searchApiParam(ctx, params)
	err := conn.Scopes(Paginate(page, pageSize)).Order("`group` ASC").Find(&list).Error
	if err != nil {
		return nil, 0, v1.ErrorSystemError("ListApi Find Error : %s", err.Error())
	}

	rv := make([]*domain.Api, 0, len(list))
	for _, record := range list {
		api := toDomainApi(record)
		rv = append(rv, api)
	}
	return rv, repo.GetApiCount(ctx, params), nil
}

func (repo ApiRepo) GetApiCount(ctx context.Context, params map[string]interface{}) int64 {
	conn := repo.searchApiParam(ctx, params)
	count := int64(0)
	conn.Count(&count)
	return count
}

func (repo ApiRepo) ListApiAll(ctx context.Context) ([]*domain.Api, error) {
	list := []*ApiEntity{}
	conn := repo.searchApiParam(ctx, map[string]interface{}{})
	err := conn.Find(&list).Error
	if err != nil {
		return nil, v1.ErrorSystemError("ListApiAll Find Error : %s", err.Error())
	}

	rv := make([]*domain.Api, 0, len(list))
	for _, record := range list {
		api := toDomainApi(record)
		rv = append(rv, api)
	}
	return rv, nil
}

func (repo ApiRepo) CreateApi(ctx context.Context, api *domain.Api) (*domain.Api, error) {
	// 名称不得重复
	count := int64(0)
	repo.searchApiParam(ctx, map[string]interface{}{
		"name": api.Name,
	}).Count(&count)
	if count > 0 {
		return nil, v1.ErrorBadRequest("名称不得重复")
	}
	entity := &ApiEntity{}
	entity.Id = api.Id
	entity.Name = api.Name
	entity.Group = api.Group
	entity.Method = api.Method
	entity.Path = api.Path
	domain := ctx.Value(requestInfo.DomainKey)
	if domain != nil {
		entity.Domain = domain.(string)
	}
	if err := repo.data.db.Create(entity).Error; err != nil {
		return nil, v1.ErrorSystemError("CreateApi Create Error : %s", err.Error())
	}
	response := toDomainApi(entity)
	return response, nil
}

func (repo ApiRepo) UpdateApi(ctx context.Context, api *domain.Api) error {
	// 名称不得重复
	count := int64(0)
	repo.searchApiParam(ctx, map[string]interface{}{
		"name":   api.Name,
		"neq_id": api.Id,
	}).Count(&count)
	if count > 0 {
		return v1.ErrorBadRequest("名称不得重复")
	}
	// 根据Id查找记录
	record, err := repo.GetApiByParams(ctx, map[string]interface{}{
		"id": api.Id,
	})
	if err != nil {
		return err
	}
	// 更新字段
	record.Name = api.Name
	record.Group = api.Group
	record.Method = api.Method
	record.Path = api.Path
	domain := ctx.Value(requestInfo.DomainKey)
	if domain != nil {
		record.Domain = domain.(string)
	}
	if err := repo.data.db.Where("id = ?", record.Id).Save(record).Error; err != nil {
		return v1.ErrorSystemError("UpdateApi Save Error : %s", err.Error())
	}
	return nil
}

func (repo ApiRepo) GetApi(ctx context.Context, params map[string]interface{}) (*domain.Api, error) {
	// 根据Id查找记录
	record, err := repo.GetApiByParams(ctx, params)
	if err != nil {
		return nil, err
	}
	// 返回数据
	response := toDomainApi(record)
	return response, nil
}

func (repo ApiRepo) DeleteApi(ctx context.Context, domain *domain.Api) error {
	// 根据Id查找记录
	record, err := repo.GetApiByParams(ctx, map[string]interface{}{
		"id": domain.Id,
	})
	if err != nil {
		return err
	}
	if domain.Id != record.Id {
		return v1.ErrorRecordNotFound("数据不存在")
	}
	if err := getDbWithDomain(ctx, repo.data.db).Unscoped().Where("Id = ?", domain.Id).Delete(&ApiEntity{}).Error; err != nil {
		return v1.ErrorSystemError("DeleteApi Delete Error : %s", err.Error())
	}
	return nil
}

func toDomainApi(api *ApiEntity) *domain.Api {
	return &domain.Api{
		Id:        api.Id,
		Domain:    api.Domain,
		Name:      api.Name,
		Group:     api.Group,
		Method:    api.Method,
		Path:      api.Path,
		CreatedAt: timeHelper.FormatYMDHIS(&api.CreatedAt),
		UpdatedAt: timeHelper.FormatYMDHIS(&api.UpdatedAt),
	}
}
