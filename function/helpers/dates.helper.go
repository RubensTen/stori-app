package helpers

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func GetMonthName(month string) (string, error) {
	mid, err := strconv.Atoi(strings.TrimSpace(month))
	if err != nil {
		return "", errors.New("invalid month")
	}
	m := time.Month(mid)
	return m.String(), nil
}
