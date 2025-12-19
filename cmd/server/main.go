package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionUrl := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")
	err = pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		fmt.Printf("could not publish message: %v", err)
	}

	// Wait for a signal (e.g. Ctrl+C) to exit the program.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Connection closed, program shutting down")
}
