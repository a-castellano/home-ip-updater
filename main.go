package main

import (
	"context"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	messagebroker "github.com/a-castellano/go-services/messagebroker"
	config "github.com/a-castellano/home-ip-updater/config"
	"github.com/a-castellano/home-ip-updater/powerdnsclient"
	updater "github.com/a-castellano/home-ip-updater/updater"
)

func main() {

	// Configure logger to write to the syslog. You could do this in init(), too.
	logwriter, e := syslog.New(syslog.LOG_INFO, "home-ip-updater")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove timestamp
		log.SetFlags(0)
	}

	// Now from anywhere else in your program, you can use this:
	log.Print("Loading config")

	appConfig, configErr := config.NewConfig()

	if configErr != nil {
		log.Print(configErr.Error())
		os.Exit(1)
	}

	log.Print("Creating RabbitMQ client")
	ctx, cancel := context.WithCancel(context.Background())

	rabbitmqClient := messagebroker.NewRabbimqClient(appConfig.RabbitmqConfig)
	messageBroker := messagebroker.MessageBroker{Client: rabbitmqClient}

	messagesReceived := make(chan []byte)
	receiveErrors := make(chan error)

	log.Print("Define os signal management")
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	log.Print("Start PowerDNS Client")
	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 seconds
	}
	powerDNSClient, clientError := powerdnsclient.NewClient(httpClient, appConfig.PowerDNSHost, appConfig.PowerDNSPort, appConfig.PowerDNSAPIKey)
	if clientError != nil {
		log.Print(clientError.Error())
		os.Exit(1)
	}

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			cancel()
		case syscall.SIGTERM:
			cancel()
		}
	}()

	go messageBroker.ReceiveMessages(ctx, appConfig.UpdateQueue, messagesReceived, receiveErrors)

	log.Print("Waiting for messages")
	for {
		select {
		case receivedError := <-receiveErrors:
			log.Print(receivedError.Error())
			os.Exit(1)
		case messageReceived := <-messagesReceived:
			ipReceived := string(messageReceived)
			log.Printf("Received new IP to update: %s.", ipReceived)
			log.Printf("Updating %s DNS record.", appConfig.Subdomain)

			powerDNSUpdater := updater.PowerDNSUpdater{
				PowerDNSClient: powerDNSClient,
				ZoneName:       appConfig.PowerDNSZoneName,
				Subdomain:      appConfig.Subdomain,
				IP:             ipReceived,
			}

			awsUpdater := updater.AWSUpdater{
				ZoneID:    appConfig.AWSZoneID,
				Subdomain: appConfig.Subdomain,
				IP:        ipReceived,
			}

			updateErr := awsUpdater.Update(ctx)
			if updateErr != nil {
				log.Print(updateErr.Error())
			}

			updateErr = powerDNSUpdater.Update(ctx)
			if updateErr != nil {
				log.Print(updateErr.Error())
			}

		case <-ctx.Done():
			log.Print("Execution finished")
			os.Exit(0)
		}
	}
}
