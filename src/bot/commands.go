package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	embedBuilder "github.com/clinet/discordgo-embed"
	"github.com/dietzy1/discordbot/src/bot/emotes"
)

type commandFn func(context.Context, *discordgo.Session, *discordgo.InteractionCreate)

func (b *bot) newCommandFn() {
	b.cmd = map[string]commandFn{
		"emotes":      b.emotes,
		"leaderboard": b.leaderboard,
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
	}
}

func (b *bot) emotes(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	opts := makeOptionMap(i.ApplicationCommandData().Options)
	targetUsername := opts["user-option"].UserValue(s).Username

	emote := &emotes.Emote{
		Guild: i.GuildID,
		User:  targetUsername,
	}

	emotes, err := b.repo.GetUserEmotes(ctx, emote)
	if err != nil {
		errorResponse(s, i, err)
		return
	}

	//Build the embed
	embed := embedBuilder.NewEmbed()
	embed.SetTitle(targetUsername + "'s emotes")
	embed.Timestamp = time.Now().Format(time.RFC3339)
	j := 0

	for _, v := range emotes {
		j++
		embed.AddField(v.Emote+fmt.Sprintf(" %d  ", v.Count), "\u200b")
		if j == 10 {
			break
		}
	}
	builtEmbed := append([]*discordgo.MessageEmbed{}, embed.MessageEmbed)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: builtEmbed,
		},
	})
	if err != nil {
		log.Println(err)
		errorResponse(s, i, err)
	}
}

func (b *bot) leaderboard(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	//Needs to retrieve the emote from the interaction
	opts := makeOptionMap(i.ApplicationCommandData().Options)
	targetEmote := opts["string-option"].StringValue()

	emote := &emotes.Emote{
		Guild: i.GuildID,
		Emote: targetEmote,
	}

	emotes, err := b.repo.GetServerEmote(ctx, emote)
	if err != nil {
		errorResponse(s, i, err)
		return
	}
	// Build the embed
	embed := embedBuilder.NewEmbed()
	embed.SetTitle(targetEmote + " Leaderboard")
	embed.Timestamp = time.Now().Format(time.RFC3339)
	j := 0
	medal := []string{"🥇", "🥈", "🥉", "🏅", "🏅", "🏅", "🏅", "🏅", "🏅", "🏅"}
	for k, v := range emotes {
		j++
		embed.AddField(fmt.Sprintf("%v - %v %v  ", medal[k], v.User, v.Count), "\u200b")
		if j == 10 {
			break
		}
	}
	builtEmbed := append([]*discordgo.MessageEmbed{}, embed.MessageEmbed)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: builtEmbed,
		},
	})
	if err != nil {
		log.Println(err)
		errorResponse(s, i, err)
	}
}
