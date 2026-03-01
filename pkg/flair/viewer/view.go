package viewer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model and renders the 2-panel layout.
func (m Model) View() string {
	// Render left panel (theme list).
	leftPanel := m.renderThemeList()

	// Render right panel (content page).
	var rightPanel string
	switch m.currentPage {
	case PageTextStatus:
		rightPanel = m.renderTextStatus()
	case PageInteractive:
		rightPanel = m.renderInteractive()
	case PageDataDisplay:
		rightPanel = m.renderDataDisplay()
	}

	// Calculate panel widths.
	leftWidth := 25
	rightWidth := m.width - leftWidth - 3 // 3 for border/padding
	if rightWidth < 40 {
		rightWidth = 40
	}

	// Style the panels.
	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		PaddingRight(1)

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		PaddingLeft(2)

	// Join panels horizontally.
	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftPanel),
		rightStyle.Render(rightPanel),
	)

	// Add help footer.
	help := m.renderHelp()
	return layout + "\n" + help
}

// renderThemeList renders the theme list for the left panel.
func (m Model) renderThemeList() string {
	var b strings.Builder

	titleStyle := m.titleStyle()
	b.WriteString(titleStyle.Render("Styles"))
	b.WriteString("\n\n")

	for i, theme := range m.themes {
		// Build prefix: cursor indicator + selection indicator
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		selected := " "
		if i == m.selectedIndex {
			selected = "*"
		}

		prefix := cursor + selected + " "
		line := prefix + theme

		// Style based on state
		if i == m.cursor && i == m.selectedIndex {
			// Both cursor and selected
			line = m.selectedCursorStyle().Render(line)
		} else if i == m.cursor {
			// Just cursor (preview)
			line = m.cursorStyle().Render(line)
		} else if i == m.selectedIndex {
			// Just selected
			line = m.selectedStyle().Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// renderTextStatus renders the Text & Status page with realistic content.
func (m Model) renderTextStatus() string {
	var b strings.Builder

	titleStyle := m.titleStyle()
	sectionStyle := m.sectionStyle()

	b.WriteString(titleStyle.Render("Text & Status"))
	b.WriteString("\n\n")

	// Text styles section.
	b.WriteString(sectionStyle.Render("Text Styles"))
	b.WriteString("\n\n")

	// Primary text.
	primaryHex := m.getTextColor("text.primary", "#c0caf5")
	primaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primaryHex))
	b.WriteString(primaryStyle.Render("Primary text: Lorem ipsum dolor sit amet, consectetur adipiscing elit."))
	b.WriteString("\n\n")

	// Secondary text.
	secondaryHex := m.getTextColor("text.secondary", "#a9b1d6")
	secondaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(secondaryHex))
	b.WriteString(secondaryStyle.Render("Secondary text: Sed do eiusmod tempor incididunt ut labore."))
	b.WriteString("\n\n")

	// Muted text.
	mutedHex := m.getTextColor("text.muted", "#565f89")
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(mutedHex))
	b.WriteString(mutedStyle.Render("Muted text: Ut enim ad minim veniam, quis nostrud exercitation."))
	b.WriteString("\n\n")

	// Status messages section.
	b.WriteString(sectionStyle.Render("Status Messages"))
	b.WriteString("\n\n")

	// Error message.
	errorHex := m.getStatusColor("status.error", "#f7768e")
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(errorHex)).Bold(true)
	b.WriteString(errorStyle.Render("Error: Failed to connect to the server. Please check your network."))
	b.WriteString("\n")

	// Warning message.
	warningHex := m.getStatusColor("status.warning", "#e0af68")
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(warningHex)).Bold(true)
	b.WriteString(warningStyle.Render("Warning: Your session will expire in 5 minutes."))
	b.WriteString("\n")

	// Success message.
	successHex := m.getStatusColor("status.success", "#9ece6a")
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(successHex)).Bold(true)
	b.WriteString(successStyle.Render("Success: File uploaded successfully."))
	b.WriteString("\n")

	// Info message.
	infoHex := m.getStatusColor("status.info", "#7dcfff")
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(infoHex)).Bold(true)
	b.WriteString(infoStyle.Render("Info: Press Ctrl+C to cancel the operation."))
	b.WriteString("\n")

	return b.String()
}

// renderInteractive renders the Interactive Components page.
func (m Model) renderInteractive() string {
	var b strings.Builder

	titleStyle := m.titleStyle()
	sectionStyle := m.sectionStyle()

	b.WriteString(titleStyle.Render("Interactive Components"))
	b.WriteString("\n\n")

	// Buttons section.
	b.WriteString(sectionStyle.Render("Buttons"))
	b.WriteString("\n\n")

	raisedHex := m.getSurfaceColor("surface.background.raised", "#24283b")
	primaryHex := m.getTextColor("text.primary", "#c0caf5")
	mutedHex := m.getTextColor("text.muted", "#565f89")
	accentHex := m.getAccentColor("#7aa2f7")
	inverseHex := m.getTextColor("text.inverse", "#1a1b26")

	// Default button.
	buttonStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(raisedHex)).
		Foreground(lipgloss.Color(primaryHex)).
		Padding(0, 2)
	b.WriteString("  ")
	b.WriteString(buttonStyle.Render("Submit"))
	b.WriteString("  ")

	// Focused button.
	buttonFocusedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(accentHex)).
		Foreground(lipgloss.Color(inverseHex)).
		Padding(0, 2).
		Bold(true)
	b.WriteString(buttonFocusedStyle.Render("Cancel"))
	b.WriteString("  ")

	// Disabled button.
	buttonDisabledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(mutedHex)).
		Padding(0, 2)
	b.WriteString(buttonDisabledStyle.Render("Disabled"))
	b.WriteString("\n\n")

	// Input fields section.
	b.WriteString(sectionStyle.Render("Input Fields"))
	b.WriteString("\n\n")

	sunkenHex := m.getSurfaceColor("surface.background.sunken", "#16161e")
	borderHex := m.getBorderColor("#565f89")
	borderFocusHex := m.getBorderFocusColor("#7aa2f7")

	// Default input.
	inputStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(sunkenHex)).
		Foreground(lipgloss.Color(primaryHex)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderHex)).
		Padding(0, 1).
		Width(46)
	b.WriteString(inputStyle.Render("Enter your name..."))
	b.WriteString("\n\n")

	// Focused input.
	inputFocusedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(sunkenHex)).
		Foreground(lipgloss.Color(primaryHex)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderFocusHex)).
		Padding(0, 1).
		Width(46)
	b.WriteString(inputFocusedStyle.Render("john@example.com"))
	b.WriteString(" (focused)")
	b.WriteString("\n\n")

	// Selection list section.
	b.WriteString(sectionStyle.Render("Selection List"))
	b.WriteString("\n\n")

	selectionHex := m.getSurfaceColor("surface.background.selection", "#33467c")

	listItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primaryHex))
	listSelectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(selectionHex)).
		Foreground(lipgloss.Color(accentHex)).
		Bold(true)

	items := []string{"Option A", "Option B", "Option C", "Option D"}
	for i, item := range items {
		prefix := "  "
		if i == 1 {
			// Selected item.
			prefix = "> "
			b.WriteString(listSelectedStyle.Render(prefix + item))
		} else {
			b.WriteString(listItemStyle.Render(prefix + item))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderDataDisplay renders the Data Display page.
func (m Model) renderDataDisplay() string {
	var b strings.Builder

	titleStyle := m.titleStyle()
	sectionStyle := m.sectionStyle()

	b.WriteString(titleStyle.Render("Data Display"))
	b.WriteString("\n\n")

	primaryHex := m.getTextColor("text.primary", "#c0caf5")
	secondaryHex := m.getTextColor("text.secondary", "#a9b1d6")
	mutedHex := m.getTextColor("text.muted", "#565f89")
	successHex := m.getStatusColor("status.success", "#9ece6a")
	warningHex := m.getStatusColor("status.warning", "#e0af68")
	errorHex := m.getStatusColor("status.error", "#f7768e")

	// Table section.
	b.WriteString(sectionStyle.Render("Table"))
	b.WriteString("\n\n")

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(secondaryHex)).
		Bold(true).
		Width(15)
	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primaryHex)).
		Width(15)

	// Table header.
	b.WriteString("  ")
	b.WriteString(headerStyle.Render("Name"))
	b.WriteString(headerStyle.Render("Status"))
	b.WriteString(headerStyle.Render("Progress"))
	b.WriteString("\n")

	// Separator.
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(mutedHex))
	b.WriteString("  ")
	b.WriteString(sepStyle.Render(strings.Repeat("-", 45)))
	b.WriteString("\n")

	// Table rows.
	tableData := []struct {
		name      string
		status    string
		statusHex string
		progress  string
	}{
		{"Build", "Complete", successHex, "100%"},
		{"Test", "Running", warningHex, "67%"},
		{"Deploy", "Failed", errorHex, "0%"},
	}

	for _, row := range tableData {
		b.WriteString("  ")
		b.WriteString(cellStyle.Render(row.name))
		statusCellStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(row.statusHex)).
			Width(15)
		b.WriteString(statusCellStyle.Render(row.status))
		b.WriteString(cellStyle.Render(row.progress))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Dialog section.
	b.WriteString(sectionStyle.Render("Dialog"))
	b.WriteString("\n\n")

	overlayHex := m.getSurfaceColor("surface.background.overlay", "#24283b")
	borderHex := m.getBorderColor("#565f89")

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(overlayHex)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderHex)).
		Padding(1, 2).
		Width(46)

	dialogContent := "Are you sure you want to continue?\n\nThis action cannot be undone."
	dialogContentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primaryHex))
	b.WriteString(dialogStyle.Render(dialogContentStyle.Render(dialogContent)))
	b.WriteString("\n\n")

	// Code block section.
	b.WriteString(sectionStyle.Render("Code Block"))
	b.WriteString("\n\n")

	keywordHex := m.getSyntaxColor("syntax.keyword", "#bb9af7")
	stringHex := m.getSyntaxColor("syntax.string", "#9ece6a")
	functionHex := m.getSyntaxColor("syntax.function", "#7aa2f7")
	commentHex := m.getSyntaxColor("syntax.comment", "#565f89")

	codeBlockStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.getSurfaceColor("surface.background.sunken", "#16161e"))).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderHex)).
		Padding(1, 2).
		Width(46)

	// Build syntax-highlighted code.
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(keywordHex))
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(stringHex))
	funcStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(functionHex))
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(commentHex))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primaryHex))

	var codeBuf strings.Builder
	codeBuf.WriteString(commentStyle.Render("// Sample function"))
	codeBuf.WriteString("\n")
	codeBuf.WriteString(keywordStyle.Render("func "))
	codeBuf.WriteString(funcStyle.Render("greet"))
	codeBuf.WriteString(normalStyle.Render("(name "))
	codeBuf.WriteString(keywordStyle.Render("string"))
	codeBuf.WriteString(normalStyle.Render(") {"))
	codeBuf.WriteString("\n")
	codeBuf.WriteString(normalStyle.Render("    "))
	codeBuf.WriteString(funcStyle.Render("println"))
	codeBuf.WriteString(normalStyle.Render("("))
	codeBuf.WriteString(stringStyle.Render("\"Hello, \""))
	codeBuf.WriteString(normalStyle.Render(" + name)"))
	codeBuf.WriteString("\n")
	codeBuf.WriteString(normalStyle.Render("}"))

	b.WriteString(codeBlockStyle.Render(codeBuf.String()))
	b.WriteString("\n")

	return b.String()
}

// renderHelp renders the help footer.
func (m Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	hints := []string{
		"↑/↓/j/k: navigate",
		"Enter: select",
		"Tab: switch page",
		"q/Esc: quit",
	}

	return helpStyle.Render(strings.Join(hints, " | "))
}

// Helper methods for styling.

func (m Model) titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)
}

func (m Model) sectionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11"))
}

func (m Model) cursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)
}

func (m Model) selectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")) // Green for selected
}

func (m Model) selectedCursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Green for selected
		Bold(true)
}

// Color helper methods that use TokenData.

func (m Model) getTextColor(key, fallback string) string {
	if hex, ok := m.tokens.Text[key]; ok {
		return hex
	}
	return fallback
}

func (m Model) getStatusColor(key, fallback string) string {
	if hex, ok := m.tokens.Status[key]; ok {
		return hex
	}
	return fallback
}

func (m Model) getSurfaceColor(key, fallback string) string {
	if hex, ok := m.tokens.Surface[key]; ok {
		return hex
	}
	return fallback
}

func (m Model) getSyntaxColor(key, fallback string) string {
	if hex, ok := m.tokens.Syntax[key]; ok {
		return hex
	}
	return fallback
}

func (m Model) getAccentColor(fallback string) string {
	// Try to find accent color in tokens or use fallback.
	return fallback
}

func (m Model) getBorderColor(fallback string) string {
	return fallback
}

func (m Model) getBorderFocusColor(fallback string) string {
	return fallback
}
