package repository

//go:generate cmd /c "if exist mocks rmdir /s /q mocks && mkdir mocks"
//go:generate minimock -i ./note.NoteRepository -o ./mocks/ -s "_minimock.go"
