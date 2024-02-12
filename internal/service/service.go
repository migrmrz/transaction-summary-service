package service

import (
	"encoding/json"
	"fmt"
	"log"

	"myservice.com/transactions/internal/clients/filereader"
	"myservice.com/transactions/internal/clients/sender"
)

type Money int64

func (m Money) Dollars() int64 { return int64(m) / 100 }
func (m Money) Cents() int64   { return int64(m) % 100 }

type fileReader interface {
	GetTransactions() ([]filereader.Transaction, error)
}

type emailSender interface {
	SendEmail(email sender.Email) error
}

type Service struct {
	fileReader  fileReader
	emailSender emailSender
}

// New creates a new Service
func New(fr fileReader, es emailSender) Service {
	return Service{
		fileReader:  fr,
		emailSender: es,
	}
}

func (s *Service) Run() error {
	// Opening and reading file to get data
	data, err := s.fileReader.GetTransactions()
	if err != nil {
		return fmt.Errorf("unable to read file: %s", err.Error())
	}

	log.Println("raw data from file:", data)

	// Get debit and credit averages
	averageDebit, averageCredit := getAverages(data)
	averageDebitStr := fmt.Sprintf("%d.%d", Money(averageDebit).Dollars(), Money(averageDebit).Cents()*-1)
	averageCreditStr := fmt.Sprintf("%d.%d", Money(averageCredit).Dollars(), Money(averageCredit).Cents())

	log.Printf(
		"averages... debit: %d.%d, credit: %d.%d\n",
		Money(averageDebit).Dollars(), Money(averageDebit).Cents()*-1,
		Money(averageCredit).Dollars(), Money(averageCredit).Cents(),
	)

	// Get total balance
	totalBalance := getTotalBalance(data)
	var totalBalanceStr string

	if totalBalance < 0 {
		log.Printf("total balance: %d.%d\n", Money(totalBalance).Dollars(), Money(totalBalance).Cents()*-1)

		totalBalanceStr = fmt.Sprintf("%d.%d", Money(totalBalance).Dollars(), Money(totalBalance).Cents()*-1)

	} else {
		log.Printf("total balance: %d.%d", Money(totalBalance).Dollars(), Money(totalBalance).Cents())

		totalBalanceStr = fmt.Sprintf("%d.%d", Money(totalBalance).Dollars(), Money(totalBalance).Cents())
	}

	// Get transaction by month
	transactions := getTransactionsByMonth(data)
	bs, _ := json.Marshal(transactions)
	log.Printf("transaction count by month: %s", string(bs))

	// Send email
	email := sender.Email{
		TotalBalance:             totalBalanceStr,
		TotalTransactionsByMonth: transactions,
		AverageCreditAmount:      averageCreditStr,
		AverageDebitAmount:       averageDebitStr,
	}

	err = s.emailSender.SendEmail(email)
	if err != nil {
		log.Printf("something happened while trying to send email: %s", err.Error())
	}

	return nil
}
