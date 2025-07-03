package converter

import (
	"github.com/GolZrd/easy-grpc/internal/repository/note/model"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromRepo(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.CreatedAt)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNoteInfoFromRepo(info model.Info) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
