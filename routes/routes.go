package routes

import (
	"library_api/config"
	"library_api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uc user.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	RouteUser(e, uc)
}

func RouteUser(e *echo.Echo, uc user.Handler) {
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.PATCH("/user/password/:id", uc.ResetPassword(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.PATCH("/user/:id", uc.UpdateUser(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.DELETE("/user/:id", uc.Delete(), echojwt.JWT([]byte(config.InitConfig().JWT)))
}
