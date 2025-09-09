
package usecase

import (
	"context"
	"fmt"
	"mymodule/internal/domain"
	"mymodule/pkg/excelutil"

	"github.com/xuri/excelize/v2"
)

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new userUsecase.
func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// ImportUsersFromExcel reads user data from an Excel file and saves it to the repository.
func (uc *userUsecase) ImportUsersFromExcel(ctx context.Context, filePath string) error {
	// Assume the first sheet is the one with user data
	sheets, err := excelutil.GetSheetNames(filePath)
	if err != nil || len(sheets) == 0 {
		return fmt.Errorf("failed to get sheets from excel file: %w", err)
	}
	sheetName := sheets[0]

	// Read file with headers
	rows, err := excelutil.ReadExcelFileWithHeaders(filePath, sheetName)
	if err != nil {
		return fmt.Errorf("failed to read excel file with headers: %w", err)
	}

	users := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		// Basic validation
		name, hasName := row["Name"]
		email, hasEmail := row["Email"]
		if !hasName || !hasEmail {
			// Skip rows that don't have the required fields
			continue
		}

		users = append(users, &domain.User{
			Username: name,
			Email:    email,
			// Password can be auto-generated or handled differently
			Password: "default-password",
		})
	}

	if len(users) == 0 {
		return fmt.Errorf("no valid user data found in the excel file")
	}

	// Save users to the repository in a batch
	return uc.userRepo.CreateBatch(ctx, users)
}

// ExportUsersToExcel gets all users from the repository and generates an Excel file.
func (uc *userUsecase) ExportUsersToExcel(ctx context.Context) (*excelize.File, string, error) {
	users, err := uc.userRepo.GetAll(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get all users: %w", err)
	}

	sheetName := "Users"
	f, err := excelutil.CreateExcelFile(sheetName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create excel file: %w", err)
	}

	// Prepare headers
	headers := []string{"ID", "Username", "Email"}
	excelutil.SetHeaderRow(f, sheetName, headers)

	// Create and apply a style for the header
	headerStyle, _ := excelutil.CreateStyle(f, map[string]interface{}{
		"font":      map[string]interface{}{"bold": true, "color": "#FFFFFF"},
		"fill":      map[string]interface{}{"type": "pattern", "color": []string{"#4F81BD"}, "pattern": 1},
		"alignment": map[string]interface{}{"horizontal": "center"},
	})
	for i := 0; i < len(headers); i++ {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		excelutil.SetCellStyle(f, sheetName, cell, headerStyle)
	}

	// Write user data
	for i, user := range users {
		rowNum := i + 2 // Start from row 2
		rowData := []interface{}{user.ID, user.Username, user.Email}
		excelutil.WriteRow(f, sheetName, rowNum, rowData)
	}

	// Auto-fit columns
	excelutil.AutoFitColumns(f, sheetName)

	fileName := fmt.Sprintf("Users_Export_%d.xlsx",_time.Now().Unix())

	return f, fileName, nil
}
