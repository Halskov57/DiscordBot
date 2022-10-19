package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/dietzy1/discordbot/src/repository"
)

type bot struct {
	s    *discordgo.Session
	repo repository.Repository
	cmd  map[string]commandFn
}

// Constructor to inject repo dependency into bot application
func New(repo repository.Repository) (*bot, error) {
	return &bot{repo: repo}, nil
}

// Creates a new discord session, registers the handlers, commands and at the end it opens a websocket connection to the discord gateway
func (b *bot) Run() error {
	s, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	b.s = s
	b.registerHandlers()

	//Open websocket connection
	err = s.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	log.Println("Connected to Discord")

	//Conditionally delete commands
	//Potentially add command line flag to delete commands - for now hardcoding is fine
	boolio := false
	if boolio {
		b.deleteCommands()
	}

	// Call function that contains constructor for commands
	b.registerCommands(commands())

	//Instanstiate map holding commands
	b.newCommandFn()

	return nil
}

func (b *bot) Close() error {
	err := b.s.Close()
	if err != nil {
		log.Fatalf("Error closing connection: %v", err)
	}
	return nil
}

func (b *bot) registerHandlers() {
	b.s.AddHandler(b.emoteHandler)
	b.s.AddHandler(b.interactionHandler)
	b.s.AddHandler(b.emoteReactionHandler)
}

// Aceepts a slice of commands and registers them to the discord api
func (b *bot) registerCommands(commands []*discordgo.ApplicationCommand) error {
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		//Should be left empty to register commands globally, for now its filled with the test server to instantly register

		cmd, err := b.s.ApplicationCommandCreate(os.Getenv("APPID"), "", v)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
			return err
		}
		registeredCommands[i] = cmd
	}
	log.Printf("Registered commands: %v", registeredCommands)
	return nil
}

func (b *bot) deleteCommands() error {
	commands, err := b.s.ApplicationCommands(os.Getenv("APPID"), "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range commands {
		err := b.s.ApplicationCommandDelete(os.Getenv("APPID"), "", v.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	return nil
}
