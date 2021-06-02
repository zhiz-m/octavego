package audio

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SongQueue struct {
	queue    []*Song
	loader   *Loader
	lock     sync.Mutex
	queueSem chan bool
}

func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue:    make([]*Song, 0),
		loader:   NewLoader(),
		queueSem: make(chan bool, AudioQueueMax),
	}
}

func (songQueue *SongQueue) Add(songs ...*Song) {
	songQueue.lock.Lock()

	songQueue.queue = append(songQueue.queue, songs...)

	songQueue.loader.Enqueue(songs...)

	songQueue.lock.Unlock()

	for range songs {
		songQueue.queueSem <- true
	}
}

func (songQueue *SongQueue) Get() *Song {
	<-songQueue.queueSem

	songQueue.lock.Lock()
	defer songQueue.lock.Unlock()

	song := songQueue.queue[0]
	songQueue.queue = songQueue.queue[1:]
	return song
}

func (songQueue *SongQueue) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	songQueue.lock.Lock()
	defer songQueue.lock.Unlock()

	rand.Shuffle(len(songQueue.queue), func(i, j int) {
		songQueue.queue[i], songQueue.queue[j] = songQueue.queue[j], songQueue.queue[i]
	})
	songQueue.loader.Clear()
	songQueue.loader.Enqueue(songQueue.queue...)
}

func (songQueue *SongQueue) Clear() {
	songQueue.lock.Lock()
	songQueue.loader.Clear()
	songQueue.queue = make([]*Song, 0)
	songQueue.lock.Unlock()

	for {
		select {
		case <-songQueue.queueSem:
		default:
			return
		}
	}
}
func (songQueue *SongQueue) String() string {
	var text string
	songQueue.lock.Lock()
	defer songQueue.lock.Unlock()
	for i, song := range songQueue.queue {
		text = text + fmt.Sprintf("%v. %v", i+1, song)
	}
	if text == "" {
		text = "*empty*"
	}
	return text
}

func (songQueue *SongQueue) Cleanup() {
	songQueue.loader.Cleanup()
}
