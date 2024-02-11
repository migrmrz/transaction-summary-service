package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"myservice.com/transactions/internal/clients/filereader"
	"myservice.com/transactions/internal/clients/sender"
	"myservice.com/transactions/internal/config"
	"myservice.com/transactions/internal/service"
)

func main() {
	// initialize configuration
	configPath := flag.String("conf", "/etc/transactions-summary-service", "directory where config file is located")

	conf, err := config.GetConfig(*configPath)
	if err != nil {
		log.Printf("unable to read config file: %s", err.Error())
	}

	log.Println("validating args...")

	// validating and gettings args
	if len(os.Args) == 1 {
		log.Println("no file name and email provided. Exiting...")
		os.Exit(1)
	}

	// Get file name and email from args
	fileName := os.Args[1]
	email := os.Args[2]

	// Create file reader and get data
	fileReader := filereader.New(fileName)

	// Create Sendgrid client
	emailSender := sender.New(conf.EmailSender, email)

	// Create database client
	fmt.Println("This part will create database client")

	// Create service with clients interfaces
	srv := service.New(fileReader, emailSender)

	err = srv.Run()
	if err != nil {
		log.Printf("an error ocurred while running the service: %s", err.Error())
	}

	log.Println("Finished...")
}
