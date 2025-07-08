package note

import (
	"context"

	"github.com/GolZrd/easy-grpc/internal/model"
)

func (s *service) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.NoteRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		// Имитируем ошибку
		// return errors.New("test error")

		// для примера сразу вызовем get
		// _, errTx = s.NoteRepository.Get(ctx, id)
		// if errTx != nil {
		// 	return errTx
		// }

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
