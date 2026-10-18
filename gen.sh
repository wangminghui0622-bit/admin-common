#!/usr/bin/env bash
# 重新生成 admin-common/dal（model + query）。
# 需要本地 MySQL 已启动；默认连 127.0.0.1:3306 的 admin_zero（开发库 root 无密码）。
# 用法：
#   ./gen.sh
#   ./gen.sh "root:pw@tcp(127.0.0.1:3306)/admin_zero?charset=utf8mb4&parseTime=True&loc=Local"
set -e
cd "$(dirname "$0")"
DSN="${1:-root:@tcp(127.0.0.1:3306)/admin_zero?charset=utf8mb4&parseTime=True&loc=Local}"
go run ./cmd/gen -dsn "$DSN"
