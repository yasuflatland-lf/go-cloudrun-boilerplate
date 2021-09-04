package main

import (
	"context"
	"fmt"
	"github.com/glassonion1/logz"
	"golang.org/x/xerrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type (
	CloudSQL interface {
		Open(ctx context.Context, dsn string) (*gorm.DB, error)
		DB() *gorm.DB
		GenerateDSNLocal(name string, username string, password string, ip string, port int64) string
		GenerateDSNForCloudDB(name string, username string, password string, cloudSqlInstances string) string
	}

	cloudSQL struct {
		mu     sync.Mutex
		dsn    string
		db     *gorm.DB
		config *applicationConfig
	}
)

func NewCloudSQL(ctx context.Context) CloudSQL {
	c := &cloudSQL{}

	// Mutex for DB connection creation
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = GetApplicationConfig(ctx)

	// Build DSN to access the database
	if c.config.IsProduction() || c.config.IsDevelopment() {
		// Production or Development
		c.dsn = c.GenerateDSNForCloudDB(
			c.config.Name,
			c.config.UserName,
			c.config.Password,
			c.config.CloudSQLInstance,
		)
	} else {
		// Test environment
		c.dsn = c.GenerateDSNLocal(
			c.config.Name,
			c.config.UserName,
			c.config.Password,
			c.config.IP,
			c.config.Port,
		)
	}

	db, err := c.Open(ctx, c.dsn)
	if err != nil {
		logz.Criticalf(ctx, "%+v\n", xerrors.Errorf(": %w", err))
	}

	// Set DB
	c.db = db

	return c
}

func (c *cloudSQL) GenerateDSNLocal(name string, username string, password string, ip string, port int64) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, ip, port, name)
}

func (c *cloudSQL) GenerateDSNForCloudDB(name string, username string, password string, cloudSqlInstances string) string {
	return fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, cloudSqlInstances, name)
}

func (c *cloudSQL) Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Caches Prepared Statement
		// https://gorm.io/docs/performance.html#Caches-Prepared-Statement
		PrepareStmt: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil || db == nil {
		logz.Errorf(ctx, "error connect to db %+v", xerrors.Errorf(": %w", err))
		return nil, xerrors.Errorf(": %w", err)
	}

	// Set Context
	db.WithContext(ctx)

	// GetDomains generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		logz.Errorf(ctx, "error to fetch sqlDB %+v", xerrors.Errorf(": %w", err))
		return nil, xerrors.Errorf(": %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(c.config.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of Open connections to the database.
	sqlDB.SetMaxOpenConns(c.config.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Database handler
func (c *cloudSQL) DB() *gorm.DB {
	return c.db
}
