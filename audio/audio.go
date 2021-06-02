package audio

import (
	"errors"
	"strings"

	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zhiz-m/octavego/util"
)

var (
	audioStates map[string]*AudioState
	prefix      string
)

func AddAudio(s *discordgo.Session, p string) {
	prefix = p
	s.AddHandler(messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix) {
		cmd := m.Content[len(prefix):]
		switch {
		case strings.HasPrefix(cmd, "hi"):
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("hi %s", m.Author.ID))
		case strings.HasPrefix(cmd, "play"):
			_play(s, m, m.Content)
		default:
			s.ChannelMessageSend(m.ChannelID, HelpPrompt)
		}
	}
}

func _play(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	audioState, err := GetAudioState(s, m)
	if err != nil {
		fmt.Println("Error _play:", err)
		return
	}
	query := util.ParseArgs(text)
	audioState.Add(query)
}

func _skip(s *discordgo.Session, m *discordgo.MessageCreate) {

}

func GetAudioState(s *discordgo.Session, m *discordgo.MessageCreate) (*AudioState, error) {
	v, c, g, err := getDiscordInfo(s, m)
	println("hi1", v, c, g, err)
	if err != nil {
		return nil, err
	}
	if audioState, ok := audioStates[g]; ok {
		return audioState, nil
	}
	if v == "" {
		s.ChannelMessageSend(c, VoiceChannelErrorPrompt)
		return nil, errors.New("caller has not joined a voice channel")
	}
	audioState, err := NewAudioState(s, g, v)
	if err != nil {
		return nil, err
	}
	println(audioState)
	return audioState, nil
}

// returns voice channel ID, text channel ID, guild ID, error
func getDiscordInfo(s *discordgo.Session, m *discordgo.MessageCreate) (string, string, string, error) {
	channel, err := s.State.Channel(m.ChannelID)

	if err != nil {
		return "", "", "", err
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return "", "", "", err
	}
	var v string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			v = vs.ChannelID
			break
		}
	}
	return v, channel.ID, guild.ID, nil
}

func Cleanup() {
	for _, audioState := range audioStates {
		audioState.Cleanup()
	}
}
