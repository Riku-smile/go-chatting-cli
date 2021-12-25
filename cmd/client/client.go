package main

import (
	"context"
	"go-chatting-cli/api"
	"go-chatting-cli/pkg"
	"log"
	"os"

	"google.golang.org/grpc"
)

func init() {
	log.SetPrefix("Client: ")
}

func main() {
	ok := pkg.ArgsValidate()
	if !ok {
		return
	}

	ctx := context.Background()

	conn, err := grpc.Dial(os.Args[1], grpc.WithInsecure())

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := api.NewChatClient(conn)
	stream, err := c.Chat(ctx)
	if err != nil {
		panic(err)
	}

	go pkg.StreamRecv(stream)
	pkg.ConnectEstablish(stream)
}
