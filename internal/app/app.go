package app

import (
	"context"
	"log"
	"net"

	"github.com/GolZrd/easy-grpc/internal/closer"
	"github.com/GolZrd/easy-grpc/internal/config"
	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

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

func (a *App) InitServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) InitGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

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
