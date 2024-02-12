package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"myservice.com/transactions/internal/clients/filereader"
	"myservice.com/transactions/internal/clients/sender"
	"myservice.com/transactions/internal/config"
	"myservice.com/transactions/internal/service"
)

func runService() {
	log.Println("transactions-summary-service has started...")
	// initialize configuration
	// configPath := flag.String("conf", "s3://transactions-summary-service/etc/transactions-summary-service", "directory where config file is located")
	err := downloadFile()
	if err != nil {
		log.Panic("unable to get file from s3:", err.Error())
	}

	configPath := "/"

	conf, err := config.GetConfig(configPath)
	if err != nil {
		log.Printf("unable to read config file: %s", err.Error())
	}

	// Create file reader and get data
	fileReader := filereader.New(conf)

	// Create Sendgrid client
	emailSender := sender.New(conf.EmailSender)

	// Create database client

	// Create service with clients interfaces
	srv := service.New(fileReader, emailSender)

	err = srv.Run()
	if err != nil {
		log.Printf("an error ocurred while running the service: %s", err.Error())
	}

	log.Println("finished...")
}

func main() {
	lambda.Start(runService)
}

func downloadFile() error {
	// sdkConfig, err := awsconfig.LoadDefaultConfig(context.TODO())
	sdkConfig := aws.Config{
		Region: "us-east-1",
	}
	//if err != nil {
	//		return fmt.Errorf("error loading aws config: %s", err)
	//	}

	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("transactions-summary-service"),
		Key:    aws.String("etc/transactions-summary-service/transactions-summary-service.yaml"),
	})
	if err != nil {
		log.Println("couldn't get object:", err.Error())
		return err
	}

	defer result.Body.Close()

	file, err := os.Create("transactions-summary-service.yaml")
	if err != nil {
		log.Println("couldn't create file:", err.Error())
		return err
	}

	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Println("unable to read object body:", err.Error())
	}

	_, err = file.Write(body)

	return err

}
