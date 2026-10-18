// gen 重新生成 admin-common/dal（model + query），生成配置与已提交的 dal 保持一致：
//   - Mode 0：不生成包级 Q 变量与 iXxxDo 接口
//   - FieldNullable：可空字段用指针
//   - json 标签策略：直接用列名（snake_case）
//   - 排除 pay_* 表：pay 查询集合由手写的 dal/query/pay_query.gen.go 维护（见该文件注释）
//
// 用法：go run ./cmd/gen -dsn "root:@tcp(127.0.0.1:3306)/admin_zero?..."
package main

import (
	"flag"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := flag.String("dsn",
		"root:@tcp(127.0.0.1:3306)/admin_zero?charset=utf8mb4&parseTime=True&loc=Local",
		"MySQL DSN，指向 admin_zero 库")
	flag.Parse()

	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db fail: %v", err)
	}

	cfg := gen.Config{
		OutPath:       "dal/query",
		ModelPkgPath:  "model",
		FieldNullable: true,
	}
	cfg.WithJSONTagNameStrategy(func(columnName string) string { return columnName })
	g := gen.NewGenerator(cfg)
	g.UseDB(db)

	// 生成除 pay_* 外的所有表（pay 的 Query 集合由手写的 dal/query/pay_query.gen.go 维护）
	tableNames, err := db.Migrator().GetTables()
	if err != nil {
		log.Fatalf("get tables fail: %v", err)
	}
	var tables []interface{}
	for _, name := range tableNames {
		if strings.HasPrefix(name, "pay_") {
			log.Printf("skip table %s (hand-written pay query)", name)
			continue
		}
		tables = append(tables, g.GenerateModel(name))
	}
	g.ApplyBasic(tables...)

	g.Execute()
}
