package audio

import (
	"fmt"
	"sync"
)

type Song struct {
	loadChan chan bool
	isWaited bool
	isLoaded bool
	lock     sync.Mutex

	SourceURL   string
	SearchQuery string
	title       string
	artist      string
	duration    int
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

func (song *Song) SetMetadata(title, artist string, numSeconds int) {
	song.title = title
	song.artist = artist
	song.duration = numSeconds
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
	author := song.artist
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
	if t > 0 {
		var minutes int = t / 60
		var seconds int = t - minutes*60
		duration = fmt.Sprintf("%d:%02d", minutes, seconds)
	}
	return fmt.Sprintf("%v by %v: %v\n", title, author, duration)
}
