package util

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ParseArgs(text string) string {
	query := strings.SplitAfterN(text, " ", 2)
	if len(query) < 2 {
		return ""
	}
	return query[1]
}

func GetChannelID(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return "", err
	}
	return channel.ID, nil
}

func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
	c, err := GetChannelID(s, m)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(c, message)
	if err != nil {
		return err
	}
	return nil
}

//#f542bf

func SendMessageEmbed(s *discordgo.Session, m *discordgo.MessageCreate, message string, color int) error {
	c, err := GetChannelID(s, m)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(c, &discordgo.MessageEmbed{
		Description: message,
		Color:       color,
	})
	if err != nil {
		return err
	}
	return nil
}

func AddReact(s *discordgo.Session, m *discordgo.MessageCreate, react string) error {
	c, err := GetChannelID(s, m)
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(c, m.ID, react)
	if err != nil {
		return err
	}
	return nil
}
