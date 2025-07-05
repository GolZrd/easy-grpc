package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/model"
)

func (s *service) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	id, err := s.NoteRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}
