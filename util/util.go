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
	s.ChannelMessageSend(c, message)
	return nil
}
