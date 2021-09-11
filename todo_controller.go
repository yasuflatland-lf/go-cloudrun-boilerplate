package main

import (
	"context"
	"fmt"
	"github.com/glassonion1/logz"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
	"net/http"
	"strconv"
	"time"
)

type (
	TodoController interface {
		List(c echo.Context) error
		Get(c echo.Context) error
		Create(c echo.Context) error
		Delete(c echo.Context) error
		Update(c echo.Context) error
	}
	todoController struct {
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
		errx := xerrors.Errorf("Missing parameter : status : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		errx := xerrors.Errorf("Missing parameter : page : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	if err != nil {
		errx := xerrors.Errorf("Missing parameter : pagesize : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	todos, rows, err := t.todoService.List(status, page, pagesize, "updated_at DESC")
	if err != nil {
		errx := xerrors.Errorf("Fetch List : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	logz.Infof(c.Request().Context(), "fetched row: %+v", rows)

	return c.JSON(http.StatusOK, todos)
}

func (t *todoController) Get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errx := xerrors.Errorf("Missing parameter : id : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	todo, err := t.todoService.Get(id)
	if err != nil {
		errx := xerrors.Errorf("Get todo: id %+w, %+w", id, err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusOK, errx)
	}

	return c.JSON(http.StatusOK, todo)
}

func (t *todoController) Create(c echo.Context) error {
	paramObj := &Todo{}
	if err := c.Bind(paramObj); err != nil {
		errx := xerrors.Errorf("Failed to bind parameter into todo object %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	// Create
	todo, err := t.todoService.Create(&Todo{
		Slug:      paramObj.Slug,
		Task:      paramObj.Task,
		Status:    paramObj.Status,
		CreatedAt: time.Time.UTC(time.Now()),
		UpdatedAt: time.Time.UTC(time.Now()),
	})

	if err != nil {
		errx := xerrors.Errorf("Create todo : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	return c.JSON(http.StatusOK, todo)
}

func (t *todoController) Delete(c echo.Context) error {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errx := xerrors.Errorf("Missing parameter : id : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	// Delete
	rowsAffected, err := t.todoService.Delete(ID)

	if err != nil {
		errx := xerrors.Errorf("Delete todo : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	return c.String(http.StatusOK,
		fmt.Sprintf("{ \"RowsAffected\": %d }", rowsAffected))
}

func (t *todoController) Update(c echo.Context) error {
	paramObj := &Todo{}
	if err := c.Bind(paramObj); err != nil {
		errx := xerrors.Errorf("Failed to bind parameter into todo object %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	orgTodo, err := t.todoService.Get(paramObj.ID)
	if err != nil {
		errx := xerrors.Errorf("ID %+w does not exist. %+w", paramObj.ID, err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	// Update
	todo, err := t.todoService.Update(&Todo{
		ID:        orgTodo.ID,
		Slug:      paramObj.Slug,
		Task:      paramObj.Task,
		Status:    paramObj.Status,
		UpdatedAt: time.Time.UTC(time.Now()),
		CreatedAt: orgTodo.CreatedAt,
	})

	if err != nil {
		errx := xerrors.Errorf("Update todo : %+w", err)
		logz.Errorf(c.Request().Context(), "%+v", errx)
		return echo.NewHTTPError(http.StatusBadRequest, errx)
	}

	return c.JSON(http.StatusOK, todo)
}
