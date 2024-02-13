package main

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"myservice.com/transactions/internal/clients/filereader"
	"myservice.com/transactions/internal/clients/sender"
	"myservice.com/transactions/internal/config"
	"myservice.com/transactions/internal/service"
)

const (
	bucketName    = "transactions-summary-service"
	configFileKey = "transactions-summary-service.yaml"
)

func runService() {
	log.Println("transactions-summary-service has started...")
	// initialize configuration
	body, err := readFileFromS3(configFileKey)
	if err != nil {
		log.Panic("unable to get config file from s3:", err.Error())
	}

	conf, err := config.GetConfigFromS3(body)
	if err != nil {
		log.Panic("error getting config data from S3:", err.Error())
	}

	txnsBody, err := readFileFromS3(conf.TransactionsFile)
	if err != nil {
		log.Panic("unable to get transactions file from s3", err.Error())
	}

	// Create file reader and get data
	fileReader := filereader.New(txnsBody)

	// Create Sendgrid client
	emailSender := sender.New(conf.EmailSender)

	// Create database client

	// Create service with clients interfaces
	srv := service.New(&fileReader, emailSender)

	err = srv.Run()
	if err != nil {
		log.Printf("an error ocurred while running the service: %s", err.Error())
	}

	log.Println("finished...")
}

func main() {
	lambda.Start(runService)
}

func readFileFromS3(fileKey string) ([]byte, error) {
	// sdkConfig, err := awsconfig.LoadDefaultConfig(context.TODO())
	sdkConfig := aws.Config{
		Region: "us-east-1",
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		log.Println("couldn't get object:", err.Error())
		return []byte{}, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Println("unable to read object body:", err.Error())

		return []byte{}, err
	}

	return body, err

}
