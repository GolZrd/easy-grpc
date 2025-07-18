package app

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"github.com/GolZrd/easy-grpc/internal/closer"
	"github.com/GolZrd/easy-grpc/internal/config"
	"github.com/GolZrd/easy-grpc/internal/interceptor"
	"github.com/GolZrd/easy-grpc/internal/logger"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Уровнь логирования задаем при старте приложения, можем в консоле указать уровень
var logLevel = flag.String("l", "info", "log level")

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// Создаем объект нашей структуры App
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.RunGRPCServer()
}

// инициализируем зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.InitConfig,
		a.InitLogger,
		a.InitServiceProvider,
		a.InitGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) InitConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) InitLogger(_ context.Context) error {
	flag.Parse()
	// Указываем уровень логирования
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	stdout := zapcore.AddSync(os.Stdout)
	// Настраиваем запись в файл
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     3, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, atomicLevel),
		zapcore.NewCore(fileEncoder, file, atomicLevel),
	)

	logger.Init(core)

	return nil
}

func (a *App) InitServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) InitGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		// Добавили интерцептор, также был добавлен еще один interceptor для логгирования
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(interceptor.LogInterceptor, interceptor.ValidateInterceptor)))

	reflection.Register(a.grpcServer)

	desc.RegisterNoteV1Server(a.grpcServer, a.serviceProvider.NoteImpl(ctx))
	return nil
}

func (a *App) RunGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}
	return nil
}
