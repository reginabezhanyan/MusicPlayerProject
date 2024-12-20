package main

import (
	"context"
	"log"
	"time"

	pb "MusicPlayerProject/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewPlaylistServiceClient(conn)

	resp, err := client.CreateSong(context.Background(), &pb.CreateSongRequest{
		Title:    "Test Song 1",
		Duration: int64(3 * time.Minute.Seconds()),
	})
	if err != nil {
		log.Fatalf("CreateSong call failed: %v", err)
	}
	log.Printf("Created song: ID=%d, Title=%s, Duration=%d seconds", resp.Id, resp.Title, resp.Duration)

	respList, errList := client.ListSongs(context.Background(), &pb.EmptyMessage{})
	if errList != nil {
		log.Fatalf("ListSongs call failed: %v", errList)
	}
	log.Printf("ListSongs: %v", respList)

	_, errDel := client.DeleteSong(context.Background(), &pb.DeleteSongRequest{
		Title: "Test Song 1",
	})
	if errDel != nil {
		log.Fatalf("DeleteSong call failed: %v", err)
	}

	client.CreateSong(context.Background(), &pb.CreateSongRequest{
		Title:    "Test Song 1",
		Duration: int64(3 * time.Minute.Seconds()),
	})
	client.CreateSong(context.Background(), &pb.CreateSongRequest{
		Title:    "Test Song 2",
		Duration: int64(3 * time.Minute.Seconds()),
	})
	resp, err = client.CreateSong(context.Background(), &pb.CreateSongRequest{
		Title:    "Test Song 2",
		Duration: int64(3 * time.Minute.Seconds()),
	})
	if err == nil {
		log.Fatalf("Must be error: The song with this title already exists in the database")
	}

	respPlay, errPlay := client.Play(context.Background(), &pb.EmptyMessage{})
	if errPlay != nil {
		log.Fatalf("Play call failed: %v", errList)
	}
	log.Printf("Play: %v", respPlay)

	_, errDel = client.DeleteSong(context.Background(), &pb.DeleteSongRequest{
		Title: "Test Song 1",
	})
	if errDel == nil {
		log.Fatalf("Must be error: The song is playing now")
	}

	_, errDel = client.DeleteSong(context.Background(), &pb.DeleteSongRequest{
		Title: "Test Song 2",
	})
	if errDel != nil {
		log.Fatalf("DeleteSong call failed: %v", err)
	}
}
