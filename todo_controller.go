package main

import (
	"context"
	"github.com/glassonion1/logz"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
	"net/http"
	"strconv"
)

type (
	TodoController interface {
		List(c echo.Context) error
		Get(c echo.Context) error
	}
	todoController struct{
		todoService TodoService
	}
)

func NewTodoController(ctx context.Context) TodoController {
	return &todoController{
		todoService: NewTodoService(ctx),
	}
}

func (t *todoController) List(c echo.Context) error {
	status, err := strconv.ParseBool(c.QueryParam("status"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			xerrors.Errorf("Missing parameter : status : %w", err))
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return xerrors.Errorf("Missing parameter : page : %w", err)
	}

	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	if err != nil {
		return xerrors.Errorf("Missing parameter : pagesize : %w", err)
	}

	results, rows, err := t.todoService.List(status, page, pagesize, "updated_at DESC")
	if err != nil {
		return xerrors.Errorf("Fetching List : %w", err)
	}

	logz.Infof(c.Request().Context(),"fetched row: %+v", rows)

	return c.JSON(http.StatusOK, results)
}

func (t *todoController) Get(c echo.Context) error {
	id := c.Param("id")
	logz.Errorf(context.Background(), "%+v", id)
	return c.JSON(http.StatusOK, &[]Todo{})
}
