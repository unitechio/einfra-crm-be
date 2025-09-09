package usecase

import (
	"mymodule/internal/domain"

	"github.com/gin-gonic/gin"
)

type auditUseCase struct {
	auditRepo domain.AuditRepository
}

func NewAuditUseCase(auditRepo domain.AuditRepository) domain.AuditService {
	return &auditUseCase{
		auditRepo: auditRepo,
	}
}

func (uc *auditUseCase) Log(c *gin.Context, entry domain.AuditEntry) (domain.AuditEntry, error) {
	return uc.auditRepo.Add(c, entry)
}

func (uc *auditUseCase) GetAll(c *gin.Context) ([]domain.AuditEntry, error) {
	return uc.auditRepo.GetAll(c)
}

func (uc *auditUseCase) GetByID(c *gin.Context, id string) (domain.AuditEntry, error) {
	return uc.auditRepo.GetByID(c, id)
}

func (uc *auditUseCase) Update(c *gin.Context, id string, entry domain.AuditEntry) (domain.AuditEntry, error) {
	return uc.auditRepo.Update(c, id, entry)
}

func (uc *auditUseCase) Delete(c *gin.Context, id string) error {
	return uc.auditRepo.Delete(c, id)
}
