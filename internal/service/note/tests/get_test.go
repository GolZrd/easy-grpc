package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/GolZrd/easy-grpc/internal/client/db"
	txMocs "github.com/GolZrd/easy-grpc/internal/client/mocks"
	"github.com/GolZrd/easy-grpc/internal/model"
	repoMocs "github.com/GolZrd/easy-grpc/internal/repository/mocks"
	noteRepository "github.com/GolZrd/easy-grpc/internal/repository/note"
	"github.com/GolZrd/easy-grpc/internal/service/note"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type noteRepositoryMockFunc func(mc *minimock.Controller) noteRepository.NoteRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		title     = gofakeit.Animal()
		content   = gofakeit.Animal()
		createdAt = time.Now()
		updatedAt = time.Now()

		repoErr = fmt.Errorf("repository error")

		res = &model.Note{
			ID: id,
			Info: model.NoteInfo{
				Title:   title,
				Content: content,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)

	tests := []struct {
		name               string
		args               args
		want               *model.Note
		err                error
		noteRepositoryMock noteRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			noteRepositoryMock: func(mc *minimock.Controller) noteRepository.NoteRepository {
				mock := repoMocs.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
			// Здесь бес передачи транзакции, так как get ее не требует, поэтому просто возвращаем
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				return txMocs.NewTxManagerMock(mc)
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			noteRepositoryMock: func(mc *minimock.Controller) noteRepository.NoteRepository {
				mock := repoMocs.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				return txMocs.NewTxManagerMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Параллельно выполняем подтесты
			t.Parallel()

			noteRepositoryMock := tt.noteRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := note.NewService(noteRepositoryMock, txManagerMock)

			res, err := service.Get(ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}
