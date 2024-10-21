package tui

import (
	"fmt"
	"gotube/internal/videx"
	"net/url"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	query   string
	results []videx.Video // parsed results
	cursor  int           // cursor pos
	state   string        // state btwn searching and displaying
	width   int           // term with
	height  int           // term height
	loading bool          // mpv delay
	spinner spinner.Model
}

func Initialize() model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	return model{
		query:   "",
		results: []videx.Video{},
		cursor:  0,
		state:   "searching", // default state
		width:   80,          // default
		height:  24,          // default
		loading: false,
		spinner: s,
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

					m.loading = true
					return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
						cmd := exec.Command("nohup", "mpv", videoURL, ">/dev/null", "2>&1", "&")
						err := cmd.Start()
						if err != nil {
							fmt.Println("[ERROR] mpv failed to open:", err)
						}
						return "mpv_loaded"
					})
				}

			}
		case "backspace":
			if m.state == "searching" && len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}
		case "k": // vim keys
			if m.state == "displaying" && m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.state == "displaying" && m.cursor < len(m.results)-1 {
				m.cursor++
			}
		default:
			// append the rest of the chars to the query
			if m.state == "searching" && len(msg.String()) == 1 {
				m.query += msg.String()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case spinner.TickMsg:
		if m.loading {
			newModel, cmd := m.spinner.Update(msg)
			m.spinner = newModel
			return m, cmd
		}
	case string:
		if msg == "mpv_loaded" {
			m.loading = false
		}
	}
	return m, nil
}

const pageSize = 10 // top 10 results

func (m model) View() string {
	// Styles P
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))
	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF"))
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF00FF")).
		Foreground(lipgloss.Color("#000000"))
	cursorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF00FF"))
	searchBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(50).
		Align(lipgloss.Center)

	blink := time.Now().Second()%2 == 0 // every second
	cursor := "|"
	if !blink {
		cursor = " "
	}

	searchBox := searchBoxStyle.Render(fmt.Sprintf("Search: %s%s\n\n", m.query, cursor))

	// Pagination logic: display only the top 10 results
	start := 0
	if m.cursor >= pageSize {
		start = m.cursor - pageSize + 1
	}
	end := start + pageSize
	if end > len(m.results) {
		end = len(m.results)
	}

	if m.loading {
		return m.spinner.View() + "\nLoading video..."
	}

	// the UI
	// if we got something, display it
	var resultsView string
	if len(m.results) > 0 {
		for i := start; i < end; i++ {
			cursor := " " // cursor icon for the selected vid
			if m.cursor == i {
				cursor = "> "
			}
			videoTitle := m.results[i].Title
			if m.cursor == i {
				videoTitle = selectedStyle.Render(videoTitle)
			}
			resultsView += fmt.Sprintf("%s%s\n%s | Duration: %s\n\n",
				cursorStyle.Render(cursor), titleStyle.Render(videoTitle),
				urlStyle.Render(m.results[i].URL), m.results[i].Length)
		}
	} else {
		resultsView = "No results found."
	}

	// Combine the search box and the results view, center it
	ui := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		Render(searchBox + "\n\n" + resultsView)

	return ui
}

func Run() {
	p := tea.NewProgram(Initialize())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	} else {
		fmt.Println("Program exited succesfully")
	}
}
