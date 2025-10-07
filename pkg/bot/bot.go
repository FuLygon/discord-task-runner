package bot

import (
	"discord-tasker-runner/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	// project id needed to send message
	projectId string
	// token for device
	deviceToken = make(map[string]string)
	// assigned task for the command
	commandTask = make(map[string]string)
)

func Run(cfg config.Config) {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	err = session.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	defer session.Close()

	// set project id
	projectId = cfg.ProjectID

	// remove existing slash commands
	err = removeCommand(session)
	if err != nil {
		log.Println(err)
		return
	}

	// register slash commands
	if err = registerCommandsFromConfig(session, cfg); err != nil {
		log.Println(err)
		return
	}

	// register commands handler
	session.AddHandler(interactions)

	fmt.Printf("Bot invitation url: https://discord.com/oauth2/authorize?client_id=%s&permissions=0&scope=bot%%20applications.commands\n", session.State.User.ID)
	fmt.Println("Bot is running. Ctrl-C to terminate")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func removeCommand(s *discordgo.Session) error {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("error getting commands: %w", err)
	}

	for _, command := range commands {
		err = s.ApplicationCommandDelete(s.State.User.ID, "", command.ID)
		if err != nil {
			return fmt.Errorf("error deleting command: %w", err)
		}
	}
	return nil
}

func registerCommandsFromConfig(s *discordgo.Session, conf config.Config) error {
	var commands []*discordgo.ApplicationCommand

	for _, cmdConfig := range conf.Commands {
		cmd := &discordgo.ApplicationCommand{
			Name:        cmdConfig.Name,
			Description: cmdConfig.Description,
			Type:        discordgo.ChatApplicationCommand,
		}

		// add device option and choices
		devicesOption := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "device",
			Description: "Device to execute task",
			Required:    true,
		}
		for _, device := range conf.Device {
			devicesOption.Choices = append(devicesOption.Choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  device.Name,
				Value: device.Name,
			})
			deviceToken[device.Name] = device.Token
		}
		cmd.Options = append(cmd.Options, devicesOption)

		// assign target task
		commandTask[cmdConfig.Name] = cmdConfig.Task

		// add extra options if configured
		if len(cmdConfig.Variables) > 0 {
			for _, variable := range cmdConfig.Variables {
				option := &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        variable.Name,
					Description: variable.Description,
					Required:    variable.Required,
				}
				cmd.Options = append(cmd.Options, option)
			}
		}

		commands = append(commands, cmd)
	}

	// Register each command
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("error creating command /%s: %w", cmd.Name, err)
		}
	}

	fmt.Printf("Successfully registered %d slash commands.\n", len(commands))
	return nil
}

func interactions(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	handleTasker(s, i)
}
