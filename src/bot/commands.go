package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type commandFn func(context.Context, *discordgo.Session, *discordgo.InteractionCreate)

func (b *bot) newCommandFn() {
	b.cmd = map[string]commandFn{
		"emotes":      b.emotes,
		"leaderboard": b.leaderboard,
		"hello":       b.hello,
	}
}

// The commands are registered correctly
func commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		//emotes command
		{
			Name:        "emotes",
			Description: "Returns a list of emotes used by a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "The user to get emotes from",
					Required:    true,
				},
			},
		},
		//leaderboard command
		{
			Name:        "leaderboard",
			Description: "Returns a leaderboard of who uses a specific emote the most",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "The emote to get the leaderboard from",
					Required:    true,
				},
			},
		},
		//hello command
		{
			Name:        "hello",
			Description: "Says hello",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "The string to say hello to",
					Required:    true,
				},
			},
		},
	}
}

func (b *bot) emotes(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
}

func (b *bot) leaderboard(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
}

func (b *bot) hello(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Hello command invoked by %v", i.Member.User.ID)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello <@" + i.Member.User.ID + ">!",
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		// Handle error
		fmt.Println("Error responding to interaction:", err)
	}
}
