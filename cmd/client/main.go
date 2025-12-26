package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	connectionUrl := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.SimpleQueueTypeTransient,
	)
	if err != nil {
		fmt.Printf("Failed to declare and bind queue: %v", err)
		return
	}

	fmt.Printf("Queue bound: %s\n", queueName)

	// Wait for a signal (e.g. Ctrl+C) to exit the program.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Connection closed, program shutting down")
}
