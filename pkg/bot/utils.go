package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func generateHelpEmbed(cmd *discordgo.ApplicationCommand) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for _, option := range cmd.Options {
		if option.Name == "device" {
			continue
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Option: `%s`", option.Name),
			Value:  fmt.Sprintf("%s\n**Required:** %s\n**Output variables:** `%%%v`", option.Description, strconv.FormatBool(option.Required), option.Name),
			Inline: false,
		})
	}

	embeds := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Command /%s", cmd.Name),
		Description: cmd.Description,
		Color:       0x00ff00,
		Fields:      fields,
	}

	return embeds
}
