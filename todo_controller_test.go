package main

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	t.Run("Get", func(t *testing.T) {
		// Setup
		e := NewRouter(ctx)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
