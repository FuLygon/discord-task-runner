package bot

import (
	"discord-tasker-runner/pkg/tasker"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handleTasker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("error deferring interaction: ", err)
		return
	}

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

	err = tasker.ExecuteTask(projectId, targetDeviceToken, targetTask, variables)
	if err != nil {
		log.Printf("error execute command /%v on device %v: %v\n", cmd, targetDevice, err)
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("**Error executing %s**", targetTask),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("error sending error message: ", err)
			return
		}
	} else {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Executing **%s**", targetTask),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("error sending success message: ", err)
			return
		}
	}
}
