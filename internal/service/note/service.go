package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/client/db"
	"github.com/GolZrd/easy-grpc/internal/model"
	"github.com/GolZrd/easy-grpc/internal/repository/note"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}

type service struct {
	NoteRepository note.NoteRepository
	txManager      db.TxManager
}

func NewService(noteRepository note.NoteRepository, txManager db.TxManager) NoteService {
	return &service{NoteRepository: noteRepository, txManager: txManager}
}
