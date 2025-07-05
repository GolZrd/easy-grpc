package note

import (
	"context"
	"log"

	"github.com/GolZrd/easy-grpc/internal/converter"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
)

// Описываем метод Get
func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
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
