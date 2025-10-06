package bot

import (
	"discord-tasker-runner/pkg/tasker"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handleTasker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// retrieve data
	cmd := i.ApplicationCommandData().Name
	targetDevice := i.ApplicationCommandData().Options[0].Value.(string)
	targetDeviceToken := deviceToken[targetDevice]
	targetTask := commandTask[cmd]

	err := tasker.ExecuteTask(projectId, targetDeviceToken, targetTask, nil)
	if err != nil {
		log.Printf("error execute command /%v on device %v: %v\n", cmd, targetDevice, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "❌ **Error**\n\nCheck logs for detail",
						Color:       0xff0000,
					},
				},
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "✅ **Success**",
						Color:       0x00ff00,
					},
				},
			},
		})
	}
}
