package pkg

import (
	"fmt"
	"go-chatting-cli/api"
	"io"
	"sync"
)

type Connection struct {
	conn api.Chat_ChatServer
	send chan *api.ChatMessage
	quit chan struct{}
}

// コントラクタ関数
func NewConnection(conn api.Chat_ChatServer) *Connection {
	c := &Connection{
		conn: conn,
		send: make(chan *api.ChatMessage),
		quit: make(chan struct{}),
	}
	go c.start()
	return c
}

func (c *Connection) Close() error {
	close(c.quit)
	close(c.send)
	return nil
}

func (c *Connection) Send(msg *api.ChatMessage) {
	defer func() {
		recover()
	}()
	c.send <- msg
}

func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			c.conn.Send(msg)
		case <-c.quit:
			running = false
		}
	}
}

// メッセージ取得
func (c *Connection) GetMessage(broadcast chan<- *api.ChatMessage) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}

		go func(msg *api.ChatMessage) {
			fmt.Println("block?")
			select {
			case broadcast <- msg:
			case <-c.quit:
			}
		}(msg)
	}
}

// ChatServerの状態を格納
type ChatServer struct {
	api.UnimplementedChatServer
	broadcast   chan *api.ChatMessage
	quit        chan struct{}
	connections []*Connection
	connLock    sync.Mutex
}

func NewChatServer() *ChatServer {
	srv := &ChatServer{
		broadcast: make(chan *api.ChatMessage),
		quit:      make(chan struct{}),
	}

	go srv.start()
	return srv
}

func (c *ChatServer) Close() error {
	close(c.quit)
	return nil
}

func (c *ChatServer) start() {
	running := true
	for running {
		select {
		case msg := <-c.broadcast:
			c.connLock.Lock()
			for _, v := range c.connections {
				go v.Send(msg)
			}
			c.connLock.Unlock()
		case <-c.quit:
			running = false
		}
	}
}

func (c *ChatServer) Chat(stream api.Chat_ChatServer) error {
	// コネクション作成
	conn := NewConnection(stream)

	// コネクションの排他制御
	c.connLock.Lock()
	c.connections = append(c.connections, conn)
	c.connLock.Unlock()

	err := conn.GetMessage(c.broadcast)

	c.connLock.Lock()
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}
	c.connLock.Unlock()

	return err
}
