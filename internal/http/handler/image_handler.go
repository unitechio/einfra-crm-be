
package http

import (
	"fmt"
	"github.com/unitechio/einfra-be/internal/constants"
	"github.com/unitechio/einfra-be/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ImageHandler handles HTTP requests for image operations.
type ImageHandler struct {
	imageUsecase domain.ImageUsecase
}

// NewImageHandler registers the image-related routes.
func NewImageHandler(e *echo.Group, uc domain.ImageUsecase) {
	h := &ImageHandler{imageUsecase: uc}
	e.POST("/images/upload", h.UploadImage)
	// The route to serve images is set up in main.go
}

// UploadImage handles the image upload request.
func (h *ImageHandler) UploadImage(c echo.Context) error {
	// 1. Get the user from the context (set by AuthMiddleware)
	user, ok := c.Request().Context().Value(constants.UserContextKey).(*domain.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	// 2. Get the file from the form
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "image file not found in request")
	}

	// 3. Call the usecase to handle the upload
	image, err := h.imageUsecase.UploadImage(c.Request().Context(), user.ID, file)
	if err != nil {
		// More specific error handling could be implemented here
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to upload image: %v", err))
	}

	// 4. Return the details of the uploaded image
	return c.JSON(http.StatusCreated, image)
}
