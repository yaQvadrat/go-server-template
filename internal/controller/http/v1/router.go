package httpapi

import (
	"app/internal/service"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: setLogsFile()}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	
	authHandlers := newAuthRoutes(services.Account)
	auth := handler.Group("/auth")
	{
		auth.POST("/sign_up", authHandlers.signUp)
		auth.POST("/sign_in", authHandlers.signIn)
	}

	authMiddleware := &AuthMiddleware{services.Account}
	api := handler.Group("/api", authMiddleware.Auth)
	{
		api.GET("/test_admin", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
		api.GET("/test_user", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/logfile.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("httpapi - setLogsFile - os.OpenFile: %w", err))
	}
	return file
}
