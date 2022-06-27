package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		log.Fatalf("start listen failed: %s", err)
	}
	defer conn.Close()

	sendMsg, revMsg := make([]byte, 1024), make([]byte, 1024)
	go func() {
		for {
			n, err := conn.Read(revMsg)
			if err != nil {
				log.Fatalf("accept failed: %s", err.Error())
			}
			fmt.Printf("\nFROM %s: %s\n%s> ", conn.RemoteAddr().String(), string(revMsg[:n]), conn.LocalAddr().String())
		}
	}()

	for {
		fmt.Printf("%s> ", conn.LocalAddr().String())
		if _, err := fmt.Scanln(&sendMsg); err != nil {
			log.Printf("scan err: %s", err.Error())
		}
		conn.Write(sendMsg)
	}
}
