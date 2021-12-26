package pkg

import (
	"bufio"
	"fmt"
	"go-chatting-cli/api"
	"io"
	"log"
	"os"
	"time"
)

var waitc = make(chan struct{})

// 引数のバリデーションチェック
func ArgsValidate() bool {
	if len(os.Args) != 3 {
		fmt.Println("第一引数: URL, 第二引数: ユーザー名が必要です。")
		return false
	}
	return true
}

// メッセージを受信・表示
func StreamRecv(stream api.Chat_ChatClient) {
	for {
		timeNow := time.Now().Format("2006-01-02 15:04")
		msg, err := stream.Recv()
		if err == io.EOF {
			close(waitc)
			return
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println(timeNow + ": " + msg.User + ": " + msg.Message)
	}
}

// コネクションが確立しました。
func ConnectEstablish(stream api.Chat_ChatClient) {
	fmt.Println("コネクションが確立しました。" +
		"\"quit\"を押下するか\"ctrl+c\"にてプログラムを停止できます。")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "quit" {
			err := stream.CloseSend()
			if err != nil {
				panic(err)
			}
			break
		}
		err := stream.Send(&api.ChatMessage{
			User:    os.Args[2],
			Message: msg,
		})
		if err != nil {
			panic(err)
		}
	}
	<-waitc
}
