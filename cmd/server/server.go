package main

import (
	"fmt"
	"go-chatting-cli/api"
	"go-chatting-cli/pkg"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	server := grpc.NewServer()
	chatSrv := pkg.NewChatServer()

	api.RegisterChatServer(server, chatSrv)
	fmt.Println("Start server :8080")
	log.Fatal(server.Serve(lst))
}
