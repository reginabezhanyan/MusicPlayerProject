package playlist

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	ErrorEmptyTitleSong       = errors.New("The title of the song cannot be empty")
	ErrorNotValidDurationSong = errors.New("The duration of the song must be greater than zero")
	ErrorEmptyPlaylist        = errors.New("The playlist is empty")
	ErrorPlayingPlaylist      = errors.New("The playlist is already playing")
	ErrorNotPlayingPlaylist   = errors.New("The playlist is not playing")
	ErrorPausedPlaylist       = errors.New("The playlist is paused")
	ErrorPlayingSong          = errors.New("The song is playing now")
	ErrorNotFoundSong         = errors.New("The song is not found")
)

type Song struct {
	Title    string
	Duration time.Duration
}

type IBasePlaybackMusicPlayer interface {
	Play() error
	Pause() error
	Next() error
	Prev() error
	AddSong(title string, duration time.Duration) error
	DeleteSong(title string) error
	UpdateSong(oldTitle string, newTitle string, newDuration time.Duration) error
}

type playlist struct {
	songs         *list.List
	currentSong   *list.Element
	isPlaying     bool
	isPaused      bool
	playbackMutex sync.Mutex
	playbackCond  *sync.Cond
	stopChan      chan struct{}
}

func NewPlaylist() IBasePlaybackMusicPlayer {
	p := &playlist{
		songs:    list.New(),
		stopChan: make(chan struct{}),
	}
	p.playbackCond = sync.NewCond(&p.playbackMutex)
	return p
}

func (p *playlist) AddSong(title string, duration time.Duration) error {
	if title == "" {
		return ErrorEmptyTitleSong
	}
	if duration <= 0 {
		return ErrorNotValidDurationSong
	}

	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	song := &Song{Title: title, Duration: duration}
	p.songs.PushBack(song)
	return nil
}

func (p *playlist) Play() error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.songs.Len() == 0 {
		return ErrorEmptyPlaylist
	}

	if p.isPlaying {
		if !p.isPaused {
			return ErrorPlayingPlaylist
		}

		p.isPaused = false
		p.playbackCond.Signal()
		return nil
	}

	if p.currentSong == nil {
		p.currentSong = p.songs.Front()
	}

	p.isPlaying = true
	p.isPaused = false

	var wg sync.WaitGroup
	wg.Add(1)
	go p.playback(&wg)
	go func() {
		wg.Wait()
	}()
	return nil
}

func (p *playlist) Pause() error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if !p.isPlaying {
		return ErrorNotPlayingPlaylist
	}

	if p.isPaused {
		return ErrorPausedPlaylist
	}

	p.isPaused = true
	return nil
}

func (p *playlist) Next() error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.currentSong == nil {
		return ErrorEmptyPlaylist
	}

	if p.currentSong.Next() == nil {
		p.currentSong = p.songs.Front()
	} else {
		p.currentSong = p.currentSong.Next()
	}

	if p.isPlaying && !p.isPaused {
		close(p.stopChan)
		p.stopChan = make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)
		go p.playback(&wg)

		go func() {
			wg.Wait()
		}()
	}
	return nil
}

func (p *playlist) Prev() error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.currentSong == nil {
		return ErrorEmptyPlaylist
	}

	if p.currentSong.Prev() == nil {
		p.currentSong = p.songs.Back()
	} else {
		p.currentSong = p.currentSong.Prev()
	}

	if p.isPlaying && !p.isPaused {
		close(p.stopChan)
		p.stopChan = make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)
		go p.playback(&wg)
		go func() {
			wg.Wait()
		}()
	}
	return nil
}

func (p *playlist) DeleteSong(title string) error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.songs.Len() == 0 {
		return ErrorEmptyPlaylist
	}

	for e := p.songs.Front(); e != nil; e = e.Next() {
		song := e.Value.(*Song)
		if song.Title == title {
			if e == p.currentSong {
				return ErrorPlayingSong
			}
			p.songs.Remove(e)
			return nil
		}
	}

	return ErrorNotFoundSong
}

func (p *playlist) UpdateSong(oldTitle string, newTitle string, newDuration time.Duration) error {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.songs.Len() == 0 {
		return ErrorEmptyPlaylist
	}

	if newTitle == "" {
		return ErrorEmptyTitleSong
	}

	for e := p.songs.Front(); e != nil; e = e.Next() {
		song := e.Value.(*Song)
		if song.Title == oldTitle {
			song.Title = newTitle
			song.Duration = newDuration
			return nil
		}
	}

	return ErrorNotFoundSong
}

func (p *playlist) playback(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		p.playbackMutex.Lock()

		for p.isPaused {
			p.playbackCond.Wait()
		}

		if !p.isPlaying {
			p.playbackMutex.Unlock()
			return
		}

		song := p.currentSong.Value.(*Song)
		p.playbackMutex.Unlock()

		select {
		case <-time.After(song.Duration):
			p.Next()
		case <-p.stopChan:
			return
		}
	}
}
