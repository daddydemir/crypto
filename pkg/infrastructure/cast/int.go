package cast

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
)

func Int64(i interface{}) (int64, error) {

	switch v := i.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return 0, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		if v != math.Trunc(v) {
			return 0, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
	}
}

func Int64WithoutError(i interface{}) int64 {
	v, err := Int64(i)
	if err != nil {
		slog.Error("Int64WithoutError", "error", err)
	}
	return v
}
