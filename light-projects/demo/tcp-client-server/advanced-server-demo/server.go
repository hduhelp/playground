package main

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

var p = tea.NewProgram(initialModel())

func main() {
	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:23333")
		if err != nil {
			log.Fatal("boot server failed:", err)
		}
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("accept connection failed:", err)
			}

			handle(conn)
		}
	}()

	if err := p.Start(); err != nil {
		log.Fatal("boot tui failed:", err)
	}
}

type msgPack struct {
	net.Conn

	addr string // 发送者的地址
	data []byte
}

func handle(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			p.Send(err)
		}
	}()

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		if n != 0 {
			p.Send(msgPack{
				Conn: conn,
				addr: conn.RemoteAddr().String(),
				data: buf[:n],
			})
		}
	}
}
