package tests

import (
	"context"
	"fmt"
	"testing"

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

func TestCreate(t *testing.T) {
	t.Parallel()

	type noteRepositoryMockFunc func(mc *minimock.Controller) noteRepository.NoteRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.NoteInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Animal()
		content = gofakeit.Animal()

		repoErr = fmt.Errorf("repository error")

		req = &model.NoteInfo{
			Title:   title,
			Content: content,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		noteRepositoryMock noteRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			noteRepositoryMock: func(mc *minimock.Controller) noteRepository.NoteRepository {
				mock := repoMocs.NewNoteRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocs.NewTxManagerMock(mc)
				// Если метод Create использует транзакции
				mock.ReadCommittedMock.Set(func(ctx context.Context, fn db.Handler) error {
					return fn(ctx)
				})

				return mock
				//Если метод Create не использует транзакции, можно упростить:
				//return txMocs.NewTxManagerMock(mc)
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			noteRepositoryMock: func(mc *minimock.Controller) noteRepository.NoteRepository {
				mock := repoMocs.NewNoteRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocs.NewTxManagerMock(mc)
				// Если метод Create использует транзакции
				mock.ReadCommittedMock.Set(func(ctx context.Context, fn db.Handler) error {
					return fn(ctx)
				})

				return mock
				//Если метод Create не использует транзакции, можно упростить:
				//return txMocs.NewTxManagerMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Параллельно выполняем подтесты
			t.Parallel()

			noteRepoMock := tt.noteRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := note.NewService(noteRepoMock, txManagerMock)

			res, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
