package bot

import (
	"discord-tasker-runner/pkg/tasker"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("error deferring interaction", err)
		return
	}

	// get all registered commands
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Println("error getting application commands", err)
		return
	}

	// get requested command
	reqCmd := i.ApplicationCommandData().Options[0].Value.(string)

	for _, command := range commands {
		if command.Name == reqCmd {
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{generateHelpEmbed(command)},
			})
			if err != nil {
				log.Println("error sending help success message: ", err)
				return
			}
			return
		}
	}

	// unknown requested command
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Unknown command",
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Println("error sending help error message: ", err)
		return
	}
}

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
	var ttl *string
	cmd := i.ApplicationCommandData().Name
	targetDevice := i.ApplicationCommandData().Options[0].Value.(string)
	targetDeviceToken := deviceToken[targetDevice]
	targetTask := commandTask[cmd]

	// handle variables
	variables, err := convertVariables(s, i)
	if err != nil {
		log.Println(err)
		return
	}

	// handle ttl
	taskTTL, ok := commandTTL[cmd]
	if ok && taskTTL > 0 {
		strTTL := fmt.Sprintf("%ds", taskTTL)
		ttl = &strTTL
	}

	err = tasker.ExecuteTask(projectId, targetDeviceToken, targetTask, ttl, variables)
	if err != nil {
		log.Printf("error execute command /%v on device %v: %v\n", cmd, targetDevice, err)
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("**Error executing %s**", targetTask),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("error sending execute error message: ", err)
			return
		}
	} else {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Executing **%s**", targetTask),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("error sending execute success message: ", err)
			return
		}
	}
}
