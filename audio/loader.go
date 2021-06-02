package audio

import (
	"fmt"
)

type Loader struct {
	work chan *Song
	kill chan bool
}

func NewLoader() *Loader {
	loader := &Loader{
		work: make(chan *Song, AudioLoaderMax),
		kill: make(chan bool, 1),
	}
	go loader.workLoop()
	return loader
}

func (loader *Loader) Cleanup() {
	loader.kill <- true
}

func (loader *Loader) workLoop() {
	for {
		select {
		case <-loader.kill:
			return
		default:
		}
		song := <-loader.work
		song.lock.Lock()
		if song.isLoaded {
			song.lock.Unlock()
			return
		}
		song.isLoaded = true
		song.lock.Unlock()

		sourceURL, err := YTDL(song)
		if err != nil {
			fmt.Printf("Error loading song %v\n", song.SearchQuery)
		}

		song.lock.Lock()
		song.SourceURL = sourceURL
		song.Load()
		// todo: add json
		song.lock.Unlock()
	}
}

func (loader *Loader) Enqueue(songs ...*Song) {
	for _, song := range songs {
		loader.work <- song
	}
}

func (loader *Loader) Clear() {
	for {
		select {
		case <-loader.work:
		default:
			return
		}
	}
}

/*
func (loader *Loader) workLoop() {
	// channels to communicate with YTDL
	// need 2 in the rare scenario that loader.loaderChan and YTDLReady are both ready
	YTDLReady := make(chan bool)
	YTDLKill := make(chan bool)
	for {
		select {
		// parent wants to kill the loader goroutine
		case <-loader.kill:
			YTDLKill <- true
			return

		// YTDL goroutine finished, load another one
		case <-YTDLReady:
			// block until one item is available
			songs := []*Song{
				<-loader.work,
			}
			var song *Song
		L:
			for {
				select {
				case song = <-loader.work:
					songs = append(songs, song)
				default:
					break L
				}
			}
			go YTDL(songs, YTDLKill, YTDLReady)
		}
	}
}
*/
