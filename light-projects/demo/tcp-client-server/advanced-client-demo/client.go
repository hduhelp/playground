package main

import (
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var p = tea.NewProgram(initialModel())

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	data := make([]byte, 1024)
	go func() {
		for {
			n, err := conn.Read(data)
			if err != nil {
				log.Println(err)
			}
			if n != 0 {
				p.Send(msgPack{
					Conn: conn,
					addr: "127.0.0.1:23333",
					data: data[:n],
				})
			}
		}
	}()

	if err := p.Start(); err != nil {
		log.Fatal("boot tui failed:", err)
	}
}

type msgPack struct {
	net.Conn

	addr string
	data []byte
}
