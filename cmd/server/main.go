package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not declarate the queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()

	i := 0
	for i < 1 {
		Input := gamelogic.GetInput()
		if len(Input) > 0 {
			if Input[0] == "pause" {
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: true,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
				fmt.Println("Pause message sent!")
			} else if Input[0] == "resume" {
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: false,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
			}	else if Input[0] == "quit" {
				i = 1
				fmt.Println("Exiting the server...")
			} else {
				fmt.Println("Unknown command.")
			}
		}
	}
}
