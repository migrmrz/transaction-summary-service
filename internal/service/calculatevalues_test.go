package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"myservice.com/transactions/internal/clients/filereader"
)

func TestConvertTransactionValue(t *testing.T) {
	cases := []struct {
		name           string
		stringValue    string
		expectedResult Money
	}{
		{
			name:           "success-debit-with-cents",
			stringValue:    "-455.23",
			expectedResult: Money(-45523),
		},
		{
			name:           "success-credit-with-cents",
			stringValue:    "455.23",
			expectedResult: Money(45523),
		},
		{
			name:           "success-credit-with-cents",
			stringValue:    "455.2",
			expectedResult: Money(45520),
		},
		{
			name:           "success-credit-with-no-cents",
			stringValue:    "455",
			expectedResult: Money(45500),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := convertTransactionValue(tc.stringValue)
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestGetAverages(t *testing.T) {
	cases := []struct {
		name                  string
		transactions          []filereader.Transaction
		expectedAverageCredit int64
		expectedAverageDebit  int64
	}{
		{
			name: "success",
			transactions: []filereader.Transaction{
				{
					Id:   "1",
					Date: "3/15",
					Txn:  "-456.33",
				},
				{
					Id:   "2",
					Date: "3/15",
					Txn:  "2432.53",
				},
				{
					Id:   "3",
					Date: "3/15",
					Txn:  "-5562.32",
				},
			},
			expectedAverageCredit: 243253,
			expectedAverageDebit:  -300932,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualAverageDebit, actualAverageCredit := getAverages(tc.transactions)
			assert.Equal(t, tc.expectedAverageDebit, actualAverageDebit)
			assert.Equal(t, tc.expectedAverageCredit, actualAverageCredit)
		})
	}
}

func TestGetTotalBalance(t *testing.T) {
	cases := []struct {
		name           string
		transactions   []filereader.Transaction
		expectedResult int64
	}{
		{
			name: "success-debit",
			transactions: []filereader.Transaction{
				{
					Id:   "1",
					Date: "3/15",
					Txn:  "-456.33",
				},
				{
					Id:   "2",
					Date: "3/15",
					Txn:  "2432.53",
				},
				{
					Id:   "3",
					Date: "3/15",
					Txn:  "-5562.32",
				},
			},
			expectedResult: -358612,
		},
		{
			name: "success-credit",
			transactions: []filereader.Transaction{
				{
					Id:   "1",
					Date: "3/15",
					Txn:  "-456.33",
				},
				{
					Id:   "2",
					Date: "3/15",
					Txn:  "2432.53",
				},
				{
					Id:   "3",
					Date: "3/15",
					Txn:  "-552.32",
				},
			},
			expectedResult: 142388,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := getTotalBalance(tc.transactions)
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestGetTransactionsByMonth(t *testing.T) {
	cases := []struct {
		name           string
		transactions   []filereader.Transaction
		expectedResult map[string]int
	}{
		{
			name: "success",
			transactions: []filereader.Transaction{
				{
					Id:   "1",
					Date: "1/15",
					Txn:  "-456.33",
				},
				{
					Id:   "2",
					Date: "2/15",
					Txn:  "2432.53",
				},
				{
					Id:   "3",
					Date: "3/15",
					Txn:  "-552.32",
				},
			},
			expectedResult: map[string]int{
				"January":  1,
				"February": 1,
				"March":    1,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := getTransactionsByMonth(tc.transactions)
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}
