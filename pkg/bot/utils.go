package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// variables format
const (
	variableUserId         = "%sId"
	variableUserName       = "%sName"
	variableChannelId      = "%sId"
	variableChannelName    = "%sName"
	variableRoleId         = "%sId"
	variableRoleName       = "%sName"
	variableAttachmentId   = "%sId"
	variableAttachmentName = "%sName"
	variableAttachmentUrl  = "%sUrl"
)

func generateHelpEmbed(cmd *discordgo.ApplicationCommand) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for _, option := range cmd.Options {
		if option.Name == "device" {
			continue
		}

		embedValue := fmt.Sprintf("%s\n**Required:** %s\n", option.Description, strconv.FormatBool(option.Required))
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser:
			varId := fmt.Sprintf(variableUserId, option.Name)
			varName := fmt.Sprintf(variableUserName, option.Name)

			embedValue += fmt.Sprintf("**Output variables:** `%%%s` `%%%s`", varId, varName)

		case discordgo.ApplicationCommandOptionChannel:
			varId := fmt.Sprintf(variableChannelId, option.Name)
			varName := fmt.Sprintf(variableChannelName, option.Name)

			embedValue += fmt.Sprintf("**Output variables:** `%%%s` `%%%s`", varId, varName)

		case discordgo.ApplicationCommandOptionRole:
			varId := fmt.Sprintf(variableRoleId, option.Name)
			varName := fmt.Sprintf(variableRoleName, option.Name)

			embedValue += fmt.Sprintf("**Output variables:** `%%%s` `%%%s`", varId, varName)

		case discordgo.ApplicationCommandOptionAttachment:
			varId := fmt.Sprintf(variableAttachmentId, option.Name)
			varName := fmt.Sprintf(variableAttachmentName, option.Name)
			varUrl := fmt.Sprintf(variableAttachmentUrl, option.Name)

			embedValue += fmt.Sprintf("**Output variables:** `%%%s` `%%%s` `%%%s`", varId, varName, varUrl)

		default:
			embedValue += fmt.Sprintf("**Output variables:** `%%%s`", option.Name)
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Option: `%s`", option.Name),
			Value:  embedValue,
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

func convertVariables(s *discordgo.Session, i *discordgo.InteractionCreate) (map[string]string, error) {
	result := make(map[string]string)
	options := i.ApplicationCommandData().Options

	// no configured variables found
	if len(options) <= 1 {
		return result, nil
	}

	for _, option := range options[1:] {
		switch option.Type {
		case discordgo.ApplicationCommandOptionString:
			result[option.Name] = option.Value.(string)

		case discordgo.ApplicationCommandOptionInteger, discordgo.ApplicationCommandOptionNumber:
			result[option.Name] = strconv.FormatFloat(option.Value.(float64), 'f', -1, 64)

		case discordgo.ApplicationCommandOptionBoolean:
			result[option.Name] = strconv.FormatBool(option.Value.(bool))

		case discordgo.ApplicationCommandOptionUser:
			id := option.Value.(string)
			member, err := s.GuildMember(i.GuildID, id)
			if err != nil {
				return nil, fmt.Errorf("error fetching member info: %w", err)
			}
			result[fmt.Sprintf(variableUserId, option.Name)] = member.User.ID
			result[fmt.Sprintf(variableUserName, option.Name)] = member.User.Username

		case discordgo.ApplicationCommandOptionChannel:
			id := option.Value.(string)
			channel, err := s.Channel(id)
			if err != nil {
				return nil, fmt.Errorf("error fetching channel info: %w", err)
			}
			result[fmt.Sprintf(variableChannelId, option.Name)] = channel.ID
			result[fmt.Sprintf(variableChannelName, option.Name)] = channel.Name

		case discordgo.ApplicationCommandOptionRole:
			id := option.Value.(string)
			role, err := s.State.Role(i.GuildID, id)
			if err != nil {
				return nil, fmt.Errorf("error fetching role info: %w", err)
			}
			result[fmt.Sprintf(variableRoleId, option.Name)] = role.ID
			result[fmt.Sprintf(variableRoleName, option.Name)] = role.Name

		case discordgo.ApplicationCommandOptionAttachment:
			id := option.Value.(string)
			attachment, ok := i.ApplicationCommandData().Resolved.Attachments[id]
			if !ok {
				return nil, fmt.Errorf("error fetching attachment info id %s", id)
			}
			result[fmt.Sprintf(variableAttachmentId, option.Name)] = attachment.ID
			result[fmt.Sprintf(variableAttachmentName, option.Name)] = attachment.Filename
			result[fmt.Sprintf(variableAttachmentUrl, option.Name)] = attachment.URL
		}
	}

	return result, nil
}
