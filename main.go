package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	correctStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // Red
	targetStyle    = lipgloss.NewStyle().Bold(true).Padding(1, 2)
	resultStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(1, 2)
	inputStyle     = lipgloss.NewStyle().Padding(1, 2)
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).Margin(1)
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
)

var sampleTexts = []string{
	"hello world",
	"go is awesome",
	"bubbletea makes TUI easy",
	"practice makes perfect",
	"fast fingers win races",
}

type model struct {
	text      string    // The target text to type
	input     string    // User's input
	finished  bool      // Whether typing is done
	startTime time.Time // Start time of typing
	isTyping  bool      // Whether the user has started typing
	width     int       // Terminal width
	height    int       // Terminal height
	cursorPos int       // Cursor position
	blink     bool      // Cursor blinking state
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, blinkCursor())
}

func blinkCursor() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return blinkMsg{}
	})
}

type blinkMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.finished && msg.String() == "r" {
			return model{text: sampleTexts[rand.Intn(len(sampleTexts))], width: m.width, height: m.height}, nil
		}
		if !m.isTyping {
			m.startTime = time.Now()
			m.isTyping = true
		}
		if msg.String() == "backspace" {
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.cursorPos--
			}
		} else if msg.String() == "ctrl+backspace" || msg.String() == "cmd+backspace" {
			if len(m.input) > 0 {
				words := strings.Fields(m.input)
				if len(words) > 0 {
					m.input = strings.Join(words[:len(words)-1], " ")
					if len(m.input) > 0 {
						m.input += " "
					}
					m.cursorPos = len(m.input)
				}
			}
		} else if msg.String() == " " {
			if len(m.input) > 0 && m.input[len(m.input)-1] != ' ' {
				index := len(m.input)
				for index < len(m.text) && m.text[index] != ' ' {
					index++
				}
				if index < len(m.text) {
					m.input = m.text[:index+1]
					m.cursorPos = len(m.input)
				}
			}
		} else {
			m.input += msg.String()
			m.cursorPos++
		}
		if m.input == m.text {
			m.finished = true
		}
	case blinkMsg:
		m.blink = !m.blink
		return m, blinkCursor()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	content := ""
	if m.finished {
		content = borderStyle.Render(
			fmt.Sprintf("%s\n%s\n%s\n%s",
				resultStyle.Render("You typed correctly!"),
				fmt.Sprintf("WPM: %d", m.calculateWPM()),
				fmt.Sprintf("Accuracy: %.2f%%", m.calculateAccuracy()),
				"Press 'r' to restart",
			),
		)
	} else {
		styledInput := ""
		for i, r := range m.text {
			if i == m.cursorPos {
				cursor := "|"
				if !m.blink {
					cursor = " "
				}
				styledInput += cursorStyle.Render(cursor)
			}
			if i < len(m.input) {
				if rune(m.input[i]) == r {
					styledInput += correctStyle.Render(string(r))
				} else {
					styledInput += incorrectStyle.Render(string(m.input[i]))
				}
			} else {
				styledInput += string(r)
			}
		}
		content = borderStyle.Render(
			fmt.Sprintf("%s\n%s\n%s",
				targetStyle.Render("Type the following:"),
				targetStyle.Render(m.text),
				inputStyle.Render("Your input: "+styledInput),
			),
		)
	}
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(content)
}

func (m model) calculateWPM() int {
	elapsed := time.Since(m.startTime).Minutes()
	if elapsed == 0 {
		return 0
	}
	wordCount := len(m.text) / 5
	return int(float64(wordCount) / elapsed)
}

func (m model) calculateAccuracy() float64 {
	correctChars := 0
	for i := 0; i < len(m.input) && i < len(m.text); i++ {
		if m.input[i] == m.text[i] {
			correctChars++
		}
	}
	return (float64(correctChars) / float64(len(m.text))) * 100
}

func main() {
	m := model{text: sampleTexts[rand.Intn(len(sampleTexts))]}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}
