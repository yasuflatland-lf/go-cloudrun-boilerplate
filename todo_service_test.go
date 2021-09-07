package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
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

		ID := createdTodo.ID
		_, err = todoService.Get(ID)
		assert.Nil(t, err)
		// Should be no data stored in the database.
		assert.Error(t, errors.New("record not found"))

		ids, err := todoService.Delete(createdTodo.ID)
		assert.Nil(t, err)
		assert.Equal(t, ids, int64(1))

		getTodo, err := todoService.Get(ID)
		assert.NotNil(t, err)
		assert.Nil(t, getTodo)

	}))

	t.Run("List and CreateInBatches", eachTestWrapper(func(t *testing.T) {
		batchAmount := 10

		// Create Dummy Todo array
		var todos = []Todo{}
		for i := 0; i < batchAmount; i++ {
			todo := Todo{}
			// Generate Fake data
			if err := faker.FakeData(&todo); err != nil {
				assert.Nil(t, err)
			}
			todo.Status = true
			todos = append(todos, todo)
		}

		_, err := todoService.CreateInBatches(todos)
		assert.Nil(t, err)

		results, rows, err := todoService.List(true, 1, 20, "updated_at asc")
		assert.Nil(t, err)
		assert.NotEmpty(t, results)
		assert.Equal(t, batchAmount, rows)
	}))

	t.Run("List Random Count True", eachTestWrapper(func(t *testing.T) {
		batchAmount := 10

		// Create Dummy Todo array
		var todos = []Todo{}
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		var trueAmount = 0
		for i := 0; i < batchAmount; i++ {
			todo := Todo{}
			// Generate Fake data
			if err := faker.FakeData(&todo); err != nil {
				assert.Nil(t, err)
			}
			var stat = r1.Intn(2)
			todo.Status = false
			if 0 < stat {
				todo.Status = true
				trueAmount++
			}

			todos = append(todos, todo)
		}

		_, err := todoService.CreateInBatches(todos)
		assert.Nil(t, err)

		results, rows, err := todoService.List(true, 1, 20, "updated_at asc")
		assert.Nil(t, err)
		assert.NotEmpty(t, results)
		assert.Equal(t, trueAmount, rows)
	}))

	t.Run("Update", eachTestWrapper(func(t *testing.T) {
		batchAmount := 10

		// Create Dummy Todo array
		var todos = []Todo{}
		for i := 0; i < batchAmount; i++ {
			todo := Todo{}
			// Generate Fake data
			if err := faker.FakeData(&todo); err != nil {
				assert.Nil(t, err)
			}
			todo.Status = true
			todo.CreatedAt = time.Time.UTC(time.Now())
			todo.UpdatedAt = time.Time.UTC(time.Now())
			todos = append(todos, todo)
		}

		createdTodos, err := todoService.CreateInBatches(todos)
		assert.Nil(t, err)

		updateTodo := createdTodos[0]
		updateTodo.Status = false
		updateTodo.Task = "Changed"
		retTodo, err := todoService.Update(&updateTodo)
		assert.NotNil(t, retTodo)

		results, err := todoService.Get(retTodo.ID)
		assert.Nil(t, err)
		assert.NotEmpty(t, results)
		assert.Equal(t, false, results.Status)
		assert.Equal(t, "Changed", results.Task)
	}))

}
