// Package migrate 提供 golang-migrate 的封装，支持进程内自动迁移。
package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run 执行迁移到最新版本。
// sourceURL 形如 "file://./migrations"，databaseURL 形如 "postgres://user:pwd@host:port/db?sslmode=disable"
func Run(sourceURL, databaseURL string) error {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("init migrate failed: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up failed: %w", err)
	}
	return nil
}

// Version 返回当前迁移版本。
func Version(sourceURL, databaseURL string) (uint, bool, error) {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return 0, false, err
	}
	defer func() { _, _ = m.Close() }()
	return m.Version()
}