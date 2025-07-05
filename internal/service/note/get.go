package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/model"
)

func (s *service) Get(ctx context.Context, id int64) (*model.Note, error) {
	note, err := s.NoteRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return note, nil
}
