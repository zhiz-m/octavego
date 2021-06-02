package audio

import (
	"sync"
)

type Song struct {
	loadChan    chan bool
	isWaited    bool
	isLoaded    bool
	lock        sync.Mutex
	SourceURL   string
	SearchQuery string
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
