package bot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dietzy1/discordbot/src/bot/emotes"
	stonkclientv1 "github.com/dietzy1/discordbot/src/proto/stonk/v1"
)

//fuck det hele mand

func (b *bot) emoteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "test" {
		log.Println(b.client.GetStonk(context.Background(), &stonkclientv1.GetStonkRequest{Stonk: "SPY"}))
		log.Println(b.client.Gainers(context.Background(), &stonkclientv1.GainersRequest{}))
		log.Println(b.client.Loosers(context.Background(), &stonkclientv1.LoosersRequest{}))
		log.Println(b.client.Compare(context.Background(), &stonkclientv1.CompareRequest{}))
	}

	re := regexp.MustCompile(`<:\w+:\d{10,45}>`)
	if len(re.FindString(m.Content)) >= 1 {
		Emote := re.FindString(m.Content)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		defer cancel()
		emote := &emotes.Emote{
			Guild: m.GuildID,
			User:  m.Author.Username,
			Emote: Emote,
		}

		err := b.repo.IncrementEmote(ctx, emote)
		if err != nil {
			return
		}
	}
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
	emote := &emotes.Emote{
		Guild: m.GuildID,
		User:  m.Member.User.Username,
		Emote: fmt.Sprintf("<:%s:%s>", m.Emoji.Name, m.Emoji.ID),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	err := b.repo.IncrementEmote(ctx, emote)
	if err != nil {
		return
	}
}
