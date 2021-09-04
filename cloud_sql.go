package main

import (
	"context"
	"fmt"
	"github.com/glassonion1/logz"
	"github.com/golang-migrate/migrate"
	golang_migrate_mysql "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"golang.org/x/xerrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type (
	CloudSQL interface {
		GetDSN() string
		Open(ctx context.Context, dsn string) (*gorm.DB, error)
		DB() *gorm.DB
		GenerateDSNLocal(name string, username string, password string, ip string, port int64) string
		GenerateDSNForCloudDB(name string, username string, password string, cloudSqlInstances string) string
		StartMigrations(ctx context.Context) error
		RollbackLastMigrations(ctx context.Context) error
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
		logz.Criticalf(ctx, "Failed to open database connection. %+v\n", xerrors.Errorf(": %w", err))
	}

	// Set DB
	c.db = db

	return c
}

func (c *cloudSQL) GetDSN() string {
	return c.dsn
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
		return nil, xerrors.Errorf(": %w", err)
	}

	// Set Context
	db.WithContext(ctx)

	// GetDomains generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
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

// https://github.dev/elsennov/guitar_collection/blob/1f869cd16ddeab778c42fa54d72cba5bdd870305/console/migrations.go
func (c *cloudSQL) StartMigrations(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//db, err := sql.Open("golang_migrate_mysql", dsn)
	//if err != nil {
	//	log.Println("Error sql.Open")
	//	return xerrors.Errorf(": %w", err)
	//}
	db, err := c.DB().DB()
	if err != nil {
		return xerrors.Errorf("Error db.DB() : %w", err)
	}

	if err := db.Ping(); err != nil {
		return xerrors.Errorf("could not ping DB...  : %w", err)
	}

	driver, err := golang_migrate_mysql.WithInstance(db, &golang_migrate_mysql.Config{})
	if err != nil {
		xerr := xerrors.Errorf("Error mysql.WithInstance : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		c.config.Name,
		driver,
	)
	if err != nil {
		xerr := xerrors.Errorf("Error migrate.NewWithDatabaseInstance : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	if migration != nil {
		err := migration.Steps(1)
		if err != nil {
			xerr := xerrors.Errorf(": %w", err)
			logz.Errorf(ctx, " %+v", xerr)
			return xerr
		}
	}
	return nil
}

func (c *cloudSQL) RollbackLastMigrations(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	db, err := c.DB().DB()
	if err != nil {
		xerr := xerrors.Errorf("Error db.DB() : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	if err := db.Ping(); err != nil {
		xerr := xerrors.Errorf("could not ping DB...  : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	driver, err := golang_migrate_mysql.WithInstance(db, &golang_migrate_mysql.Config{})
	if err != nil {
		xerr := xerrors.Errorf("Error mysql.WithInstance : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		c.config.Name,
		driver,
	)
	if err != nil {
		xerr := xerrors.Errorf("Error migrate.NewWithDatabaseInstance : %w", err)
		logz.Errorf(ctx, " %+v", xerr)
		return xerr
	}

	if migration != nil {
		err := migration.Steps(-1)
		if err != nil {
			xerr := xerrors.Errorf(": %w", err)
			logz.Errorf(ctx, " %+v", xerr)
			return xerr
		}
	}
	return nil
}
