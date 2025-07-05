package converter

import (
	"github.com/GolZrd/easy-grpc/internal/model"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromService(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.CreatedAt)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromService(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNoteInfoFromService(info model.NoteInfo) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}

func ToNoteInfoFromDesc(info *desc.NoteInfo) *model.NoteInfo {
	return &model.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
