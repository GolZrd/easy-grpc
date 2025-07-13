package tests

import (
	"context"
	"fmt"
	"testing"

	noteAPI "github.com/GolZrd/easy-grpc/internal/api/note"
	"github.com/GolZrd/easy-grpc/internal/logger"
	"github.com/GolZrd/easy-grpc/internal/model"
	serviceMocks "github.com/GolZrd/easy-grpc/internal/service/mocks"
	noteService "github.com/GolZrd/easy-grpc/internal/service/note"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	// Инициализируем логгер
	logger.Init(zap.NewNop().Core())

	t.Parallel()

	type noteServiceMockFunc func(mc *minimock.Controller) noteService.NoteService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	// Данные для теста
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Animal()
		content = gofakeit.Animal()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.NoteInfo{
				Title:   title,
				Content: content,
			},
		}

		info = &model.NoteInfo{
			Title:   title,
			Content: content,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "success case",
			// Входные параметры в этом кейсе
			args: args{
				ctx: ctx,
				req: req,
			},
			// Ожидаемый результат
			want: res,
			// В этом успешном кейсе не должно быть ошибки
			err: nil,
			noteServiceMock: func(mc *minimock.Controller) noteService.NoteService {
				// Создаем экземпляр (реализацию) mock-объекта(сервиса)
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			// Входные параметры в этом кейсе
			args: args{
				ctx: ctx,
				req: req,
			},
			// Ожидаемый результат
			want: nil,
			// В этом кейсе должна вернуться ошибка
			err: serviceErr,
			noteServiceMock: func(mc *minimock.Controller) noteService.NoteService {
				// Создаем экземпляр (реализацию) mock-объекта(сервиса)
				mock := serviceMocks.NewNoteServiceMock(mc)
				// В этом кейсе ожидаем ошибку
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.noteServiceMock(mc)
			// Создаем структуру newImplementation, которую определили в пакете api
			api := noteAPI.NewImplementation(noteServiceMock)

			// Вызываем метод у объекта api и передаем в него тестовые аргументы
			res, err := api.Create(tt.args.ctx, tt.args.req)
			// Проверяем, что полученные данные соответствует ожидаемым
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
