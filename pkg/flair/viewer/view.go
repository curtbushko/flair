package viewer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Slot names for base24 palette display.
var slotNames = [24]string{
	"base00", "base01", "base02", "base03",
	"base04", "base05", "base06", "base07",
	"base08", "base09", "base0A", "base0B",
	"base0C", "base0D", "base0E", "base0F",
	"base10", "base11", "base12", "base13",
	"base14", "base15", "base16", "base17",
}

// View implements tea.Model and renders the current page.
func (m Model) View() string {
	var content string

	switch m.currentPage {
	case PageSelector:
		content = m.renderSelector()
	case PagePalette:
		content = m.renderPalette()
	case PageTokens:
		content = m.renderTokens()
	case PageComponents:
		content = m.renderComponents()
	}

	// Add help footer.
	help := m.renderHelp()
	return content + "\n" + help
}

// renderSelector renders the theme selection list.
func (m Model) renderSelector() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Theme Selector"))
	b.WriteString("\n\n")

	for i, theme := range m.themes {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		selected := ""
		if theme == m.selectedTheme {
			selected = " [selected]"
		}

		line := fmt.Sprintf("%s%s%s", cursor, theme, selected)
		if i == m.cursor {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("14")).
				Bold(true).
				Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// renderPalette renders the base24 palette colors.
func (m Model) renderPalette() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Palette Colors"))
	b.WriteString("\n\n")

	// Render in a grid: 4 columns x 6 rows.
	for row := 0; row < 6; row++ {
		for col := 0; col < 4; col++ {
			idx := row*4 + col
			if idx >= 24 {
				break
			}

			name := slotNames[idx]
			hex := m.palette.Colors[idx]
			if hex == "" {
				hex = "#000000"
			}

			// Color swatch using background.
			swatch := lipgloss.NewStyle().
				Background(lipgloss.Color(hex)).
				Padding(0, 2).
				Render("  ")

			// Label with hex value.
			label := fmt.Sprintf("%s %s", name, hex)

			b.WriteString(swatch)
			b.WriteString(" ")
			b.WriteString(label)

			if col < 3 {
				b.WriteString("   ")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderTokens renders semantic tokens grouped by category.
func (m Model) renderTokens() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		MarginTop(1)

	b.WriteString(titleStyle.Render("Semantic Tokens"))
	b.WriteString("\n\n")

	// Surface tokens.
	b.WriteString(sectionStyle.Render("Surface"))
	b.WriteString("\n")
	for name, hex := range m.tokens.Surface {
		b.WriteString(m.renderTokenLine(name, hex))
	}

	// Text tokens.
	b.WriteString(sectionStyle.Render("Text"))
	b.WriteString("\n")
	for name, hex := range m.tokens.Text {
		b.WriteString(m.renderTokenLine(name, hex))
	}

	// Status tokens.
	b.WriteString(sectionStyle.Render("Status"))
	b.WriteString("\n")
	for name, hex := range m.tokens.Status {
		b.WriteString(m.renderTokenLine(name, hex))
	}

	// Syntax tokens.
	if len(m.tokens.Syntax) > 0 {
		b.WriteString(sectionStyle.Render("Syntax"))
		b.WriteString("\n")
		for name, hex := range m.tokens.Syntax {
			b.WriteString(m.renderTokenLine(name, hex))
		}
	}

	// Diff tokens.
	if len(m.tokens.Diff) > 0 {
		b.WriteString(sectionStyle.Render("Diff"))
		b.WriteString("\n")
		for name, hex := range m.tokens.Diff {
			b.WriteString(m.renderTokenLine(name, hex))
		}
	}

	return b.String()
}

// renderTokenLine renders a single token with color swatch.
func (m Model) renderTokenLine(name, hex string) string {
	swatch := lipgloss.NewStyle().
		Background(lipgloss.Color(hex)).
		Padding(0, 1).
		Render(" ")

	return fmt.Sprintf("  %s %s %s\n", swatch, name, hex)
}

// renderComponents renders styled component examples.
func (m Model) renderComponents() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Component Showcase"))
	b.WriteString("\n\n")

	// Status buttons using token colors.
	b.WriteString("Status Buttons:\n")
	for name, hex := range m.tokens.Status {
		btn := lipgloss.NewStyle().
			Background(lipgloss.Color(hex)).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 2).
			Bold(true).
			Render(name)

		b.WriteString("  ")
		b.WriteString(btn)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Text examples using text tokens.
	if len(m.tokens.Text) > 0 {
		b.WriteString("Text Styles:\n")
		for name, hex := range m.tokens.Text {
			txt := lipgloss.NewStyle().
				Foreground(lipgloss.Color(hex)).
				Render(name + ": Example text")

			b.WriteString("  ")
			b.WriteString(txt)
			b.WriteString("\n")
		}
	}

	// Syntax highlighting examples.
	if len(m.tokens.Syntax) > 0 {
		b.WriteString("\nSyntax Highlighting:\n")
		for name, hex := range m.tokens.Syntax {
			txt := lipgloss.NewStyle().
				Foreground(lipgloss.Color(hex)).
				Render(name)

			b.WriteString("  ")
			b.WriteString(txt)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderHelp renders the help footer.
func (m Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var hints []string
	hints = append(hints, "Tab: switch page")
	hints = append(hints, "j/k: navigate")

	if m.currentPage == PageSelector {
		hints = append(hints, "Enter: select theme")
	}

	hints = append(hints, "q: quit")

	return helpStyle.Render(strings.Join(hints, " | "))
}
