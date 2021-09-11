package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTodoController(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	t.Run("List", eachTestWrapper(func(t *testing.T) {
		// Setup
		router := NewRouter(ctx)
		q := make(url.Values)
		q.Set("status", "false")
		q.Set("page", "1")
		q.Set("pagesize", "10")
		req := httptest.NewRequest(http.MethodGet, "/todos?"+q.Encode(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())

		var responceJson = &[]Todo{}
		err := json.Unmarshal([]byte(rec.Body.String()), &responceJson)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotNil(t, responceJson)
		}

	}))

	t.Run("Create Get and Delete", eachTestWrapper(func(t *testing.T) {
		t.Run("1 Create", func(t *testing.T) {
			// Create a todo
			todo := &Todo{}
			// Generate Fake data
			if err := faker.FakeData(&todo); err != nil {
				assert.Nil(t, err)
			}
			// Stringify object
			todoStr, err := json.Marshal(todo)
			if err != nil {
				assert.Nil(t, err)
			}

			// fmt.Printf("%+v", string(todoStr))
			// Setup
			router := NewRouter(ctx)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(todoStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
		})

		t.Run("2 Get", func(t *testing.T) {
			// Setup
			router := NewRouter(ctx)

			// ID 1 record Should be created in the above Create
			req := httptest.NewRequest(http.MethodGet, "/1", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
		})

		t.Run("3 Delete", func(t *testing.T) {
			// Setup
			router := NewRouter(ctx)

			// ID 1 record Should be created in the above Create
			req := httptest.NewRequest(http.MethodDelete, "/1", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{ \"RowsAffected\": 1 }", rec.Body.String())
		})

		t.Run("4 Make sure the data is deleted", func(t *testing.T) {
			// Setup
			router := NewRouter(ctx)

			// ID 1 record Should be created in the above Create
			req := httptest.NewRequest(http.MethodGet, "/1", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{}\n", rec.Body.String())
		})
	}))

	t.Run("Update", eachTestWrapper(func(t *testing.T) {
		t.Run("1 Create", func(t *testing.T) {
			// Create a todo
			todo := &Todo{}
			// Generate Fake data
			if err := faker.FakeData(&todo); err != nil {
				assert.Nil(t, err)
			}
			// Stringify object
			todoStr, err := json.Marshal(todo)
			if err != nil {
				assert.Nil(t, err)
			}

			// Setup
			router := NewRouter(ctx)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(todoStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
			// fmt.Printf("%+v", rec.Body.String())
		})

		t.Run("2 Get and Update", func(t *testing.T) {
			// Setup
			router := NewRouter(ctx)
			todo := &Todo{
				ID:     1,
				Slug:   "test-slug",
				Task:   "Changed",
				Status: true,
			}

			// Stringify object
			todoStr, err := json.Marshal(todo)
			if err != nil {
				assert.Nil(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(todoStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
			fmt.Printf("%+v", rec.Body.String())
		})

		t.Run("3 Update Fail", func(t *testing.T) {
			// Setup
			router := NewRouter(ctx)
			todo := &Todo{
				ID:     2,
				Slug:   "test-slug",
				Task:   "Changed",
				Status: true,
			}

			// Stringify object
			todoStr, err := json.Marshal(todo)
			if err != nil {
				assert.Nil(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(todoStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
			//fmt.Printf("%+v", rec.Body.String())
		})
	}))
}
