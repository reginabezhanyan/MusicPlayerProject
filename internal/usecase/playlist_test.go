package usecase

import (
	"MusicPlayerProject/internal/data"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSongDB struct {
	mock.Mock
}

func (m *MockSongDB) Create(ctx context.Context, song *data.Song) (int, error) {
	args := m.Called(ctx, song)
	return args.Int(0), args.Error(1)
}

func (m *MockSongDB) Get(ctx context.Context, title string) (*data.Song, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(*data.Song), args.Error(1)
}

func (m *MockSongDB) Update(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error {
	args := m.Called(ctx, oldTitle, newTitle, duration)
	return args.Error(0)
}

func (m *MockSongDB) Delete(ctx context.Context, title string) error {
	args := m.Called(ctx, title)
	return args.Error(0)
}

func (m *MockSongDB) List(ctx context.Context) ([]*data.Song, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*data.Song), args.Error(1)
}

type MockPlaybackMusicPlayer struct {
	mock.Mock
}

func (m *MockPlaybackMusicPlayer) Play() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlaybackMusicPlayer) Pause() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlaybackMusicPlayer) Next() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlaybackMusicPlayer) Prev() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlaybackMusicPlayer) AddSong(title string, duration time.Duration) error {
	args := m.Called(title, duration)
	return args.Error(0)
}

func TestCreateSong(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()
	song := &data.Song{
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}

	mockRepo.On("Create", ctx, song).Return(1, nil)
	mockRepo.On("Get", ctx, "Test Song").Return((*data.Song)(nil), nil)

	id, err := controller.CreateSong(ctx, song.Title, song.Duration)
	assert.NoError(t, err, "expected no error, but got: %v", err)
	assert.Equal(t, 1, id, "expected song ID to be 1, but got: %d", id)

	mockRepo.AssertCalled(t, "Create", ctx, song)
}
func TestGetSong(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()

	expectedSong := &data.Song{
		ID:       1,
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}

	mockRepo.On("Get", ctx, "Test Song").Return(expectedSong, nil)

	song, err := controller.GetSong(ctx, "Test Song")
	assert.NoError(t, err, "expected no error, but got: %v", err)
	assert.Equal(t, expectedSong, song, "expected song to match, but got a mismatch")

	mockRepo.AssertCalled(t, "Get", ctx, "Test Song")
}
func TestUpdateSong(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()

	song := &data.Song{
		Title:    "Old Title",
		Duration: 2 * time.Minute,
	}
	updatedSong := *song
	updatedSong.Title = "New Title"

	mockRepo.On("Create", ctx, song).Return(1, nil)
	mockRepo.On("Update", ctx, song.Title, updatedSong.Title, updatedSong.Duration).Return(nil)
	mockRepo.On("Get", ctx, "Old Title").Return((*data.Song)(nil), nil)
	controller.CreateSong(ctx, song.Title, song.Duration)
	mockRepo.On("Get", ctx, "Old Title").Return(song, nil)
	mockRepo.On("Get", ctx, "New Title").Return((*data.Song)(nil), nil)

	err := controller.UpdateSong(ctx, "Old Title", "New Title", 2*time.Minute)
	assert.Error(t, err, "expected no error, but got: %v", err)
}

func TestDeleteSong(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()

	song := &data.Song{
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}
	mockRepo.On("Delete", ctx, "Test Song").Return(nil)
	mockRepo.On("Create", ctx, song).Return(1, nil)
	mockRepo.On("Get", ctx, "Test Song").Return((*data.Song)(nil), nil)

	controller.CreateSong(ctx, song.Title, song.Duration)

	mockRepo.On("Get", ctx, "Test Song").Return(&data.Song{Title: "Test Song", Duration: 3 * time.Minute}, nil)

	err := controller.DeleteSong(ctx, "Test Song")
	assert.Error(t, err, "expected no error, but got: %v", err)

	mockRepo.AssertCalled(t, "Get", ctx, "Test Song")
}

func TestListSongs(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()

	expectedSongs := []*data.Song{
		{ID: 1, Title: "Song 1", Duration: 2 * time.Minute},
		{ID: 2, Title: "Song 2", Duration: 3 * time.Minute},
	}

	mockRepo.On("List", ctx).Return(expectedSongs, nil)

	songs, err := controller.ListSongs(ctx)
	assert.NoError(t, err, "expected no error, but got: %v", err)
	assert.Equal(t, expectedSongs, songs, "expected songs list to match")

	mockRepo.AssertCalled(t, "List", ctx)
}

func TestPlayPause(t *testing.T) {
	mockRepo := new(MockSongDB)
	controller := NewPlaylistController(mockRepo)

	ctx := context.Background()

	song := &data.Song{
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}
	mockRepo.On("Create", ctx, song).Return(1, nil)
	mockRepo.On("Get", ctx, "Test Song").Return((*data.Song)(nil), nil)
	controller.CreateSong(ctx, song.Title, song.Duration)

	err := controller.PlaySong(context.Background())
	assert.NoError(t, err, "expected no error on Play, but got: %v", err)

	err = controller.PauseSong(context.Background())
	assert.NoError(t, err, "expected no error on Pause, but got: %v", err)
}
