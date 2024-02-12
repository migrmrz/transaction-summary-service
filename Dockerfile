FROM golang:1.20.12-alpine3.19

WORKDIR /app/transactions-summary-service

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd

COPY internal/config/transactions-summary-service.yaml /etc/transactions-summary-service/

COPY txns/ /mnt/transactions/

COPY internal/ ./internal

COPY scripts/ ./scripts

COPY txns/ ./txns

RUN go build -v -o ./build/transactions-summary-service ./cmd/transactions-summary-service/main.go

CMD ["./build/transactions-summary-service"]
