package repository

//go:generate cmd /c "if exist mocks rmdir /s /q mocks && mkdir mocks"
//go:generate minimock -i ./note.NoteService -o ./mocks/ -s "_minimock.go"
