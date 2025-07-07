package api

import (
	"github.com/labstack/echo/v4"
)

func CreateRouters(e *echo.Echo, handler *Handler) {

	e.GET("/tokens", handler.GetTokens)
	e.POST("/refresh", handler.RefreshTokens)

	// Группа для защищенных эндпоинтов
	protected := e.Group("")
	// Применяем наш AuthMiddleware к этой группе
	protected.Use(AuthMiddleware)
	protected.GET("/me", handler.GetMe)
	protected.POST("/logout", handler.Logout)

}
