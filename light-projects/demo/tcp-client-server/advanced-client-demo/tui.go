package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	historyMsg []string
	textInput  textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "You can input whatever you want (only support ASCII chars currently)"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30
	return model{
		historyMsg: []string{},
		textInput:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		panic(msg)
	case msgPack:
		m.historyMsg = append(m.historyMsg, fmt.Sprintf("FROM %s: %s", msg.addr, string(msg.data)))

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if _, err := conn.Write([]byte(m.textInput.Value())); err != nil {
				m.textInput.Reset()
				return m, func() tea.Msg {
					return err
				}
			}
			m.textInput.Reset()
			return m, nil

		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	b := "TcpClient powered by sslime336 with BubbleTea\n\n"
	// TODO: 修复服务器断线后疯狂刷新空消息的bug
	for _, msg := range m.historyMsg {
		b += msg + "\n"
	}
	b += m.textInput.View()
	return b
}
