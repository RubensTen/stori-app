package helpers

import (
	"strconv"
)

type Transactions struct {
	Email       string
	Name        string
	Id          string
	Date        string
	Transaction float64
}

func MapToTransactions(data [][]string) []Transactions {
	var tinfo []Transactions
	for _, row := range data {
		transactionInfo := &Transactions{
			Email: row[0],
			Name:  row[1],
			Id:    row[2],
			Date:  row[3],
		}
		t, err := strconv.ParseFloat(row[4], 2)
		if err != nil {
			// skip record, transaction invalid
			continue
		}
		transactionInfo.Transaction = t
		tinfo = append(tinfo, *transactionInfo)
	}
	return tinfo
}
