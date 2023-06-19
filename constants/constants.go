package constants

var ExcelFunctions = map[string]bool{
	"concat":  true,
	"text":    true,
	"incFrom": true,
	"sum":     true,
	"spread":  true,
	"split":   true,
}

func IsExcelFunction(s string) bool {
	_, exists := ExcelFunctions[s]
	return exists
}
