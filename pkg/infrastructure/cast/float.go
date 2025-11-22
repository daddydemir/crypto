package cast

import (
	"fmt"
	"log/slog"
	"strconv"
)

func Float64(i interface{}) (float64, error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot cast %v to float64, type: %T", v, i)
	}
}

func Float64WithoutError(i interface{}) float64 {
	v, err := Float64(i)
	if err != nil {
		slog.Error("Float64WithoutError", "error", err)
	}
	return v
}
