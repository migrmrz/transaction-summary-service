package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"my-service.com/transactions/internal/clients/filereader"
	"my-service.com/transactions/internal/config"
	"my-service.com/transactions/internal/service"
)

func main() {
	fmt.Println("This it the main code")

	// initialize configuration
	configPath := flag.String("conf", "/etc/transactions-summary-service", "directory where config file is located")
	flag.Parse()

	_, err := config.GetConfig(*configPath)
	if err != nil {
		log.Println("unable to read config file")
	}

	log.Println("validating args...")

	// validating and gettings args
	if len(os.Args) == 1 {
		log.Println("no file name and email provided. Exiting...")
		os.Exit(1)
	}

	// Get file name and email from args
	fileName := os.Args[1]
	// email := os.Args[2]
	_ = os.Args[2]

	// Create file reader and get data
	fileReader := filereader.New(fileName)

	// Create Sendgrid client
	fmt.Println("This part wil create sendgrid client")

	// Create database client
	fmt.Println("This part will create database client")

	// Create service with clients interfaces
	service := service.New(fileReader)

	err = service.Run()
	if err != nil {
		log.Printf("an error ocurred while running the service: %s", err.Error())
	}

	log.Println("Finished...")
}
