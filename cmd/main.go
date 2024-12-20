package main

import (
	"log"
	"net"

	db_song "MusicPlayerProject/internal/db"
	"MusicPlayerProject/internal/grpcserver"
	"MusicPlayerProject/internal/usecase"
	pb "MusicPlayerProject/proto"
	"database/sql"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@db:5432/playlist?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Successfully connected to the database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	} else {
		log.Println("Successfully ping to the database")
	}

	repo := db_song.NewSongDB(db)
	controller := usecase.NewPlaylistController(repo)
	grpcServerInstance := grpcserver.NewGRPCServer(controller)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	} else {
		log.Println("Successfully listen on port 8080")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPlaylistServiceServer(grpcServer, grpcServerInstance)

	reflection.Register(grpcServer)

	log.Println("gRPC server is running on port 8080...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
