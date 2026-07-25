// Package main provides the home-ip-updater service that monitors a RabbitMQ queue
// for IP address updates and automatically updates DNS records in AWS Route53.
// This service is designed to work with the home-ip-monitor system to keep
// a subdomain pointing to the current home IP address.
package main

import (
	//"context"
	"log"
	"log/syslog"
	"os"
	//"os/signal"
	//"syscall"

	//messagebroker "github.com/a-castellano/go-services/services/messagebroker"
	//updater "github.com/a-castellano/home-ip-updater/internal/app/updater"
	config "github.com/a-castellano/home-ip-updater/internal/infra/config"
)

// main is the entry point of the home-ip-updater service.
// It sets up logging, configuration, message broker connection,
// signal handling, and starts the main message processing loop.
func main() {

	// Configure logger to write to syslog for system integration
	// This allows the service to integrate with system logging infrastructure
	logwriter, e := syslog.New(syslog.LOG_INFO, "home-ip-updater")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove timestamp from log messages as syslog already provides timestamps
		log.SetFlags(0)
	}

	log.Print("Loading configuration from environment variables")

	// Load application configuration from environment variables
	// This validates all required AWS, RabbitMQ, and domain settings
	//appConfig, configErr := config.NewConfig()
	_, configErr := config.NewConfig()

	if configErr != nil {
		log.Print(configErr.Error())
		os.Exit(1)
	}

	log.Print("Creating RabbitMQ client for message consumption")

	// Create a cancellable context for graceful shutdown
	//	ctx, cancel := context.WithCancel(context.Background())

	// // Initialize RabbitMQ client and message broker
	// rabbitmqClient := messagebroker.NewRabbitmqClient(appConfig.RabbitmqConfig)
	// messageBroker := messagebroker.MessageBroker{Client: rabbitmqClient}
	//
	// // Create channels for message processing and error handling
	// messagesReceived := make(chan []byte)
	// receiveErrors := make(chan error)
	//
	// log.Print("Setting up OS signal handling for graceful shutdown")
	//
	// // Set up signal handling for graceful shutdown (SIGTERM, SIGINT)
	// signalChannel := make(chan os.Signal, 2)
	// signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	//
	// // Start signal handler goroutine
	//
	//	go func() {
	//		sig := <-signalChannel
	//		switch sig {
	//		case os.Interrupt:
	//			log.Print("Received SIGINT, initiating graceful shutdown")
	//			cancel()
	//		case syscall.SIGTERM:
	//			log.Print("Received SIGTERM, initiating graceful shutdown")
	//			cancel()
	//		}
	//	}()
	//
	// // Start message consumer in background
	// go messageBroker.ReceiveMessages(ctx, appConfig.UpdateQueue, messagesReceived, receiveErrors)
	//
	// log.Print("Starting main message processing loop")
	//
	// // Main processing loop - handles messages and errors
	//
	//	for {
	//		select {
	//		case receivedError := <-receiveErrors:
	//			// Handle RabbitMQ connection or message processing errors
	//			log.Print(receivedError.Error())
	//			os.Exit(1)
	//		case messageReceived := <-messagesReceived:
	//			// Process received IP address update
	//			ipReceived := string(messageReceived)
	//			log.Printf("Received new IP address to update: %s", ipReceived)
	//			log.Printf("Updating DNS record for subdomain: %s", appConfig.Subdomain)
	//
	//			// Create AWS updater instance with current configuration
	//			awsUpdater := updater.AWSUpdater{
	//				ZoneID:    appConfig.AWSZoneID,
	//				Subdomain: appConfig.Subdomain,
	//				IP:        ipReceived,
	//			}
	//
	//			// Update DNS record in AWS Route53
	//			updateErr := awsUpdater.Update(ctx)
	//			if updateErr != nil {
	//				log.Printf("Failed to update DNS record: %s", updateErr.Error())
	//			} else {
	//				log.Printf("Successfully updated DNS record for %s to IP %s", appConfig.Subdomain, ipReceived)
	//			}
	//
	//		case <-ctx.Done():
	//			// Handle graceful shutdown
	//			log.Print("Shutdown signal received, terminating service")
	//			os.Exit(0)
	//		}
	//	}
}
