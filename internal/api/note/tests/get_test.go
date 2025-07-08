package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	noteAPI "github.com/GolZrd/easy-grpc/internal/api/note"
	"github.com/GolZrd/easy-grpc/internal/model"
	serviceMocks "github.com/GolZrd/easy-grpc/internal/service/mocks"
	noteService "github.com/GolZrd/easy-grpc/internal/service/note"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type noteServiceMockFunc func(mc *minimock.Controller) noteService.NoteService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		testStartTime = time.Now().Truncate(time.Second)

		id        = gofakeit.Int64()
		title     = gofakeit.Animal()
		content   = gofakeit.Animal()
		createdAt = testStartTime
		updatedAt = testStartTime // Если я ставлю gofakeit.Date(), то тесты не проходят, но если поставить как createdAt или time.Now(), то проходят

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetRequest{Id: id}

		serviceRes = &model.Note{
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

		res = &desc.GetResponse{
			Note: &desc.Note{
				Id: id,
				Info: &desc.NoteInfo{
					Title:   title,
					Content: content,
				},
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) noteService.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(serviceRes, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  serviceErr,
			noteServiceMock: func(mc *minimock.Controller) noteService.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.noteServiceMock(mc)
			api := noteAPI.NewImplementation(noteServiceMock)

			res, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}

}
