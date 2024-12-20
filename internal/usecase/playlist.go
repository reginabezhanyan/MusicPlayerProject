package usecase

import (
	"MusicPlayerProject/internal/data"
	db_song "MusicPlayerProject/internal/db"
	"MusicPlayerProject/internal/playlist"
	"context"
	"errors"
	"time"
)

type IPlaylistController interface {
	CreateSong(ctx context.Context, title string, duration time.Duration) (int, error)
	GetSong(ctx context.Context, title string) (*data.Song, error)
	UpdateSong(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error
	DeleteSong(ctx context.Context, title string) error
	ListSongs(ctx context.Context) ([]*data.Song, error)
	PlaySong(ctx context.Context) error
	PauseSong(ctx context.Context) error
	NextSong(ctx context.Context) error
	PrevSong(ctx context.Context) error
}

type playlistController struct {
	db       db_song.SongDB
	playlist playlist.IBasePlaybackMusicPlayer
}

func NewPlaylistController(db db_song.SongDB) IPlaylistController {
	playlist := playlist.NewPlaylist()
	return &playlistController{db: db, playlist: playlist}
}

var (
	ErrSongNotPlaying       = errors.New("No song is currently playing")
	ErrSongStillPlaying     = errors.New("The song is still playing, stop it before deleting")
	ErrorSongExised         = errors.New("The song with this title already exists in the database")
	ErrorNotFoundSongOnBase = errors.New("The song is not found on database")
	ErrorNilPlaylist        = errors.New("The playlist cannot be nil")
)

func (c *playlistController) CreateSong(ctx context.Context, title string, duration time.Duration) (int, error) {
	if title == "" {
		return 0, playlist.ErrorEmptyTitleSong
	}
	if duration <= 0 {
		return 0, playlist.ErrorNotValidDurationSong
	}

	song, err := c.db.Get(ctx, title)
	if err != nil {
		return 0, err
	}
	if song != nil {
		return 0, ErrorSongExised
	}

	song = &data.Song{
		Title:    title,
		Duration: duration,
	}

	id, err := c.db.Create(ctx, song)
	if err != nil {
		return 0, err
	}

	err = c.playlist.AddSong(title, duration)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *playlistController) GetSong(ctx context.Context, title string) (*data.Song, error) {
	song, err := c.db.Get(ctx, title)
	if err != nil {
		return nil, err
	}
	if song == nil {
		return nil, playlist.ErrorNotFoundSong
	}
	return song, nil
}

func (c *playlistController) UpdateSong(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error {
	if newTitle == "" {
		return playlist.ErrorEmptyTitleSong
	}
	if duration <= 0 {
		return playlist.ErrorNotValidDurationSong
	}

	song, err := c.db.Get(ctx, oldTitle)
	if err != nil {
		return err
	}
	if song == nil {
		return ErrorNotFoundSongOnBase
	}

	song, err = c.db.Get(ctx, newTitle)
	if err != nil {
		return err
	}
	if song != nil {
		return ErrorSongExised
	}

	err = c.db.Update(ctx, oldTitle, newTitle, duration)
	if err != nil {
		return err
	}

	err = c.playlist.UpdateSong(oldTitle, newTitle, duration)
	if err != nil {
		return err
	}

	return nil
}

func (c *playlistController) DeleteSong(ctx context.Context, title string) error {
	song, err := c.db.Get(ctx, title)
	if err != nil {
		return err
	}
	if song == nil {
		return ErrorNotFoundSongOnBase
	}

	err = c.playlist.DeleteSong(song.Title)
	if err != nil {
		return err
	}

	err = c.db.Delete(ctx, title)
	if err != nil {
		return err
	}

	return nil
}

func (c *playlistController) ListSongs(ctx context.Context) ([]*data.Song, error) {
	return c.db.List(ctx)
}

func (c *playlistController) PlaySong(ctx context.Context) error {
	if c.playlist == nil {
		return ErrorNilPlaylist
	}
	return c.playlist.Play()
}

func (c *playlistController) PauseSong(ctx context.Context) error {
	if c.playlist == nil {
		return ErrorNilPlaylist
	}
	return c.playlist.Pause()
}

func (c *playlistController) NextSong(ctx context.Context) error {
	if c.playlist == nil {
		return ErrorNilPlaylist
	}
	return c.playlist.Next()
}

func (c *playlistController) PrevSong(ctx context.Context) error {
	if c.playlist == nil {
		return ErrorNilPlaylist
	}
	return c.playlist.Prev()
}
