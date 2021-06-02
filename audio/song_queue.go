package audio

import (
	"math/rand"
	"sync"
	"time"
)

type SongQueue struct {
	queue     []*Song
	loader    *Loader
	queueLock sync.Mutex
	queueSem  chan bool
}

func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue:    make([]*Song, 0),
		loader:   NewLoader(),
		queueSem: make(chan bool, AudioQueueMax),
	}
}

func (songQueue *SongQueue) Add(songs ...*Song) {
	songQueue.queueLock.Lock()

	songQueue.queue = append(songQueue.queue, songs...)

	songQueue.loader.Enqueue(songs...)

	songQueue.queueLock.Unlock()

	for range songs {
		songQueue.queueSem <- true
	}
}

func (songQueue *SongQueue) Get() *Song {
	<-songQueue.queueSem

	songQueue.queueLock.Lock()
	defer songQueue.queueLock.Unlock()

	song := songQueue.queue[0]
	songQueue.queue = songQueue.queue[1:]
	return song
}

func (songQueue *SongQueue) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	songQueue.queueLock.Lock()
	defer songQueue.queueLock.Unlock()

	rand.Shuffle(len(songQueue.queue), func(i, j int) {
		songQueue.queue[i], songQueue.queue[j] = songQueue.queue[j], songQueue.queue[i]
	})
	songQueue.loader.Clear()
	songQueue.loader.Enqueue(songQueue.queue...)
}

func (songQueue *SongQueue) Clear() {
	songQueue.queueLock.Lock()
	defer songQueue.queueLock.Unlock()
	songQueue.loader.Clear()
	songQueue.queue = make([]*Song, 0)
}

func (songQueue *SongQueue) Cleanup() {
	songQueue.loader.Cleanup()
}
