package audio

import (
	"fmt"
	"sync"
	"time"
)

type Song struct {
	loadChan chan bool
	isWaited bool
	isLoaded bool
	lock     sync.Mutex

	SourceURL   string
	SearchQuery string
	title       string
	author      string
	duration    *time.Time
}

func NewSongBySearch(query string) *Song {
	return &Song{
		SearchQuery: query,
		loadChan:    make(chan bool, 1),
	}
}

func NewSongByURL(URL string) *Song {
	return &Song{
		SourceURL: URL,
		loadChan:  make(chan bool, 1),
		isWaited:  true,
		isLoaded:  true,
	}
}

// waits for the song to be fully loaded
func (song *Song) WaitLoad() {
	song.lock.Lock()
	isWaited := song.isWaited
	song.lock.Unlock()
	if !isWaited {
		<-song.loadChan
		song.lock.Lock()
		song.isWaited = true
		song.lock.Unlock()
	}
}

// signals that the song is fully loaded. Must only be called once
func (song *Song) Load() {
	song.loadChan <- true
}

func (song *Song) String() string {
	song.lock.Lock()
	defer song.lock.Unlock()
	title := song.title
	author := song.author
	t := song.duration
	duration := "unknown duration"
	if title == "" {
		title = "unknown"
		if song.SearchQuery != "" {
			title = song.SearchQuery
		}
	}
	if author == "" {
		author = "unknown"
	}
	if t != nil {
		duration = string(t.String())
	}
	return fmt.Sprintf("%v by %v: %v\n", title, author, duration)
}
