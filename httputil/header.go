package httputil

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetHeaderInt64(r *http.Request, key string) (int64, error) {
	value := r.Header.Get(key)
	if value == "" {
		return 0, fmt.Errorf("header %s is required", key)
	}
	return strconv.ParseInt(value, 10, 64)
}

func GetHeaderString(r *http.Request, key string) (string, error) {
	value := r.Header.Get(key)
	if value == "" {
		return "", fmt.Errorf("header %s is required", key)
	}
	return value, nil
}
