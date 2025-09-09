// Package excelutil provides a comprehensive set of wrapper functions
// to simplify working with Excel files using the excelize library.
package excelutil

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// --- File Operations ---

// CreateExcelFile creates a new in-memory Excel file with a default sheet.
func CreateExcelFile(sheetName string) (*excelize.File, error) {
	f := excelize.NewFile()
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	// Delete the default "Sheet1" that is created automatically.
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(0)
	return f, nil
}

// OpenExcel opens an existing Excel file from the given path.
func OpenExcel(path string) (*excelize.File, error) {
	return excelize.OpenFile(path)
}

// Save saves the Excel file to the specified path.
func Save(f *excelize.File, path string) error {
	return f.SaveAs(path)
}

// --- Cell Operations ---

// WriteCell writes a value to a specific cell in a sheet.
func WriteCell(f *excelize.File, sheet, cell string, value interface{}) error {
	return f.SetCellValue(sheet, cell, value)
}

// ReadCell reads the value from a specific cell in a sheet.
func ReadCell(f *excelize.File, sheet, cell string) (string, error) {
	return f.GetCellValue(sheet, cell)
}

// --- Row Operations ---

// WriteRow writes a slice of values to a specific row in a sheet, starting from column A.
func WriteRow(f *excelize.File, sheet string, row int, values []interface{}) error {
	startCell := fmt.Sprintf("A%d", row)
	return f.SetSheetRow(sheet, startCell, &values)
}

// ReadRow reads all cells in a specific row.
func ReadRow(f *excelize.File, sheet string, row int) ([]string, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if row-1 >= len(rows) {
		return nil, fmt.Errorf("row %d does not exist in sheet '%s'", row, sheet)
	}
	return rows[row-1], nil
}

// --- Column Operations ---

// WriteColumn writes a slice of values to a specific column in a sheet, starting from row 1.
func WriteColumn(f *excelize.File, sheet string, col string, values []interface{}) error {
	return f.SetSheetCol(sheet, col+"1", &values)
}

// ReadColumn reads all cells in a specific column.
func ReadColumn(f *excelize.File, sheet string, col string) ([]string, error) {
	colNum, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		return nil, err
	}
	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}
	if colNum-1 >= len(cols) {
		return nil, fmt.Errorf("column %s does not exist in sheet '%s'", col, sheet)
	}
	return cols[colNum-1], nil
}

// --- Sheet Operations ---

// AddSheet adds a new sheet to the Excel file.
func AddSheet(f *excelize.File, sheetName string) (int, error) {
	return f.NewSheet(sheetName)
}

// DeleteSheet removes a sheet from the Excel file.
func DeleteSheet(f *excelize.File, sheetName string) error {
	return f.DeleteSheet(sheetName)
}

// ListSheets returns a slice of sheet names in the Excel file.
func ListSheets(f *excelize.File) []string {
	return f.GetSheetList()
}

// --- Data Import/Export ---

// ExportStructsToExcel writes a slice of structs to a sheet.
// It uses reflection to create a header row from the struct field names.
func ExportStructsToExcel(f *excelize.File, sheet string, data interface{}) error {
	slice := reflect.ValueOf(data)
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice of structs")
	}
	if slice.Len() == 0 {
		return nil // Nothing to write
	}

	// Get headers from struct fields
	structType := slice.Index(0).Type()
	headers := make([]string, structType.NumField())
	for i := range headers {
		headers[i] = structType.Field(i).Name
	}
	SetHeaderRow(f, sheet, headers)

	// Write data rows
	for i := 0; i < slice.Len(); i++ {
		row := make([]interface{}, len(headers))
		structVal := slice.Index(i)
		for j := range headers {
			row[j] = structVal.Field(j).Interface()
		}
		WriteRow(f, sheet, i+2, row) // +2 because header is row 1
	}
	return nil
}

// ImportExcelToStructs reads an Excel sheet into a slice of structs.
// The target must be a pointer to a slice of structs.
func ImportExcelToStructs(f *excelize.File, sheet string, target interface{}) error {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target must be a pointer to a slice of structs")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return fmt.Errorf("sheet must have a header row and at least one data row")
	}

	headers := rows[0]
	sliceElemType := targetVal.Elem().Type().Elem()
	newSlice := reflect.MakeSlice(targetVal.Elem().Type(), 0, len(rows)-1)

	for i := 1; i < len(rows); i++ {
		newStruct := reflect.New(sliceElemType).Elem()
		for j, header := range headers {
			field := newStruct.FieldByName(header)
			if field.IsValid() && field.CanSet() && j < len(rows[i]) {
				// This is a simplified conversion, a real implementation might need type conversions
				field.SetString(rows[i][j])
			}
		}
		newSlice = reflect.Append(newSlice, newStruct)
	}

	targetVal.Elem().Set(newSlice)
	return nil
}

// --- Styling ---

// CreateStyle creates a new style from a definition and returns its ID.
// The styleDef map is marshaled to JSON and then unmarshaled into an excelize.Style struct.
func CreateStyle(f *excelize.File, styleDef map[string]interface{}) (int, error) {
	jsonDef, err := json.Marshal(styleDef)
	if err != nil {
		return 0, err
	}

	var style excelize.Style
	if err := json.Unmarshal(jsonDef, &style); err != nil {
		return 0, err
	}

	return f.NewStyle(&style)
}

// SetCellStyle applies a pre-defined style to a cell.
func SetCellStyle(f *excelize.File, sheet, cell string, styleID int) error {
	return f.SetCellStyle(sheet, cell, cell, styleID)
}

// --- Utility ---

// AutoFitColumns adjusts the width of all columns in a sheet to fit the content.
func AutoFitColumns(f *excelize.File, sheet string) error {
	cols, err := f.GetCols(sheet)
	if err != nil {
		return err
	}
	for i := range cols {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		if err := f.AutoFitCol(sheet, colName); err != nil {
			return err // Or log the error and continue
		}
	}
	return nil
}

// MergeCells merges a range of cells.
func MergeCells(f *excelize.File, sheet, startCell, endCell string) error {
	return f.MergeCell(sheet, startCell, endCell)
}

// SetHeaderRow writes a slice of strings to the first row of a sheet.
func SetHeaderRow(f *excelize.File, sheet string, headers []string) error {
	// Convert []string to []interface{}
	interfaceHeaders := make([]interface{}, len(headers))
	for i, v := range headers {
		interfaceHeaders[i] = v
	}
	return WriteRow(f, sheet, 1, interfaceHeaders)
}
