
package http

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/unitechio/einfra-be/internal/usecase"
)

// ExcelHandler handles HTTP requests for Excel import/export.
type ExcelHandler struct {
	userUsecase usecase.UserUsecase
}

// NewExcelHandler registers the Excel import/export routes.
func NewExcelHandler(group *echo.Group, userUsecase usecase.UserUsecase) {
	h := &ExcelHandler{userUsecase: userUsecase}
	group.POST("/users/import", h.ImportUsers)
	group.GET("/users/export", h.ExportUsers)
}

// ImportUsers handles the Excel file upload and user import.
func (h *ExcelHandler) ImportUsers(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file not found in request")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file")
	}
	defer src.Close()

	// We need to save the file temporarily to pass its path to the usecase.
	// In a real application, you might use a shared volume or object storage.
	tempFilePath := "temp-" + file.Filename
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create temp file")
	}
	defer dst.Close()
	defer os.Remove(tempFilePath) // Clean up the temp file

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save temp file")
	}

	if err := h.userUsecase.ImportUsersFromExcel(c.Request().Context(), tempFilePath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to import users: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Users imported successfully"})
}

// ExportUsers handles the request to export users to an Excel file.
func (h *ExcelHandler) ExportUsers(c echo.Context) error {
	file, filename, err := h.userUsecase.ExportUsersToExcel(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to export users: "+err.Error())
	}

	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)

	if err := file.Write(c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to write excel file to response")
	}

	return c.NoContent(http.StatusOK)
}
