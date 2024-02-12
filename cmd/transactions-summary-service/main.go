package main

import (
	"flag"
	"log"

	"myservice.com/transactions/internal/clients/filereader"
	"myservice.com/transactions/internal/clients/sender"
	"myservice.com/transactions/internal/config"
	"myservice.com/transactions/internal/service"
)

func main() {
	log.Println("transactions-summary-service has started...")
	// initialize configuration
	configPath := flag.String("conf", "/etc/transactions-summary-service", "directory where config file is located")

	conf, err := config.GetConfig(*configPath)
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
