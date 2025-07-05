package note

import (
	"context"
	"log"

	"github.com/GolZrd/easy-grpc/internal/converter"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
)

// Описываем метод Create
func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("Inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
