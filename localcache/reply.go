package localcache

import (
	"strconv"
	"errors"
	"fmt"
	"strings"
)

var ErrNil = errors.New("localcache: nil returned")

var ErrType = errors.New("localcache: type err")

var ErrNotFound = errors.New("localcache: not found")

func GetInt64(key string) (int64, error) {
	if reply, err := Get(key); err == nil {
		switch reply := reply.(type) {
		case int64:
			return reply, nil
		case string:
			n, err := strconv.ParseInt(reply, 10, 64)
			return n, err
		case nil:
			return 0, nil
		default:
			return 0, ErrType
		}
	} else {
		return 0, err
	}
}

func GetFloat64(key string) (float64, error) {
	if reply, err := Get(key); err == nil {
		switch reply := reply.(type) {
		case float64:
			return reply, nil
		case string:
			n, err := strconv.ParseFloat(reply, 64)
			return n, err
		case nil:
			return 0, nil
		case int:
			return float64(reply), nil
		case int8:
			return float64(reply), nil
		case int32:
			return float64(reply), nil
		case int64:
			return float64(reply), nil
		case uint:
			return float64(reply), nil
		case uint8:
			return float64(reply), nil
		case uint32:
			return float64(reply), nil
		case uint64:
			return float64(reply), nil
		default:
			return 0, ErrType
		}

	} else {
		return 0, err
	}
}

func GetInt(key string) (int, error) {
	if reply, err := Get(key); err == nil {
		switch reply := reply.(type) {
		case int:
			return reply, nil
		case string:
			n, err := strconv.Atoi(reply)
			return n, err
		case nil:
			return 0, nil
		default:
			return 0, ErrType
		}
	} else {
		return 0, err
	}
}


func GetString(key string) (string, error) {
	if reply, err := Get(key); err == nil {
		return fmt.Sprint(reply), nil
	} else {
		return "", err
	}
}

func GetStrings(key string) ([]string, error) {
	if reply, err := Get(key); err == nil {
		switch reply := reply.(type) {
		case []string:
			return reply, nil
		case nil:
			return []string{}, nil
		default:
			return nil, ErrType
		}
	} else {
		return nil, err
	}
	return nil, ErrNotFound
}

func isNumeric(v string) bool {
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return true
	}
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
		if _, err := strconv.ParseInt(v, 0, 64); err == nil {
			return true
		}
	}
	return false
}
