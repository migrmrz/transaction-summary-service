package service

import "myservice.com/transactions/internal/clients/filereader"

type mockFileReader struct{}

func (m mockFileReader) GetTransactions() ([]filereader.Transaction, error) {
	return []filereader.Transaction{
		{
			Id:   "1",
			Date: "4/20",
			Txn:  "+34.33",
		},
		{
			Id:   "2",
			Date: "4/21",
			Txn:  "-14.53",
		},
		{
			Id:   "3",
			Date: "4/22",
			Txn:  "+4.3",
		},
	}, nil
}
