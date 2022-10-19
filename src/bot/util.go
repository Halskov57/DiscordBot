package bot

import "github.com/bwmarrin/discordgo"

type OptionsMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func makeOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) (m OptionsMap) {
	m = make(OptionsMap, len(options))

	for _, option := range options {
		m[option.Name] = option
	}

	return
}

func errorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + err.Error(),
		},
	})
}
