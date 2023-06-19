package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

func GetMaxRowsAndCols(cellMap map[string]string) (int, int) {
	maxRows, maxCols := 0, 0
	for key := range cellMap {
		col := key[:1]
		row, _ := strconv.Atoi(key[1:])

		// If this row is larger than the current max, update maxRows
		if row > maxRows {
			maxRows = row
		}

		// Convert the column from a letter to a number, assuming 'A' is 1, 'B' is 2, etc.
		// If this column is larger than the current max, update maxCols
		colNum := int(col[0] - 'A' + 1)
		if colNum > maxCols {
			maxCols = colNum
		}
	}

	return maxRows, maxCols
}

// check if the cell contains standalone formula - Formulas that do not require evaluation of other cells
func isFormula(value string) bool {
	// Check if it starts with '='
	if !strings.HasPrefix(value, "=") {
		return false
	}

	// Check for DoubleCaret, we leave those for later
	if strings.HasPrefix(value, "=^^") {
		return false
	}

	return true
}

// check if the cell contains standalone formula - Formulas that do not require evaluation of other cells
func isStandaloneFormula(value string) bool {
	// Check if it starts with '='
	if !strings.HasPrefix(value, "=") {
		return false
	}
	// Check for DoubleCaret, we leave those for later
	if strings.HasPrefix(value, "=^^") {
		return false
	}

	// Check if it contains A^, B^, etc. or A^v, B^v, etc.
	matches, _ := regexp.MatchString(`[A-Z]\^v?`, value)
	if matches {
		return false
	}

	return true
}
func GetStandaloneFormulas(cellMap map[string]string) map[string]string {
	standaloneFormulas := make(map[string]string)

	// Check if it contains A^, B^, etc. or A^v, B^v, etc.
	for key, value := range cellMap {
		if isStandaloneFormula(value) {
			standaloneFormulas[key] = value
		}
	}

	return standaloneFormulas
}

// extract standalone formulas from cell map
func GetFormulas(cellMap map[string]string) map[string]string {
	formulas := make(map[string]string)

	for key, value := range cellMap {
		if isFormula(value) {
			formulas[key] = value
		}
	}

	return formulas
}

// This only replaces a cell value with its pointer in a formula
func MapCellsInFormulas(cellMap map[string]string) map[string]string {
	cellRefFormulas := make(map[string]string)
	re := regexp.MustCompile(`[A-Z]\d+`)

	for cell, value := range cellMap {
		if re.MatchString(value) {
			cellRefFormulas[cell] = value
		}
	}

	return cellRefFormulas
}

func MapFormulasToCellMap(cellMap map[string]string, standaloneFormulas map[string]string) {
	for cell, formula := range standaloneFormulas {
		cellMap[cell] = formula
	}
}

func CleanFormula(formula string) string {
	formula = strings.ReplaceAll(formula, `\`, "")
	formula = strings.ReplaceAll(formula, `"`, "")
	formula = strings.Replace(formula, "=", "", 1)
	return formula
}
