package audio

import (
	"time"
)

const (
	AudioLoaderMax = 1000
	AudioQueueMax  = 1000

	AudioTimeout = 10 * time.Minute

	VoiceChannelErrorPrompt     = "**Please join a voice channel**"
	RemoveAudioStateErrorPrompt = "**Not currently in a channel**"
	HelpPrompt                  = "**Invalid command**"
)

var ()
