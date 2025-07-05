package converter

import (
	"github.com/GolZrd/easy-grpc/internal/model"
	modelRepo "github.com/GolZrd/easy-grpc/internal/repository/note/model"
)

// Описываем конвертеры между репозиторием и сервисом.
// modelRepo - модель внутри репо слоя, а model - внутри бизнес логики, то есть слоя Service
func ToNoteFromRepo(note *modelRepo.Note) *model.Note {

	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(info modelRepo.NoteInfo) model.NoteInfo {
	return model.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
