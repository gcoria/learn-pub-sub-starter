package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		command := words[0]
		switch command {
		case "pause":
			fmt.Println("Sending pause message")
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
		case "resume":
			fmt.Println("Sending resume message")
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				fmt.Printf("could not publish message: %v", err)
			}
		case "quit":
			fmt.Println("Quitting...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}
