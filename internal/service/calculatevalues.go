package service

import (
	"strconv"
	"strings"
	"time"

	"my-service.com/transactions/internal/clients/filereader"
)

// GetAverages calculates the credit and debit averages from transactions
func (s *Service) GetAverages(transactions []filereader.Transaction) (averageDebit, averageCredit int64) {
	var debitTotal int64
	var creditTotal int64

	var creditCounter int
	var debitCounter int

	// Get the credit and debit values from each transaction
	for _, transaction := range transactions {
		// convert to money (int64) type
		money := convertTransactionValue(transaction.Txn)
		if money < 0 { // this is debit
			debitTotal += int64(money)
			debitCounter += 1
		} else { // this is credit
			creditTotal += int64(money)
			creditCounter += 1
		}
	}

	return debitTotal / int64(debitCounter), creditTotal / int64(creditCounter)
}

// GetTotalBalance calculates the total amount from transactions
func (s *Service) GetTotalBalance(transactions []filereader.Transaction) (totalBalance int64) {
	for _, transaction := range transactions {
		money := convertTransactionValue(transaction.Txn)
		totalBalance += int64(money)
	}

	return totalBalance
}

// GetTransactionsByMonth gets the number of transactions by month from records
func (s *Service) GetTransactionsByMonth(transactions []filereader.Transaction) map[string]int {
	monthsCount := make(map[string]int)

	for _, transaction := range transactions {
		month, _ := strconv.Atoi(strings.Split(transaction.Date, "/")[0])
		monthsCount[time.Month(month).String()] += 1
	}

	return monthsCount
}

// convertTransactionValue takes in the transaction value as string to convert it to int64 as cents
func convertTransactionValue(stringValue string) Money {
	var value int64
	var cents int64

	splitValues := strings.Split(stringValue, ".") // split dollars and cents with separator
	value, _ = strconv.ParseInt(splitValues[0], 10, 64)
	if len(splitValues) > 1 {
		cents, _ = strconv.ParseInt(splitValues[1], 10, 64)
		if cents < 10 {
			cents = cents * 10
		}
	}

	if stringValue[0] == '-' { // debit value
		return Money((value * 100) - cents)
	}

	// credit value
	return Money((value * 100) + cents)
}
