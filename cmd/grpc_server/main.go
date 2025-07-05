package main

import (
	"context"
	"log"
	"net"

	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/jackc/pgx/v4/pgxpool"

	noteAPI "github.com/GolZrd/easy-grpc/internal/api/note"
	"github.com/GolZrd/easy-grpc/internal/config"
	noteRepository "github.com/GolZrd/easy-grpc/internal/repository/note"
	noteService "github.com/GolZrd/easy-grpc/internal/service/note"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с БД
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	noteRepo := noteRepository.NewRepository(pool)
	//Добавляем сервис
	noteService := noteService.NewService(noteRepo)

	s := grpc.NewServer()
	// включаем reflection на сервере
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, noteAPI.NewImplementation(noteService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
