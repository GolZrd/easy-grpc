package main

import (
	"context"
	"log"
	"time"

	desc "github.com/GolZrd/easy-grpc/pkg/note_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	//Конектимся к серверу
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	defer conn.Close()

	client := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetNote()))
}
