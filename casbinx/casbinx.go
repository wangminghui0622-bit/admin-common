// Package casbinx 装配 admin-zero 的 casbin 执行器：
// 模型与 simple-admin 同构，策略存 admin_zero.casbin_rules（列名 ptype）。
package casbinx

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinRule 映射 casbin_rules 表（simple-admin 用 ptype 而非 p_type 列名）。
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
	V3    string `gorm:"column:v3"`
	V4    string `gorm:"column:v4"`
	V5    string `gorm:"column:v5"`
}

// ModelText 与 simple-admin 的 rbac_model.conf 一致。
const ModelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// casbinRulesDDL 与 core 的 init.sql 中 casbin_rules 建表语句一致。
// 空库首次启动时 init.sql 还没执行（InitDatabase 尚未被调用），
// 这里先按相同结构建表，避免适配器 LoadPolicy 因表不存在而启动失败（鸡生蛋问题）。
// 唯一索引是 init.sql INSERT IGNORE 去重的依据，不能省略。
const casbinRulesDDL = `
CREATE TABLE IF NOT EXISTS ` + "`casbin_rules`" + ` (
  ` + "`id`" + ` bigint NOT NULL AUTO_INCREMENT,
  ` + "`ptype`" + ` varchar(24) NOT NULL DEFAULT '',
  ` + "`v0`" + ` varchar(24) NOT NULL DEFAULT '',
  ` + "`v1`" + ` varchar(255) NOT NULL DEFAULT '',
  ` + "`v2`" + ` varchar(24) NOT NULL DEFAULT '',
  ` + "`v3`" + ` varchar(32) NOT NULL DEFAULT '',
  ` + "`v4`" + ` varchar(12) NOT NULL DEFAULT '',
  ` + "`v5`" + ` varchar(12) NOT NULL DEFAULT '',
  PRIMARY KEY (` + "`id`" + `),
  UNIQUE KEY ` + "`casbinrule_ptype_v0_v1_v2_v3_v4_v5`" + ` (` + "`ptype`" + `,` + "`v0`" + `,` + "`v1`" + `,` + "`v2`" + `,` + "`v3`" + `,` + "`v4`" + `,` + "`v5`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin`

// NewEnforcer 用 gorm adapter 连接 casbin_rules 表，启动时自动 LoadPolicy。
// 表结构由 init.sql 创建（ptype varchar(24) 等带唯一索引）；空库启动时本函数
// 用相同 DDL 先建表（CREATE TABLE IF NOT EXISTS，与 init.sql 语义一致）。
// 关闭适配器的 AutoMigrate，避免它把列改成 longtext 撞上索引报 1170。
func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	gormadapter.TurnOffAutoMigrate(db)
	if !db.Migrator().HasTable(&CasbinRule{}) {
		if err := db.Exec(casbinRulesDDL).Error; err != nil {
			return nil, err
		}
	}
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &CasbinRule{}, "casbin_rules")
	if err != nil {
		return nil, err
	}
	// NewEnforcer 的第一个字符串参数是模型文件路径，内存模型要用 NewModelFromString
	m, err := model.NewModelFromString(ModelText)
	if err != nil {
		return nil, err
	}
	return casbin.NewEnforcer(m, adapter)
}
