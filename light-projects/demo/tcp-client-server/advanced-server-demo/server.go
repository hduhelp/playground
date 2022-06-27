package main

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

var p = tea.NewProgram(initialModel())

func init() {
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

			p.Send(connection{
				remoteAddr: conn.RemoteAddr().String(),
				conn:       conn,
			})
		}
	}()
}

type connection struct {
	remoteAddr string
	conn       net.Conn
}

func main() {
	if err := p.Start(); err != nil {
		log.Fatal("boot tui failed:", err)
	}
}
