// Example: charm-bubbletea demonstrates using flair with bubbletea for TUI applications.
//
// This example shows a simple progress indicator that uses themed styles
// for spinners, progress bars, and text.
//
// Run with: go run ./examples/charm-bubbletea
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	flairbubble "github.com/curtbushko/flair/pkg/charm/bubbletea"
	"github.com/curtbushko/flair/pkg/flair"
)

// Spinner frames for the loading animation.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type model struct {
	styles   *flairbubble.Styles
	progress float64
	frame    int
	done     bool

	// Custom styles built from theme colors.
	titleStyle    lipgloss.Style
	spinnerStyle  lipgloss.Style
	labelStyle    lipgloss.Style
	progressFull  lipgloss.Style
	progressEmpty lipgloss.Style
	successStyle  lipgloss.Style
	helpStyle     lipgloss.Style
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tickMsg:
		m.frame = (m.frame + 1) % len(spinnerFrames)
		m.progress += 0.02
		if m.progress >= 1.0 {
			m.progress = 1.0
			m.done = true
			return m, tea.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.titleStyle.Render("Flair + Bubbletea Example"))
	b.WriteString("\n\n")

	if m.done {
		b.WriteString(m.successStyle.Render("  Done!"))
		b.WriteString("\n\n")
	} else {
		// Spinner
		spinner := m.spinnerStyle.Render(spinnerFrames[m.frame])
		label := m.labelStyle.Render(" Loading...")
		b.WriteString("  " + spinner + label + "\n\n")
	}

	// Progress bar
	barWidth := 40
	filled := int(float64(barWidth) * m.progress)
	empty := barWidth - filled

	b.WriteString("  ")
	b.WriteString(m.progressFull.Render(strings.Repeat("█", filled)))
	b.WriteString(m.progressEmpty.Render(strings.Repeat("░", empty)))
	b.WriteString(" ")
	b.WriteString(m.labelStyle.Render(fmt.Sprintf("%.0f%%", m.progress*100)))
	b.WriteString("\n\n")

	// Help text
	b.WriteString(m.helpStyle.Render("  Press q to quit"))
	b.WriteString("\n")

	return tea.NewView(b.String())
}

func main() {
	// Load theme and create styles.
	theme := flair.MustLoad()
	styles := flairbubble.NewStyles(theme)

	// Build custom styles using the base styles from flair.
	m := model{
		styles:        styles,
		titleStyle:    styles.Text.Primary.Bold(true).MarginBottom(1),
		spinnerStyle:  styles.Accent.Primary,
		labelStyle:    styles.Text.Primary,
		progressFull:  styles.Status.Success,
		progressEmpty: styles.Text.Muted,
		successStyle:  styles.Status.Success.Bold(true),
		helpStyle:     styles.Text.Muted,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
