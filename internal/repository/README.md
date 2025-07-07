# Для генерации моков использовался minimock и чтобы указать путь до интерфейса, который лежал в папке note, использовалась следующая команда
# -i ./note.NoteRepository
# Полностью команда выглядит так: go:generate minimock -i ./note.NoteRepository -o ./mocks/ -s "_minimock.go"