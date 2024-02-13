package filereader

import (
	"log"

	"github.com/gocarina/gocsv"
)

type Transaction struct {
	Id   string `csv:"Id"`
	Date string `csv:"Date"`
	Txn  string `csv:"Transaction"`
}

type Reader struct {
	body []byte
}

// New creates a new Reader client
func New(body []byte) Reader {
	return Reader{
		body: body,
	}
}

// GetTransactions opens and reads file from fileName to return data
func (r Reader) GetTransactions() (data []Transaction, err error) {
	log.Println("reading transactions file and getting data...")

	err = gocsv.UnmarshalBytes(r.body, &data)
	if err != nil {
		return data, err
	}

	return data, nil

}
