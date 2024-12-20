package db_song

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"MusicPlayerProject/internal/data"
)

type SongDB interface {
	Create(ctx context.Context, song *data.Song) (int, error)
	Get(ctx context.Context, title string) (*data.Song, error)
	Update(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error
	Delete(ctx context.Context, title string) error
	List(ctx context.Context) ([]*data.Song, error)
}

type songPostgreSQL struct {
	db *sql.DB
}

func NewSongDB(db *sql.DB) SongDB {
	return &songPostgreSQL{db: db}
}

func (r *songPostgreSQL) Create(ctx context.Context, song *data.Song) (int, error) {
	var id int
	query := `
		INSERT INTO songs (title, duration)
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, song.Title, song.Duration.Seconds()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *songPostgreSQL) Get(ctx context.Context, title string) (*data.Song, error) {
	query := `
		SELECT id, title, duration
		FROM songs
		WHERE title = $1
	`

	var song data.Song
	var durationSeconds int64

	err := r.db.QueryRowContext(ctx, query, title).Scan(&song.ID, &song.Title, &durationSeconds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	song.Duration = time.Duration(durationSeconds) * time.Second
	return &song, nil
}

func (r *songPostgreSQL) Update(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error {
	query := `
		UPDATE songs
		SET title = $2, duration = $3
		WHERE title = $1
	`

	res, err := r.db.ExecContext(ctx, query, oldTitle, newTitle, duration.Seconds())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows updated, check the song title")
	}

	return nil
}

func (r *songPostgreSQL) Delete(ctx context.Context, title string) error {
	query := `
		DELETE FROM songs
		WHERE title = $1
	`

	res, err := r.db.ExecContext(ctx, query, title)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows deleted, check the song ID")
	}

	return nil
}

func (r *songPostgreSQL) List(ctx context.Context) ([]*data.Song, error) {
	query := `
		SELECT id, title, duration
		FROM songs
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*data.Song

	for rows.Next() {
		var song data.Song
		var durationSeconds int64

		err = rows.Scan(&song.ID, &song.Title, &durationSeconds)
		if err != nil {
			return nil, err
		}

		song.Duration = time.Duration(durationSeconds) * time.Second
		songs = append(songs, &song)
	}

	return songs, rows.Err()
}
