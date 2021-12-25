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
	// バリデーション
	ok := pkg.ArgsValidate()
	if !ok {
		return
	}

	// context作成
	ctx := context.Background()
	// grpx接続ダイアル設定
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

	// メッセージ受付処理
	go pkg.StreamRecv(stream)
	// コネクション確立
	pkg.ConnectEstablish(stream)
}
