package note

import (
	service "github.com/GolZrd/easy-grpc/internal/service/note"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
