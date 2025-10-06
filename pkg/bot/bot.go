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

var commandHandlers = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate))

func Run(cfg config.Config) {
	session, err := discordgo.New("Bot " + cfg.Token)
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

		if len(cmdConfig.Variables) > 0 {
			var cmdOptions []*discordgo.ApplicationCommandOption
			for _, variable := range cmdConfig.Variables {
				options := &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        variable.Name,
					Description: variable.Description,
					Required:    variable.Required,
				}
				cmdOptions = append(cmdOptions, options)
			}
			cmd.Options = cmdOptions
		}

		commandHandlers[cmd.Name] = handleTest
		commands = append(commands, cmd)
	}

	// Register each command
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("error creating command %s: %w", cmd.Name, err)
		}
	}

	fmt.Printf("Successfully registered %d slash commands.\n", len(commands))
	return nil
}

func handleTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "hello",
		},
	})
}

func interactions(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	if handler, ok := commandHandlers[commandName]; ok {
		handler(s, i)
	}
}
