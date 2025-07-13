package interceptor

import (
	"context"
	"time"

	"github.com/GolZrd/easy-grpc/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Засекаем время
	now := time.Now()

	// Делаем запрос
	res, err := handler(ctx, req)
	// Если ошибка, то логируем на уровне Error
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}
	//Если запрос прошел успешно, то логируем на уровне Info
	logger.Info("request success", zap.String("method", info.FullMethod), zap.Any("req", req), zap.Any("res", res), zap.Duration("duration", time.Since(now)))

	return res, err
}
