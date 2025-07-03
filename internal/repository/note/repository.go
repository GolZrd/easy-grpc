package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/repository/note/converter"
	"github.com/GolZrd/easy-grpc/internal/repository/note/model"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	CreatedAtColumn = "created_at"
	UpdatedAtColumn = "updated_at"
)

type NoteRepository interface {
	Create(ctx context.Context, info *desc.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.Note, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) NoteRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, info *desc.NoteInfo) (int64, error) {
	builder := squirrel.Insert(tableName).PlaceholderFormat(squirrel.Dollar).Columns(titleColumn, contentColumn).Values(info.Title, info.Content).Suffix("RETURNING id")
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.Note, error) {
	builder := squirrel.Select(idColumn, titleColumn, contentColumn, CreatedAtColumn, UpdatedAtColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(tableName).Where(squirrel.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var note model.Note
	err = r.db.QueryRow(ctx, query, args...).Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
