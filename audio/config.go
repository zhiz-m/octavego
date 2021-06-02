package audio

import (
	"time"
)

const (
	AudioLoaderMax = 1000
	AudioQueueMax  = 1000

	AudioTimeout = 10 * time.Minute

	EmbedColor = 0xf542bf

	VoiceChannelErrorPrompt     = "**Error: please join a voice channel**"
	RemoveAudioStateErrorPrompt = "**Error: currently in a channel**"
	PlayErrorPrompt             = "**Error: resource not found**"
	SkipErrorPrompt             = "**Error: no audio is playing**"
	PauseErrorPrompt            = "**Error: no audio is playing**"
	ResumeErrorPrompt           = "**Error: audio is already playing**"
	LoopErrorPrompt             = "**Error: failed to loop queue**"

	HelpPrompt = "**Invalid command**"
)

var ()
