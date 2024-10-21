package tui

import (
	"fmt"
	"gotube/internal/videx"
	"net/url"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	query   string
	results []videx.Video
	cursor  int
	state   string
}

func Initialize() model {
	return model{
		query:   "",
		results: []videx.Video{},
		cursor:  0,
		state:   "searching", // default state
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update gets called on every every event (i.e. keypress)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == "searching" { // search mode
				// Construct the YouTube search URL from the query
				searchURL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(m.query)
				results, err := videx.ExtractVideos(searchURL)
				if err != nil {
					fmt.Println("[ERROR] fetch failed", err)
					return m, nil
				}
				m.results = results // update results view
				m.cursor = 0
				m.state = "displaying"
			} else if m.state == "displaying" {
				if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
					video := m.results[m.cursor]
					videoURL := "https://www.youtube.com" + video.URL
					cmd := exec.Command("nohup", "mpv", videoURL, ">/dev/null", "2>&1", "&")
					err := cmd.Start()
					if err != nil {
						fmt.Println("[ERROR] mpv failed to open:", err)
					}
					fmt.Printf("Playing video: %s\n", video.Title)
				}

			}
		case "k": // vim keys
			if m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
		case "backspace":
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}
		default:
			// append the rest of the chars to the query
			if len(msg.String()) == 1 {
				m.query += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))
	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF"))
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF00FF"))

	// let's build the UI
	s := fmt.Sprintf("Search: %s\n\n", m.query)

	// if we got something, display it
	if len(m.results) > 0 {
		s += "Search results:\n\n"
		for i, video := range m.results {
			cursor := " " // cursor icon for the selected vid
			if m.cursor == i {
				cursor = "> "
			}
			videoTitle := video.Title
			if m.cursor == i {
				videoTitle = selectedStyle.Render(videoTitle)
			}
			s += fmt.Sprintf("%s %s\n%s | Duration: %s\n\n", cursor, titleStyle.Render(videoTitle), urlStyle.Render(video.URL), video.Length)
		}
	} else {
		s += "No results found\n"
	}
	return s
}

func Run() {
	p := tea.NewProgram(Initialize())
	if finalModel, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	} else {
		fmt.Println("Program exited succesfully", finalModel)
	}
}
