package playlist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddSong(t *testing.T) {
	p := NewPlaylist().(*playlist)

	// check add song
	err := p.AddSong("Song 1", 3*time.Second)
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to have 1 song")
	assert.Equal(t, "Song 1", p.songs.Front().Value.(*Song).Title, "expected the title of the first song to be 'Song 1'")

	// check add song with bad title
	err = p.AddSong("", time.Second)
	assert.Equal(t, ErrorEmptyTitleSong, err, "expected error %v, but get: %v", ErrorEmptyTitleSong, err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to still have 1 song after failed addition")

	// check add song with bad duration
	err = p.AddSong("Song 2", 0)
	assert.Equal(t, ErrorNotValidDurationSong, err, "expected error %v, but get: %v", ErrorNotValidDurationSong, err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to still have 1 song after failed addition")

	err = p.AddSong("Song 2", -5*time.Second)
	assert.Equal(t, ErrorNotValidDurationSong, err, "expected error %v, but get: %v", ErrorNotValidDurationSong, err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to still have 1 song after failed addition")

	err = p.AddSong("Song 2", 5*time.Second)
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, 2, p.songs.Len(), "expected playlist to still have 1 song after failed addition")
}

func TestUpdateSong(t *testing.T) {
	p := NewPlaylist().(*playlist)

	err := p.AddSong("Song 1", 3*time.Second)

	err = p.UpdateSong("Song 1", "Song 11", 3*time.Second)
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to have 1 song")
	assert.Equal(t, "Song 11", p.songs.Front().Value.(*Song).Title, "expected the title of the first song to be 'Song 11'")

	err = p.UpdateSong("Song 11", "Song 11", 4*time.Second)
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to have 1 song")
	assert.Equal(t, 4*time.Second, p.songs.Front().Value.(*Song).Duration, "expected the duration of the first song to be 4 second")

	// check update song with bad title
	err = p.UpdateSong("Song 111", "Song", 4*time.Second)
	assert.Equal(t, ErrorNotFoundSong, err, "expected error %v, but get: %v", ErrorEmptyTitleSong, err)

	err = p.UpdateSong("Song 11", "", 4*time.Second)
	assert.Equal(t, ErrorEmptyTitleSong, err, "expected error %v, but get: %v", ErrorEmptyTitleSong, err)
}

func TestPlay(t *testing.T) {
	p := NewPlaylist().(*playlist)

	p.AddSong("Song 1", 1*time.Second)
	p.AddSong("Song 2", 1*time.Second)
	p.AddSong("Song 3", 1*time.Second)

	err := p.Play()
	assert.NoError(t, err, "expected no error, but got: %v", err)

	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, "Song 2", p.currentSong.Value.(*Song).Title, "expected 'Song 2' to be playing")

	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, "Song 3", p.currentSong.Value.(*Song).Title, "expected 'Song 3' to be playing")
}
func TestPlayAndNext(t *testing.T) {
	p := NewPlaylist().(*playlist)

	p.AddSong("Song 1", 150*time.Second)
	p.AddSong("Song 2", 150*time.Second)

	err := p.Play()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 1", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 1'")

	err = p.Next()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 2", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 2'")

	// check circle playing
	err = p.Next()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 1", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 1'")
}

func TestPlayAndPrev(t *testing.T) {
	p := NewPlaylist().(*playlist)

	p.AddSong("Song 1", 150*time.Second)
	p.AddSong("Song 2", 150*time.Second)

	err := p.Play()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 1", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 1'")

	err = p.Next()
	err = p.Prev()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 1", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 1'")

	// check circle playing
	err = p.Prev()
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, "Song 2", p.currentSong.Value.(*Song).Title, "expected the title of the current song to be 'Song 2'")
}

func TestDeleteSong(t *testing.T) {
	p := NewPlaylist().(*playlist)

	p.AddSong("Song 1", 10*time.Second)
	err := p.DeleteSong("Song 1")
	assert.NoError(t, err, "expected no error, but get: %v", err)
	assert.Equal(t, 0, p.songs.Len(), "expected playlist to have 0 song")

	p.AddSong("Song 1", 10*time.Second)
	p.Play()

	err = p.DeleteSong("Song 1")
	assert.Equal(t, ErrorPlayingSong, err, "expected error %v, but get: %v", ErrorEmptyTitleSong, err)
	assert.Equal(t, 1, p.songs.Len(), "expected playlist to still have 1 song after failed addition")
}
