package data

import (
	"context"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/conf"
)

type CasbinEntity struct {
	ID     int64  `json:"id"`
	PType  string `json:"ptype" gorm:"column:p_type" `
	Role   string `json:"role_name" gorm:"column:v0" `
	Path   string `json:"path" gorm:"column:v1" `
	Method string `json:"method" gorm:"column:v2" `
}

func (CasbinEntity) TableName() string {
	return "casbin_rule"
}

type CasbinRepo struct {
	data     *Data
	enforcer *casbin.Enforcer
	log      *log.Helper
}

func NewCasbinRepo(data *Data, conf *conf.Bootstrap, logger log.Logger) biz.CasbinRepo {
	// 初始化基础数据库 casbin权限控制策略,连接基础库
	db, _ := gormadapter.NewAdapterByDB(data.db)
	// 加载权限配置文件
	m, err := model.NewModelFromString(conf.Casbin.RbacModel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	enforcer, err := casbin.NewEnforcer(m, db)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// 开启权限日志
	if conf.Env.Mode == "dev" {
		enforcer.EnableLog(true)
	}
	// 从DB加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf(err.Error())
	}

	repo := &CasbinRepo{
		data:     data,
		enforcer: enforcer,
		log:      log.NewHelper(log.With(logger, "module", "data/Casbin")),
	}
	return repo
}

func (r CasbinRepo) SetRolesForUser(ctx context.Context, username string, Roles []string) (bool, error) {
	domain := getDomain(ctx)
	// 检查角色是否存在
	if !r.checkRoleExist(ctx, domain, Roles) {
		return false, v1.ErrorRecordNotFound("数据不存在")
	}
	// 删除用户所有角色
	success, err := r.enforcer.DeleteRolesForUser(username, domain)
	if err != nil {
		return false, v1.ErrorSystemError("SetRolesForUser DeleteRolesForUser Error : %s", err.Error())
	}
	// 添加用户角色
	success, err = r.enforcer.AddRolesForUser(username, Roles, domain)
	if err != nil {
		return false, v1.ErrorSystemError("SetRolesForUser AddRolesForUser Error : %s", err.Error())
	}
	return success, nil
}

func (r CasbinRepo) GetRolesForUser(ctx context.Context, username string) ([]string, error) {
	domain := getDomain(ctx)
	roles, err := r.enforcer.GetRolesForUser(username, domain)
	if err != nil {
		return nil, v1.ErrorSystemError("GetRolesForUser GetRolesForUser Error : %s", err.Error())
	}
	return roles, nil
}

func (r CasbinRepo) GetUsersForRole(ctx context.Context, Casbin string) ([]string, error) {
	domain := getDomain(ctx)
	users, err := r.enforcer.GetUsersForRole(Casbin, domain)
	if err != nil {
		return nil, v1.ErrorSystemError("GetUsersForRole GetUsersForRole Error : %s", err.Error())
	}
	return users, nil
}

func (r CasbinRepo) DeleteRoleForUser(ctx context.Context, username string, Casbin string) (bool, error) {
	domain := getDomain(ctx)
	success, err := r.enforcer.DeleteRoleForUser(username, Casbin, domain)
	if err != nil {
		return false, v1.ErrorSystemError("DeleteRoleForUser DeleteRoleForUser Error : %s", err.Error())
	}
	return success, nil
}

func (r CasbinRepo) DeleteRolesForUser(ctx context.Context, username string) (bool, error) {
	domain := getDomain(ctx)
	success, err := r.enforcer.DeleteRolesForUser(username, domain)
	if err != nil {
		return false, v1.ErrorSystemError("DeleteRolesForUser DeleteRolesForUser Error : %s", err.Error())
	}
	return success, nil
}

func (r CasbinRepo) GetPolicies(ctx context.Context, role string) ([]*biz.PolicyRules, error) {
	domain := getDomain(ctx)
	rules := []*biz.PolicyRules{}

	// 检查角色是否存在
	if !r.checkRoleExist(ctx, domain, []string{role}) {
		return rules, v1.ErrorBadRequest("角色不存在")
	}

	// 查询已有策略规则
	policies := r.enforcer.GetFilteredPolicy(0, role, domain)
	for _, v := range policies {
		rules = append(rules, &biz.PolicyRules{
			Path:   v[2],
			Method: v[3],
		})
	}

	return rules, nil
}

func (r CasbinRepo) UpdatePolicies(ctx context.Context, role string, rules []biz.PolicyRules) (bool, error) {
	domain := getDomain(ctx)
	// 检查角色是否存在
	if !r.checkRoleExist(ctx, domain, []string{role}) {
		return false, v1.ErrorBadRequest("角色不存在")
	}
	policies := [][]string{}
	for _, v := range rules {
		// method需要为全部大写
		policies = append(policies, []string{role, domain, v.Path, strings.ToUpper(v.Method)})
	}
	// 移除已有策略规则
	_, err := r.enforcer.RemoveFilteredPolicy(0, role, domain)
	if err != nil {
		return false, v1.ErrorSystemError("UpdatePolicies RemoveFilteredPolicy Error : %s", err.Error())
	}

	success, err := r.enforcer.AddPolicies(policies)
	if err != nil {
		return false, v1.ErrorSystemError("UpdatePolicies AddPolicies Error : %s", err.Error())
	}
	if !success {
		return false, v1.ErrorBadRequest("存在相同api,添加失败")
	}
	return true, nil
}

// 检查角色是否存在
func (r CasbinRepo) checkRoleExist(ctx context.Context, domain string, role []string) bool {
	roleCount := int64(len(role))
	if roleCount == 0 {
		return false
	}
	var count int64
	if len(role) == 1 {
		getDbWithDomain(ctx, r.data.db).Model(RoleEntity{}).Where("name = ? AND domain = ?", role[0], domain).Count(&count)
	} else {
		getDbWithDomain(ctx, r.data.db).Model(RoleEntity{}).Where("name IN (?) AND domain = ?", role, domain).Count(&count)
	}
	return count == roleCount
}

func (a CasbinRepo) CheckAuthorization(ctx context.Context, sub, obj, act string) (bool, error) {
	domain := getDomain(ctx)
	res, err := a.enforcer.Enforce(sub, domain, obj, act)
	return res, err
}
