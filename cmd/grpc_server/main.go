package main

import (
	"context"
	"log"
	"net"

	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/GolZrd/easy-grpc/internal/config"
	"github.com/GolZrd/easy-grpc/internal/converter"
	"github.com/GolZrd/easy-grpc/internal/repository/note"
	service "github.com/GolZrd/easy-grpc/internal/service/note"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Заменили repository на noteService, так как добавил новый слой бизнес логики
type server struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

// Описываем метод Create
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("Inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Описываем метод Get
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := s.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", noteObj.ID, noteObj.Info.Title,
		noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)
	return &desc.GetResponse{
		Note: converter.ToNoteFromService(noteObj),
	}, nil
}

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

	noteRepo := note.NewRepository(pool)
	//Добавляем сервис
	noteService := service.NewService(noteRepo)

	s := grpc.NewServer()
	// включаем reflection на сервере
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{noteService: noteService})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
