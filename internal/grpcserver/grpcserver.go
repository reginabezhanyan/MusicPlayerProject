package grpcserver

import (
	"MusicPlayerProject/internal/usecase"
	pb "MusicPlayerProject/proto"
	"context"
	"time"
)

type GRPCServer struct {
	controller usecase.IPlaylistController
	pb.UnimplementedPlaylistServiceServer
}

func NewGRPCServer(controller usecase.IPlaylistController) *GRPCServer {
	return &GRPCServer{controller: controller}
}

func (s *GRPCServer) CreateSong(ctx context.Context, req *pb.CreateSongRequest) (*pb.SongResponse, error) {
	id, err := s.controller.CreateSong(ctx, req.Title, time.Duration(req.Duration)*time.Second)
	if err != nil {
		return nil, err
	}

	return &pb.SongResponse{Id: int32(id), Title: req.Title, Duration: req.Duration}, nil
}

func (s *GRPCServer) GetSong(ctx context.Context, req *pb.GetSongRequest) (*pb.SongResponse, error) {
	song, err := s.controller.GetSong(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.SongResponse{
		Id:       int32(song.ID),
		Title:    song.Title,
		Duration: int64(song.Duration.Seconds()),
	}, nil
}

func (s *GRPCServer) UpdateSong(ctx context.Context, req *pb.UpdateSongRequest) (*pb.SongResponse, error) {
	err := s.controller.UpdateSong(ctx, req.OldTitle, req.NewTitle, time.Duration(req.Duration)*time.Second)
	if err != nil {
		return nil, err
	}

	return &pb.SongResponse{
		Title:    req.NewTitle,
		Duration: req.Duration,
	}, nil
}

func (s *GRPCServer) DeleteSong(ctx context.Context, req *pb.DeleteSongRequest) (*pb.EmptyMessage, error) {
	err := s.controller.DeleteSong(ctx, req.Title)
	if err != nil {
		return nil, err
	}

	return &pb.EmptyMessage{}, nil
}

func (s *GRPCServer) ListSongs(ctx context.Context, req *pb.EmptyMessage) (*pb.ListSongsResponse, error) {
	songs, err := s.controller.ListSongs(ctx)
	if err != nil {
		return nil, err
	}

	var songResponses []*pb.SongResponse
	for _, song := range songs {
		songResponses = append(songResponses, &pb.SongResponse{
			Id:       int32(song.ID),
			Title:    song.Title,
			Duration: int64(song.Duration.Seconds()),
		})
	}

	return &pb.ListSongsResponse{Songs: songResponses}, nil
}

func (s *GRPCServer) Play(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.controller.PlaySong(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *GRPCServer) Pause(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.controller.PauseSong(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *GRPCServer) Next(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.controller.NextSong(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *GRPCServer) Prev(ctx context.Context, req *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.controller.PrevSong(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}
