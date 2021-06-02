package audio

import (
	"fmt"
	"io"
	"sync"

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

// encodes PCM music and sends it to the discordgo Voice Connection
func Encode(in chan []int16, wg *sync.WaitGroup, out *discordgo.VoiceConnection) {
	defer wg.Done()
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
		if err == io.EOF {
			return
		}
		if err != io.ErrUnexpectedEOF && err != nil {
			fmt.Println("Error encode:", err)
			return
		}
		out.OpusSend <- encoded
	}
}
