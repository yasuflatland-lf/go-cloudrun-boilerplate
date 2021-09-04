package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationConfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Fetching test environment values", func(t *testing.T) {
		t.Parallel()
		var config = GetApplicationConfig(ctx)

		assert.NotNil(t, config)
		assert.NotEmpty(t, config.ProjectId)
	})
}
