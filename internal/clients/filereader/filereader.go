package filereader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Transaction struct {
	Id   string
	Date string
	Txn  string
}

//type ReaderHandler interface {
//	GetTransactions() ([]Transaction, error)
//}

type Reader struct {
	fileName string
}

// New creates a new Reader client
func New(fileName string) Reader {
	return Reader{
		fileName: fileName,
	}
}

// GetTransactions opens and reads file from fileName to return data
func (r Reader) GetTransactions() (data []Transaction, err error) {
	// Reads file
	file, err := os.Open(r.fileName)
	if err != nil {
		return data, fmt.Errorf("an error ocurred while opening the file: %s\n", err.Error())
	}

	defer file.Close()

	lineCounter := 0

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return data, fmt.Errorf("error: %s\n", err.Error())
			}

			break
		}

		if lineCounter == 0 { // exclude header
			lineCounter += 1

			continue
		}

		txn := Transaction{
			Id:   record[0],
			Date: record[1],
			Txn:  record[2],
		}

		data = append(data, txn)

	}

	return data, nil

}
