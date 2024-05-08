package helpers

import (
	"fmt"
	"strconv"
)

func FixFloatPrecision(number float64) (float64, error)  {
	stringNumber := fmt.Sprintf("%.1f", number)
	numberParsed, err := strconv.ParseFloat(stringNumber, 2)
	if err != nil {
		return 0.0, err
	}
	return numberParsed, nil
}