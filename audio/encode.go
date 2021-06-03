package audio

import (
	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	FrameRate = 48000
	FrameSize = 960
	Channels  = 2
)

// encodes PCM music and sends it to the discordgo Voice Connection
func Encode(in <-chan []int16, status chan<- error, out *discordgo.VoiceConnection) {

	encoder, err := gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		status <- err
		return
	}
	encoder.SetBitrate(96000)
	for {
		raw, ok := <-in
		if !ok {
			status <- nil
			return
		}
		encoded, err := encoder.Encode(raw, FrameSize, FrameSize*4)
		if err != nil {
			status <- err
			return
		}
		out.OpusSend <- encoded
	}
}
