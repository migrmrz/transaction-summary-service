# transactions-summary-service

Simple service that reads a transactions file in a csv format, processes it and sends a report email with the data.
When processing, the following data will be collected to be sent:
- Total balance
- Number of transactions by month
- Average credit amount
- Average debit amount

## How to run the code

### Locally

The most useful make targets for working locally are:

* `make build`: Builds the service.
* `make run FILE=/path/to/file EMAIL=email@example.com`: Starts the service. FILE is the path where the file to be processed is and EMAIL is the user's email to which the report will be sent.
* `make clean`: Clean temporary files.

Configuration file located in `internal/config/transactions-summary.yaml` has a key to turn on/off the sending of the email action so this can be switched back and forth depending on what the user wants to happen.
```yaml
send-email: false
```

### Docker

### Example

`make run FILE=txns/migrmrz/txns.csv EMAIL=migrmrz@gmail.com`

## Dependencies

This project has dependencies on:
* go (`1.20.12`)
