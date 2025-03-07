package bot

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

//fuck det hele mand

func (b *bot) emoteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf(m.Author.ID)

}

// Inject the commands map into the interaction handler
func (b *bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if fn, ok := b.cmd[i.ApplicationCommandData().Name]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		defer cancel()
		fn(ctx, s, i)
	}
}

func (b *bot) emoteReactionHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

}

func (b *bot) helloHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello <@"+m.Author.ID+">")
	}
}
