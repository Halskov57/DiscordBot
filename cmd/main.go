package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dietzy1/discordbot/src/bot"
	"github.com/dietzy1/discordbot/src/config"
	"github.com/dietzy1/discordbot/src/repository"
)

func main() {
	//Cannot proceed without a .env file so must be fatal.
	err := config.ReadEnvfile()
	if err != nil {
		log.Fatal(err)
	}

	//Create a new repository
	repo, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}

	//Inject repo dependency into bot application
	bot, err := bot.New(repo)
	if err != nil {
		log.Fatal(err)
	}

	//Run the bot
	err = bot.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Wait for interrupt signal to gracefully shutdown the bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	//Close the connection
	err = bot.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection closed, bot is shutting down")
}
