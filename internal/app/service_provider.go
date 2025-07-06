package app

import (
	"context"
	"log"

	noteAPI "github.com/GolZrd/easy-grpc/internal/api/note"
	"github.com/GolZrd/easy-grpc/internal/client/db"
	"github.com/GolZrd/easy-grpc/internal/client/db/pg"
	"github.com/GolZrd/easy-grpc/internal/closer"
	"github.com/GolZrd/easy-grpc/internal/config"
	noteRepository "github.com/GolZrd/easy-grpc/internal/repository/note"
	noteService "github.com/GolZrd/easy-grpc/internal/service/note"
)

// Описываем структуру для хранения зависимостей
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	noteRepository noteRepository.NoteRepository

	noteService noteService.NoteService

	noteImpl *noteAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Функция, которая проверяет был ли уже инициализирован конфиг для БД. Если нет, то инициализируем его. и возвращаем.
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) NoteRepository(ctx context.Context) noteRepository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepository.NewRepository(s.PgPool(ctx))
	}
	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) noteService.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(s.NoteRepository(ctx))
	}

	return s.noteService
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *noteAPI.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = noteAPI.NewImplementation(s.NoteService(ctx))
	}

	return s.noteImpl
}
