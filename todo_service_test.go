package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTodoService(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	var todoService = NewTodoService(ctx)

	t.Run("Create and Delete", eachTestWrapper(func(t *testing.T) {

		todo := &Todo{
			Task:   "test task",
			Status: false,
		}
		createdTodo, err := todoService.Create(todo)

		assert.Nil(t, err)
		assert.NotNil(t, createdTodo)
		fmt.Printf("data %+v", createdTodo)

		ID := int(createdTodo.ID)
		_, err = todoService.Get(ID)
		assert.Nil(t, err)
		// Should be no data stored in the database.
		assert.Error(t, errors.New("record not found"))

		ids, err := todoService.Delete(createdTodo)
		assert.Nil(t, err)
		assert.Equal(t, ids, int64(1))

		getTodo, err := todoService.Get(ID)
		assert.NotNil(t, err)
		assert.Nil(t, getTodo)

	}))
}
