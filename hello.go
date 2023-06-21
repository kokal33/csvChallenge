package main

import (
	"fmt"
	"io/ioutil"
	"kokal/functions"
	"kokal/helpers"
	"log"
	"strings"
)

func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// replace all \r with an empty string
	dataStr := strings.ReplaceAll(string(data), "\r", "")
	return dataStr
}

func createCellMap(data string) map[string]string {
	rows := strings.Split(data, "\n")

	// creating a map to hold the cell mappings
	cellMap := make(map[string]string)

	// calculate the maximum number of columns
	maxCols := 0
	for _, row := range rows {
		cols := strings.Split(row, "|")
		if len(cols) > maxCols {
			maxCols = len(cols)
		}
	}

	// iterate through the rows
	rownum := 1
	firstHeader := true
	for _, row := range rows {
		// If a header starts and it's not the first one, give it 3 lines of space
		if strings.HasPrefix(row, "!") {
			if firstHeader {
				firstHeader = false
			} else {
				rownum += 3
			}
		}
		// split the row into columns
		cols := strings.Split(row, "|")
		// iterate through the columns
		for j, col := range cols {
			// map each cell to its respective data
			cellMap[fmt.Sprintf("%c%d", 'A'+j, rownum)] = col
		}
		rownum++
	}

	return cellMap
}

func printData(cellMap map[string]string, maxRows int, maxCols int) {
	// print the column headers dynamically
	for i := 0; i < maxCols; i++ {
		fmt.Printf("                %c", 'A'+i)
	}
	fmt.Println()

	for i := 1; i <= maxRows; i++ {
		fmt.Printf("%d", i)
		for j := 0; j < maxCols; j++ {
			cellKey := fmt.Sprintf("%c%d", 'A'+j, i)
			cellVal := cellMap[cellKey]
			fmt.Printf("        %s     |", cellVal)
		}
		fmt.Println()
	}
}

func main() {
	data := readFile("transactions.csv")
	// Create cell map and print it
	cellMap := createCellMap(data)
	maxRows, maxCols := helpers.GetMaxRowsAndCols(cellMap)

	// Get the all formulas
	allFormulas := helpers.GetAllFormulas(cellMap)
	// Process the formulas, mapping the cell numbers with their value
	processedAllFormulas := functions.SimplifyFormulas(cellMap, allFormulas)
	helpers.MapFormulasToCellMap(cellMap, processedAllFormulas)
	// With those processed, we can apply the doubleCaret formulas
	functions.ProcessDoubleCaret(cellMap)

	// Get standalone formulas, to solve them before solving the Evaluated ^ ones
	fmt.Println("****** SOLVING STANDALONE FORMULAS ******")
	standaloneFormulas := helpers.GetStandaloneFormulas(cellMap)
	for key, formula := range standaloneFormulas {
		fmt.Println("Start processing formula: ", formula, " on cell: ", key)
		processed := functions.ProcessFormula(key, formula, &standaloneFormulas)
		fmt.Println("Processed formula, result: ", processed)
	}
	helpers.MapFormulasToCellMap(cellMap, standaloneFormulas)

	// Get the rest of the formulas
	fmt.Println("****** SOLVING REST OF THE FORMULAS ******")
	lastFormulas := helpers.GetAllFormulas(cellMap)
	processedLastFormulas := helpers.MapEvaluatedCellsToFormula(maxRows, cellMap, lastFormulas)
	simplifiedLastFormulas := functions.SimplifyFormulas(cellMap, processedLastFormulas)
	for key, formula := range simplifiedLastFormulas {
		fmt.Println("Start processing formula: ", formula, " on cell: ", key)
		processed := functions.ProcessFormula(key, formula, &simplifiedLastFormulas)
		fmt.Println("Processed formula, result: ", processed)
	}

	// Solve the final expressions
	fmt.Println("****** SOLVING FINAL EXPRESSIONS ******")
	for key, formula := range simplifiedLastFormulas {
		fmt.Println("Start processing expression: ", formula, " on cell: ", key)
		solvedExpressions := functions.SolveExpression(key, formula, &simplifiedLastFormulas)
		fmt.Println("Solved final expression, result: ", solvedExpressions)
	}

	helpers.MapFormulasToCellMap(cellMap, simplifiedLastFormulas)
	printData(cellMap, maxRows, maxCols)
}
