package cast

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
)

func intToInt64(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int16:
		return int64(v), true
	case int8:
		return int64(v), true
	}
	return 0, false
}

func uintToInt64(i interface{}) (int64, bool, error) {
	switch v := i.(type) {
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0, false, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), true, nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, false, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), true, nil
	case uint32:
		return int64(v), true, nil
	case uint16:
		return int64(v), true, nil
	case uint8:
		return int64(v), true, nil
	}
	return 0, false, nil
}

func floatToInt64(i interface{}) (int64, bool, error) {
	switch v := i.(type) {
	case float64:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return 0, false, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		if v != math.Trunc(v) {
			return 0, false, fmt.Errorf("cannot cast %v to int64, type: %T", v, i)
		}
		return int64(v), true, nil
	case float32:
		return int64(v), true, nil
	}
	return 0, false, nil
}

func Int64(i interface{}) (int64, error) {
	if v, ok := intToInt64(i); ok {
		return v, nil
	}
	if v, ok, err := uintToInt64(i); ok {
		return v, err
	} else if err != nil {
		return 0, err
	}
	if v, ok, err := floatToInt64(i); ok {
		return v, err
	} else if err != nil {
		return 0, err
	}
	switch v := i.(type) {
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
