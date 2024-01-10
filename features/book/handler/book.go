package handler

import (
	"context"
	"errors"
	"library_api/features/book"
	"library_api/helper/cld"
	"net/http"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type BookHandler struct {
	s      book.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s book.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) book.Handler {
	return &BookHandler{
		s:      s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

// AddBook implements book.Handler.
func (bh *BookHandler) AddBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(BookRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}
		formHeader, err := c.FormFile("picture")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				inputProcess := &book.Book{
					Tittle:    input.Tittle,
					Publisher: input.Publisher,
					Author:    input.Author,
					Picture:   "",
					Category:  input.Category,
					Stock:     uint(input.Stock),
				}

				result, err := bh.s.AddBook(c.Get("user").(*golangjwt.Token), *inputProcess)

				if err != nil {
					c.Logger().Error("ERROR Register, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "terdaftar") {
						statusCode = http.StatusBadRequest
						message = "data yang diinputkan sudah terdaftar ada sistem"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(BookResponse)
				response.ID = result.ID
				response.Tittle = result.Tittle
				response.Publisher = result.Publisher
				response.Author = result.Author
				response.Picture = result.Picture
				response.Category = result.Category
				response.Stock = int(result.Stock)

				return c.JSON(http.StatusCreated, map[string]any{
					"message": "Success Create Book Data",
					"data":    response,
				})

			}
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formheader error",
				})

		}
		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formfile error",
				})
		}

		link, err := cld.UploadImage(bh.cl, bh.ct, formFile, bh.folder)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "harap pilih gambar",
					"data":    nil,
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"message": "kesalahan pada server",
					"data":    nil,
				})
			}
		}

		var inputProcess = new(book.Book)
		inputProcess.Tittle = input.Tittle
		inputProcess.Publisher = input.Publisher
		inputProcess.Author = input.Author
		inputProcess.Tittle = input.Tittle
		inputProcess.Category = input.Category
		inputProcess.Picture = link
		inputProcess.Stock = uint(input.Stock)

		result, err := bh.s.AddBook(c.Get("user").(*golangjwt.Token), *inputProcess)
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())

			if strings.Contains(err.Error(), "terdaftar") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data duplicate",
				})
			}

			if strings.Contains(err.Error(), "tidak memiliki izin") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "tidak memiliki izin",
				})
			} else if strings.Contains(err.Error(), "tidak memiliki izin") {
				return c.JSON(http.StatusForbidden, map[string]any{
					"message": "tidak memiliki izin",
				})
			}

			return c.JSON(http.StatusForbidden, map[string]any{
				"message": "Tidak memiliki izin",
			})
		}

		var response = new(BookResponse)
		response.ID = result.ID
		response.Tittle = result.Tittle
		response.Publisher = result.Publisher
		response.Author = result.Author
		response.Picture = result.Picture
		response.Category = result.Category
		response.Stock = int(result.Stock)

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Success Created Book Data",
			"data":    response,
		})
	}
}