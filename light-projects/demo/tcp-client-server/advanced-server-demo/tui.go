package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	sendMsg, recvMsg tcpMsg
	connections      []connection
	historyMsg       []string
	textInput        textinput.Model
}

type tcpMsg struct {
	prefix string
	buf    []byte
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "input whatever you want in this line"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{
		sendMsg:     tcpMsg{prefix: "HOST", buf: make([]byte, 1024)},
		recvMsg:     tcpMsg{prefix: "", buf: make([]byte, 1024)},
		connections: []connection{},
		historyMsg:  []string{},
		textInput:   ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// TODO: 处理多余的空位 刷新收到的消息
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case connection:
		m.connections = append(m.connections, msg)
		n, err := msg.conn.Read(m.recvMsg.buf)
		if err != nil {
			log.Println("read data from connection failed:", err)
		}

		if n != 0 {
			m.recvMsg.prefix = msg.remoteAddr
			m.recvMsg.buf = m.recvMsg.buf[:n]
		}

		// 广播消息
		for _, cn := range m.connections {
			if _, err := cn.conn.Write(m.recvMsg.buf); err != nil {
				log.Println("failed to boardcast msg:", err)
			}
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			copy(m.sendMsg.buf, []byte(m.textInput.Value()))
			m.textInput.Reset()

			// 广播消息
			for _, cn := range m.connections {
				if _, err := cn.conn.Write(m.sendMsg.buf); err != nil {
					log.Println("failed to send msg:", err)
				}
			}
			return m, nil

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	content := "TcpServer powered by sslime336 with BubbleTea\n\n"
	var b strings.Builder
	for _, msg := range m.historyMsg {
		b.WriteString(msg)
		b.WriteString("\n")
	}
	b.WriteString(m.textInput.View())
	content += b.String()
	return content
}
