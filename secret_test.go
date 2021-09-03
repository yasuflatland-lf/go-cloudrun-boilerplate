package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSecret(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	PROJECT_UUID := os.Getenv("PROJECT_UUID")
	PROJECT_ID := os.Getenv("PROJECT_ID")

	t.Run("Create, Delete and Get A Secret", func(t *testing.T) {
		storeString := "test"

		// Create instance
		obj := NewSecret(ctx, PROJECT_ID, PROJECT_UUID)
		assert.NotNil(t, obj)

		// Generate Secret Id
		secretId := fmt.Sprint(uuid.New())

		// Create Secret
		secret, err := obj.CreateSecret(secretId, storeString)
		assert.Nil(t, err)
		assert.NotNil(t, secret)

		// Get Secret
		ret, err := obj.GetSecret(secretId)
		assert.Nil(t, err)
		assert.Equal(t, storeString, ret)

		// Delete Secret
		err = obj.DeleteSecret(secretId)
		assert.Nil(t, err)

	})
}
