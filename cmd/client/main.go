package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	_ = username

	connectionUrl := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	fmt.Println("Connected to RabbitMQ")

	// Wait for a signal (e.g. Ctrl+C) to exit the program.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Connection closed, program shutting down")
}
