package utils

import (
	"fmt"
	"strings"
	"time"
)

func MakeValue(value interface{}) (string, bool) {
	switch value.(type) {
	case float32:
		result := fmt.Sprintf("%f", value)
		return result[:len(result)-2], true
	case float64:
		return fmt.Sprintf("%f", value), true
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", value), true
	case bool:
		if value == true {
			return "1", true
		}
		return "0", true
	case string:
		if result, ok := value.(string); ok {
			return fmt.Sprintf("'%s'", strings.ReplaceAll(result, "'", "''")), true
		}
	case time.Time:
		if result, ok := value.(time.Time); ok {
			return fmt.Sprintf("'%s'", TimeToSQL(result)), true
		}
	}
	return "", false
}
