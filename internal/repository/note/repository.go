package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/model"
	"github.com/GolZrd/easy-grpc/internal/repository/note/converter"
	modelRepo "github.com/GolZrd/easy-grpc/internal/repository/note/model"
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
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) NoteRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
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

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	builder := squirrel.Select(idColumn, titleColumn, contentColumn, CreatedAtColumn, UpdatedAtColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(tableName).Where(squirrel.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var note modelRepo.Note
	err = r.db.QueryRow(ctx, query, args...).Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
