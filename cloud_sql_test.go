package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test MySQL Smoke
func TestCloudSQL(t *testing.T) {
	t.Helper()
	t.Parallel()
	ctx := context.Background()

	t.Run("GORM Connection Open Test", func(t *testing.T) {
		t.Parallel()

		dao := NewCloudSQL(ctx)
		db := dao.DB()
		assert.NotNil(t, db)
	})

	t.Run("Migration create and delete", func(t *testing.T) {
		t.Parallel()

		dao := NewCloudSQL(ctx)
		db := dao.DB()
		assert.NotNil(t, db)

		err := dao.StartMigrations(ctx)
		assert.Nil(t, err)

		err = dao.RollbackLastMigrations(ctx)
		assert.Nil(t, err)

	})
	// Remove comments to generate model in the database automatically.
	//t.Run("Generate Models From Tables", func(t *testing.T) {
	//	seedDataPath, _ := os.Getwd()
	//	err = converter.NewTable2Struct().
	//		SavePath(seedDataPath + "/model.go").
	//		Dsn(dsn).
	//		Run()
	//})
}
