package audio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AudioState struct {
	songQueue   *SongQueue
	vc          *discordgo.VoiceConnection
	kill        chan bool
	skip        chan bool
	pause       chan bool
	resume      chan bool
	isPaused    bool
	isLooping   bool
	currentSong *Song
	lock        sync.Mutex
}

func NewAudioState(s *discordgo.Session, guildID, channelID string) (*AudioState, error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}
	state := &AudioState{
		songQueue:   NewSongQueue(),
		vc:          vc,
		kill:        make(chan bool, 1),
		skip:        make(chan bool, 1),
		pause:       make(chan bool, 1),
		resume:      make(chan bool, 1),
		currentSong: nil,
	}
	go state.workLoop()
	return state, nil
}

func (state *AudioState) workLoop() {
	songChan := make(chan *Song, 1)
	getSong := func() {
		if state.isLooping {
			songChan <- state.currentSong
		} else {
			songChan <- state.songQueue.Get()
		}
	}
	for {
		go getSong()
		select {
		case song := <-songChan:
			state.lock.Lock()
			state.currentSong = song
			state.lock.Unlock()
			song.WaitLoad()
			killed, err := state.play(song.SourceURL)
			if killed {
				return
			}
			if err != nil {
				fmt.Println("Error PlaySong:", err)
			}
		case <-time.After(AudioTimeout):
			return
		case <-state.kill:
			return
		}
	}
}

// returns a boolean representing whether the function was killed, and an error
func (state *AudioState) play(URL string) (bool, error) {
	state.vc.Speaking(true)
	defer state.vc.Speaking(false)

	cmd := exec.Command("ffmpeg", "-i", URL, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	buffer := bufio.NewReaderSize(pipe, 16384)

	var wg sync.WaitGroup
	wg.Add(1)
	inChan := make(chan []int16, 2)
	go Encode(inChan, &wg, state.vc)

	err = cmd.Start()
	if err != nil {
		return false, err
	}
	for {
		// very ugly, probably could be improved
		// key difficulty is that we must wait for signals (kill, skip) whether paused or not
		select {
		case <-state.kill:
			close(inChan)
			return true, nil
		case <-state.skip:
			close(inChan)
			return false, nil
		case <-state.pause:
			state.lock.Lock()
			state.isPaused = true
			state.lock.Unlock()
		L:
			for {
				select {
				case <-state.kill:
					close(inChan)
					return true, nil
				case <-state.skip:
					close(inChan)
					return false, nil
				case <-state.resume:
					state.lock.Lock()
					state.isPaused = false
					state.lock.Unlock()
					break L
				}
			}
		default:
			buf := make([]int16, FrameSize*Channels)
			err = binary.Read(buffer, binary.LittleEndian, &buf)
			if err == io.EOF {
				wg.Wait()
				return false, nil
			}
			if err != io.ErrUnexpectedEOF && err != nil {
				return false, err
			}
			inChan <- buf
		}
	}
}

func (state *AudioState) Cleanup() {
	state.kill <- true
	err := state.vc.Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting from voice connection")
	}
}

func (state *AudioState) Add(query string) bool {
	songs := ProcessQuery(query)
	state.songQueue.Add(songs...)
	return true
}

func (state *AudioState) Clear() bool {
	state.songQueue.Clear()
	return true
}

func (state *AudioState) Shuffle() bool {
	state.songQueue.Shuffle()
	return true
}

func (state *AudioState) Skip() bool {
	state.lock.Lock()
	skip := state.currentSong != nil && !state.isPaused
	state.lock.Unlock()
	if skip {
		state.skip <- true
		return true
	}
	return false
}

func (state *AudioState) IsPaused() bool {
	state.lock.Lock()
	defer state.lock.Unlock()
	return state.isPaused
}

func (state *AudioState) Pause() bool {
	if state.IsPaused() || state.currentSong == nil {
		return false
	}
	state.pause <- true

	return true
}

func (state *AudioState) Resume() bool {
	if !state.IsPaused() {
		return false
	}
	state.resume <- true
	return true
}

func (state *AudioState) Loop() bool {
	state.lock.Lock()
	defer state.lock.Unlock()
	if state.currentSong == nil {
		return false
	}
	state.isLooping = !state.isLooping
	return true
}

func (state *AudioState) String() string {
	state.lock.Lock()
	defer state.lock.Unlock()
	currentSong := "*Not playing*\n"
	if state.currentSong != nil {
		currentSong = fmt.Sprintf("%v", state.currentSong)
	}
	return fmt.Sprintf("**Current Song:**\n%v\n**Queue:**\n%v", currentSong, state.songQueue)
}
