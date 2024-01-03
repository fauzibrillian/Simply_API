package handler

import (
	"library_api/features/user"
	"library_api/helper/jwt"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type UserController struct {
	srv user.Service
}

func New(s user.Service) user.Handler {
	return &UserController{
		srv: s,
	}
}

// Login implements user.Handler.
func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := uc.srv.Login(input.Phone, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]any{
					"message": "phone tidak ditemukan",
				})
			}
			if strings.Contains(err.Error(), "password salah") {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "password salah",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "phone tidak di temukan",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID, result.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var response = new(LoginResponse)
		response.ID = result.ID
		response.Name = result.Name
		response.Phone = result.Phone
		response.Role = result.Role
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Success Login Data",
			"data":    response,
		})
	}
}

// Register implements user.Handler.
func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input tidak sesuai",
			})
		}
		var inputProses = new(user.User)
		inputProses.Phone = input.Phone
		inputProses.Name = input.Name
		inputProses.Password = input.Password

		result, err := uc.srv.Register(*inputProses)
		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "email telah terdaftar",
			})
		}
		var response = new(RegisterResponse)
		response.ID = result.ID
		response.Phone = result.Phone
		response.Name = result.Name

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Register Data",
			"data":    response,
		})
	}
}