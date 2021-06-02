package audio

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	FrameRate = 48000
	FrameSize = 960
	Channels  = 2
)

type KillError struct {
	error
}

func (error *KillError) Error() string {
	return "PlaySong was killed"
}

func Encode(in chan []int16, out *discordgo.VoiceConnection) {
	encoder, err := gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		print("error creating pcm encoder")
		return
	}
	for {
		raw, ok := <-in
		if !ok {
			return
		}
		encoded, err := encoder.Encode(raw, FrameSize, FrameSize)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("Error encode:", err)
			return
		}
		out.OpusSend <- encoded
	}
}
