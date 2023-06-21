package functions

import (
	"fmt"
	"kokal/helpers"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

func ProcessDoubleCaret(cellMap map[string]string) {
	changed := true
	for changed {
		changed = false
		for key, value := range cellMap {
			if value == "=^^" {
				column := key[:1]
				row, _ := strconv.Atoi(key[1:])
				aboveKey := fmt.Sprintf("%s%d", column, row-1)
				if aboveVal, ok := cellMap[aboveKey]; ok {
					cellMap[key] = aboveVal
					changed = true
				}
			}
		}
	}
}

func SimplifyFormulas(cellMap map[string]string, standaloneFormulas map[string]string) map[string]string {
	processedFormulas := make(map[string]string)
	// Regular expression to match cell references
	cellRe := regexp.MustCompile(`([A-Z]\d+)`)
	// Regular expression to match headers with index
	headerRe := regexp.MustCompile(`@(.*?)<1>`)

	for cell, formula := range standaloneFormulas {
		// Process header references
		headerMatches := headerRe.FindAllStringSubmatch(formula, -1)
		for _, match := range headerMatches {
			// match[0] is the full string (e.g., "@adjusted_cost<1>"), match[1] is the header (e.g., "adjusted_cost")
			headerCell := findCellForHeader(cellMap, match[1])
			value, exists := cellMap[headerCell]
			if exists {
				// Replace the header reference with the value from the cell
				formula = strings.Replace(formula, match[0], value, -1)
			}
		}

		// Process cell references
		cellMatches := cellRe.FindAllString(formula, -1)
		for _, match := range cellMatches {
			value, exists := cellMap[match]
			if exists {
				formula = strings.Replace(formula, match, value, -1)
			} else {
				// If the value is not found or empty string, we write 0
				formula = strings.Replace(formula, match, "0", -1)
			}
		}

		processedFormulas[cell] = formula
	}

	return processedFormulas
}

func findCellForHeader(cellMap map[string]string, header string) string {
	for cell, value := range cellMap {
		// If the cell value starts with "!", it's a header
		if strings.HasPrefix(value, "!") && strings.TrimLeft(value, "!") == header {
			// Split the cell into column and row parts
			column := cell[:1]
			row, err := strconv.Atoi(cell[1:])
			if err != nil {
				continue // If the row part is not a number, skip this cell
			}
			// The cell for the header is in the next row
			return fmt.Sprintf("%s%d", column, row+1)
		}
	}
	return "" // Return an empty string if the header was not found
}

// SOLVING
func getInnerMostFunction(formula string) (string, bool) {
	re := regexp.MustCompile(`[a-zA-Z]+\([^()]*\)`)

	matches := re.FindAllString(formula, -1)

	// If there are no matches, return the input string and false.
	if matches == nil {
		return formula, false
	}

	// Return the last match (which will be the innermost function) and true.
	return matches[len(matches)-1], true
}

func ProcessFormula(key string, toProcess string, standaloneFormulas *map[string]string) string {
	formula := helpers.CleanFormula(toProcess)
	innerFunction, found := getInnerMostFunction(formula)

	for found {
		result, err := solveFunction(innerFunction)
		if err != nil {
			continue
		}
		formula = strings.Replace(formula, innerFunction, result, -1)
		innerFunction, found = getInnerMostFunction(formula)
	}

	(*standaloneFormulas)[key] = formula
	return formula
}

func solveFunction(fn string) (string, error) {
	// Extract function name and parameters string
	nameRegex := regexp.MustCompile(`^([a-zA-Z]+)\((.*)\)$`)
	match := nameRegex.FindStringSubmatch(fn)
	if len(match) != 3 {
		return "", fmt.Errorf("invalid function format")
	}
	funcName := match[1]
	paramsStr := match[2]

	// Split parameters string into a slice
	params := strings.Split(paramsStr, ",")

	// Check if the function name exists in the constant map
	_, found := funcMap[funcName]
	if !found {
		return "", fmt.Errorf("unknown function: %s", funcName)
	}

	// Call the function with the parameters
	result, err := callFunction(funcName, params)
	if err != nil {
		return "", err
	}

	return result, nil
}

func SolveExpression(key string, expressionString string, expressions *map[string]string) float64 {
	expression, err := govaluate.NewEvaluableExpression(expressionString)
	if err != nil {
		fmt.Println("Error creating expression:", err)
		return 0
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Println("Error evaluating expression:", err)
		return 0
	}

	(*expressions)[key] = fmt.Sprintf("%f", result.(float64))
	return result.(float64)
}

func callFunction(funcName string, params []string) (string, error) {
	switch funcName {
	case "bte":
		if len(params) != 2 {
			return "", fmt.Errorf("bte expects exactly 2 arguments")
		}
		a, err1 := strconv.ParseFloat(params[0], 64)
		b, err2 := strconv.ParseFloat(params[1], 64)
		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("bte expects arguments that can be parsed to float64")
		}
		return strconv.FormatBool(bte(a, b)), nil
	case "incFrom":
		paramInt, err := strconv.Atoi(params[0])
		if err != nil {
			return "", fmt.Errorf("invalid parameter: %s", params[0])
		}
		return strconv.Itoa(incFrom(paramInt)), nil
	case "text":
		paramInt, err := strconv.Atoi(params[0])
		if err != nil {
			return "", fmt.Errorf("invalid parameter: %s", params[0])
		}
		return text(paramInt), nil
	case "concat":
		return concat(params...), nil
	case "split":
		// The params are already splitted from the solve function, we just return them here
		filteredParams := helpers.FilterEmptyStrings(params)
		return strings.Join(filteredParams, ","), nil
	case "spread":
		spreadResult, err := spread(params)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", spreadResult), nil
	case "sum":
		nums := make([]float64, len(params))
		for i, param := range params {
			num, err := strconv.ParseFloat(param, 64)
			if err != nil {
				return "", fmt.Errorf("invalid parameter: %s", param)
			}
			nums[i] = num
		}
		return strconv.FormatFloat(sum(nums), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported function: %s", funcName)
	}
}

// EXCEL FUNCTIONS

func incFrom(n int) int {
	return n + 1
}

func text(n int) string {
	return strconv.Itoa(n)
}

func concat(strs ...string) string {
	return strings.Join(strs, "")
}

func split(s, delimiter string) []string {
	return strings.Split(s, delimiter)
}

func spread(params []string) (float64, error) {
	var numbers []float64

	for _, param := range params {
		num, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return 0, err
		}
		numbers = append(numbers, num)
	}

	min, max := numbers[0], numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	return max - min, nil
}

func sum(nums []float64) float64 {
	total := 0.0
	for _, num := range nums {
		total += num
	}
	return total
}

func bte(a float64, b float64) bool {
	return a >= b
}

var funcMap = map[string]interface{}{
	"incFrom": incFrom,
	"text":    text,
	"concat":  concat,
	"split":   split,
	"spread":  spread,
	"sum":     sum,
	"bte":     bte,
}

func GetFunction(name string) (interface{}, bool) {
	function, found := funcMap[name]
	return function, found
}
