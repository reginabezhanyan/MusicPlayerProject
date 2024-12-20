package db_song

import (
	"MusicPlayerProject/internal/data"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbsong := NewSongDB(db)

	ctx := context.Background()
	song := &data.Song{
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}

	mock.ExpectQuery("INSERT INTO songs").
		WithArgs(song.Title, song.Duration.Seconds()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := dbsong.Create(ctx, song)
	assert.NoError(t, err, "unexpected error when creating a song")
	assert.Equal(t, 1, id, "expected song ID to be 1")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbsong := NewSongDB(db)

	ctx := context.Background()
	expectedSong := &data.Song{
		ID:       1,
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}

	mock.ExpectQuery("SELECT id, title, duration FROM songs WHERE title = \\$1").
		WithArgs("Test Song").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "duration"}).
			AddRow(expectedSong.ID, expectedSong.Title, int64(expectedSong.Duration.Seconds())))

	song, err := dbsong.Get(ctx, "Test Song")
	assert.NoError(t, err, "unexpected error when getting a song")
	assert.Equal(t, expectedSong, song, "expected song to match")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbsong := NewSongDB(db)

	ctx := context.Background()
	song := &data.Song{
		ID:       1,
		Title:    "Test Song",
		Duration: 4 * time.Minute,
	}

	mock.ExpectExec("UPDATE songs SET title = \\$2, duration = \\$3 WHERE title = \\$1").
		WithArgs(song.Title, "Updated Test Song", song.Duration.Seconds()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dbsong.Create(ctx, song)
	err = dbsong.Update(ctx, song.Title, "Updated Test Song", 4*time.Minute)
	assert.NoError(t, err, "unexpected error when updating a song")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbsong := NewSongDB(db)

	ctx := context.Background()
	songTitle := "song"

	mock.ExpectExec("DELETE FROM songs WHERE title = \\$1").
		WithArgs(songTitle).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Одна строка удалена.

	err = dbsong.Delete(ctx, songTitle)
	assert.NoError(t, err, "unexpected error when deleting a song")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListSongs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbsong := NewSongDB(db)

	ctx := context.Background()
	expectedSongs := []*data.Song{
		{ID: 1, Title: "Song 1", Duration: 2 * time.Minute},
		{ID: 2, Title: "Song 2", Duration: 4 * time.Minute},
	}

	mock.ExpectQuery("SELECT id, title, duration FROM songs").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "duration"}).
			AddRow(expectedSongs[0].ID, expectedSongs[0].Title, int64(expectedSongs[0].Duration.Seconds())).
			AddRow(expectedSongs[1].ID, expectedSongs[1].Title, int64(expectedSongs[1].Duration.Seconds())))

	songs, err := dbsong.List(ctx)
	assert.NoError(t, err, "unexpected error when listing songs")
	assert.Equal(t, expectedSongs, songs, "expected songs list to match")

	assert.NoError(t, mock.ExpectationsWereMet())
}
