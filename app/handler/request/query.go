package request

import (
	"fmt"
	"net/http"
	"strconv"
)

type Option struct {
	Name         string
	DefaultValue int64
	MinValue     int64
	MaxValue     int64
}

func GetOptionParams(r *http.Request, options []Option) (map[string]int64, error) {
	params := make(map[string]int64)
	var err error

	for _, op := range options {
		params[op.Name], err = paramOf(r.URL.Query().Get(string(op.Name)), op.DefaultValue, op.MinValue, op.MaxValue)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}

func paramOf(strParam string, defaultValue, min, max int64) (int64, error) {
	if strParam == "" {
		return defaultValue, nil
	}

	param, err := strconv.ParseInt(strParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%q is invalid format for option", strParam)
	}
	if param < min || max < param {
		return 0, fmt.Errorf("%d is over valid range [%d, %d]", param, min, max)
	}

	return param, nil
}
