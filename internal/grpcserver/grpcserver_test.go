package grpcserver

import (
	"MusicPlayerProject/internal/data"
	pb "MusicPlayerProject/proto"
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type MockPlaylistController struct {
	mock.Mock
}

func (m *MockPlaylistController) CreateSong(ctx context.Context, title string, duration time.Duration) (int, error) {
	args := m.Called(ctx, title, duration)
	return args.Int(0), args.Error(1)
}

func (m *MockPlaylistController) GetSong(ctx context.Context, title string) (*data.Song, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(*data.Song), args.Error(1)
}

func (m *MockPlaylistController) UpdateSong(ctx context.Context, oldTitle string, newTitle string, duration time.Duration) error {
	args := m.Called(ctx, oldTitle, newTitle, duration)
	return args.Error(0)
}

func (m *MockPlaylistController) DeleteSong(ctx context.Context, title string) error {
	args := m.Called(ctx, title)
	return args.Error(0)
}

func (m *MockPlaylistController) ListSongs(ctx context.Context) ([]*data.Song, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*data.Song), args.Error(1)
}

func (m *MockPlaylistController) PlaySong(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPlaylistController) PauseSong(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPlaylistController) NextSong(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPlaylistController) PrevSong(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func bufDialer(mockController *MockPlaylistController) (*grpc.ClientConn, func(), error) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	grpcServer := NewGRPCServer(mockController)

	pb.RegisterPlaylistServiceServer(server, grpcServer)

	go func() {
		_ = server.Serve(lis)
	}()

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	closeConn := func() {
		conn.Close()
		server.Stop()
	}

	return conn, closeConn, err
}

func TestCreateSong(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	expectedID := 1
	mockController.On("CreateSong",
		mock.Anything,
		"Test Song",
		3*time.Minute,
	).Return(expectedID, nil)

	req := &pb.CreateSongRequest{
		Title:    "Test Song",
		Duration: int64(3 * time.Minute.Seconds()),
	}

	resp, err := client.CreateSong(context.Background(), req)

	assert.NoError(t, err, "unexpected error during CreateSong gRPC call")
	assert.Equal(t, expectedID, int(resp.Id), "expected song ID to match")
	mockController.AssertCalled(t, "CreateSong", mock.Anything, mock.Anything, mock.Anything)
}
func TestGetSong(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	expectedSong := &data.Song{
		ID:       1,
		Title:    "Test Song",
		Duration: 3 * time.Minute,
	}
	mockController.On("GetSong", mock.Anything, "Test Song").Return(expectedSong, nil)

	req := &pb.GetSongRequest{Title: "Test Song"}

	resp, err := client.GetSong(context.Background(), req)

	assert.NoError(t, err, "unexpected error during GetSong gRPC call")
	assert.Equal(t, expectedSong.ID, int(resp.Id), "expected song ID to match")
	assert.Equal(t, expectedSong.Title, resp.Title, "expected song Title to match")
	assert.Equal(t, expectedSong.Duration.Seconds(), float64(resp.Duration), "expected song Duration to match")
	mockController.AssertCalled(t, "GetSong", mock.Anything, "Test Song")
}

func TestListSongs(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	expectedSongs := []*data.Song{
		{ID: 1, Title: "Song 1", Duration: 2 * time.Minute},
		{ID: 2, Title: "Song 2", Duration: 4 * time.Minute},
	}
	mockController.On("ListSongs", mock.Anything).Return(expectedSongs, nil)

	req := &pb.EmptyMessage{}

	resp, err := client.ListSongs(context.Background(), req)

	assert.NoError(t, err, "unexpected error during ListSongs gRPC call")
	assert.Len(t, resp.Songs, 2, "expected two songs in the list")
	assert.Equal(t, "Song 1", resp.Songs[0].Title, "expected Song 1 to match")
	assert.Equal(t, "Song 2", resp.Songs[1].Title, "expected Song 2 to match")

	mockController.AssertCalled(t, "ListSongs", mock.Anything)
}

func TestDeleteSong(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("DeleteSong", mock.Anything, "Test Song").Return(nil)

	req := &pb.DeleteSongRequest{Title: "Test Song"}

	_, err = client.DeleteSong(context.Background(), req)
	assert.NoError(t, err, "unexpected error during DeleteSong gRPC call")

	mockController.AssertCalled(t, "DeleteSong", mock.Anything, "Test Song")
}

func TestUpdateSong(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("UpdateSong", mock.Anything, "Old Title", "Updated Title", 240*time.Second).
		Return(nil)

	req := &pb.UpdateSongRequest{
		OldTitle: "Old Title",
		NewTitle: "Updated Title",
		Duration: int64(4 * time.Minute.Seconds()),
	}

	resp, err := client.UpdateSong(context.Background(), req)
	assert.NoError(t, err, "unexpected error during UpdateSong gRPC call")
	assert.Equal(t, "Updated Title", resp.Title, "expected the response Title to match")

	mockController.AssertCalled(t, "UpdateSong", mock.Anything, "Old Title", "Updated Title", 240*time.Second)
}

func TestPlay(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("PlaySong", mock.Anything).Return(nil)

	_, err = client.Play(context.Background(), &pb.EmptyMessage{})
	assert.NoError(t, err, "unexpected error during Play gRPC call")

	mockController.AssertCalled(t, "PlaySong", mock.Anything)
}

func TestPause(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("PauseSong", mock.Anything).Return(nil)

	_, err = client.Pause(context.Background(), &pb.EmptyMessage{})
	assert.NoError(t, err, "unexpected error during Pause gRPC call")

	mockController.AssertCalled(t, "PauseSong", mock.Anything)
}

func TestNext(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("NextSong", mock.Anything).Return(nil)

	_, err = client.Next(context.Background(), &pb.EmptyMessage{})
	assert.NoError(t, err, "unexpected error during Next gRPC call")

	mockController.AssertCalled(t, "NextSong", mock.Anything)
}

func TestPrev(t *testing.T) {
	mockController := new(MockPlaylistController)
	conn, cleanup, err := bufDialer(mockController)
	assert.NoError(t, err)
	defer cleanup()

	client := pb.NewPlaylistServiceClient(conn)

	mockController.On("PrevSong", mock.Anything).
		Return(nil)

	_, err = client.Prev(context.Background(), &pb.EmptyMessage{})
	assert.NoError(t, err, "unexpected error during Prev gRPC call")

	mockController.AssertCalled(t, "PrevSong", mock.Anything)
}
