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
	variables := make(map[string]string)

	// handle variables
	if len(i.ApplicationCommandData().Options) > 1 {
		for _, options := range i.ApplicationCommandData().Options[1:] {
			variables[options.Name] = options.Value.(string)
		}
	}

	err := tasker.ExecuteTask(projectId, targetDeviceToken, targetTask, variables)
	if err != nil {
		log.Printf("error execute command /%v on device %v: %v\n", cmd, targetDevice, err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		if err != nil {
			log.Println("error sending error message", err)
			return
		}
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		if err != nil {
			log.Println("error sending success message", err)
			return
		}
	}
}
