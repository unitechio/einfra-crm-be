
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mymodule/internal/domain"
	"mymodule/internal/usecase"
)

// SystemSettingHandler handles HTTP requests related to system settings.
type SystemSettingHandler struct {
	uc *usecase.SystemSettingUseCase
}

// NewSystemSettingHandler creates a new instance of SystemSettingHandler.
func NewSystemSettingHandler(uc *usecase.SystemSettingUseCase) *SystemSettingHandler {
	return &SystemSettingHandler{uc: uc}
}

// CreateSystemSetting handles the creation of a new system setting.
func (h *SystemSettingHandler) CreateSystemSetting(c *gin.Context) {
	var setting domain.SystemSetting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSetting, err := h.uc.CreateSystemSetting(&setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdSetting)
}

// GetSystemSettingByKey handles the retrieval of a system setting by its key.
func (h *SystemSettingHandler) GetSystemSettingByKey(c *gin.Context) {
	key := c.Param("key")
	setting, err := h.uc.GetSystemSettingByKey(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "system setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// GetAllSystemSettings handles the retrieval of all system settings.
func (h *SystemSettingHandler) GetAllSystemSettings(c *gin.Context) {
	settings, err := h.uc.GetAllSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// GetSystemSettingsByCategory handles the retrieval of all system settings of a specific category.
func (h *SystemSettingHandler) GetSystemSettingsByCategory(c *gin.Context) {
    category := c.Param("category")
	settings, err := h.uc.GetSystemSettingsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}


// UpdateSystemSetting handles the update of an existing system setting.
func (h *SystemSettingHandler) UpdateSystemSetting(c *gin.Context) {
	var setting domain.SystemSetting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSetting, err := h.uc.UpdateSystemSetting(&setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSetting)
}

// DeleteSystemSetting handles the deletion of a system setting by its ID.
func (h *SystemSettingHandler) DeleteSystemSetting(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.DeleteSystemSetting(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
