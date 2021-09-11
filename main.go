package main

import (
	"context"
	"github.com/glassonion1/logz"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

func main() {
	ctx := context.Background()
	config := GetApplicationConfig(ctx)

	logz.InitTracer()

	router := NewRouter(ctx)

	// Start server
	router.Logger.Fatal(router.Start(":" + strconv.Itoa(config.Port)))
}

func NewRouter(ctx context.Context) *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	todoController := NewTodoController(ctx)

	// Routes
	e.GET("/todos", todoController.List)
	e.GET("/:id", todoController.Get)
	e.POST("/", todoController.Create)
	e.DELETE("/:id", todoController.Delete)
	e.PUT("/", todoController.Update)

	return e
}
