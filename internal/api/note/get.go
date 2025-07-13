package note

import (
	"context"
	"log"

	"github.com/GolZrd/easy-grpc/internal/converter"
	"github.com/GolZrd/easy-grpc/internal/logger"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"go.uber.org/zap"
)

// Описываем метод Get
func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// Добавим логирование
	logger.Info("Getting note...", zap.Int64("id", req.Id))

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
