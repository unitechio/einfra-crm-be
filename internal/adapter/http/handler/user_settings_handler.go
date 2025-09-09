
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
)

// UserSettingsHandler handles user settings-related requests.
type UserSettingsHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserSettingsHandler creates a new UserSettingsHandler.
func NewUserSettingsHandler(userUsecase domain.UserUsecase) *UserSettingsHandler {
	return &UserSettingsHandler{userUsecase: userUsecase}
}

// UpdateSettings handles the request to update user settings.
func (h *UserSettingsHandler) UpdateSettings(c *gin.Context) {
	var settings domain.UserSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("id") // Or get from authenticated user

	if err := h.userUsecase.UpdateUserSettings(c.Request.Context(), userID, settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
