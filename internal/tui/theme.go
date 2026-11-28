package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Keep the same layout and meaning on terminals without box-drawing glyphs.
func asciiArt(view string) string {
	return strings.NewReplacer("─", "-", "━", "-", "│", "|", "┌", "+", "┐", "+", "└", "+", "┘", "+", "├", "+", "┤", "+", "╭", "+", "╮", "+", "╰", "+", "╯", "+", "┬", "+", "┴", "+", "┼", "+", "█", "#", "░", ".", "✓", "+", "✦", "*", "◈", "*", "◆", "*", "◇", "o", "○", "o", "◐", "~", "▸", ">", "→", ">", "←", "<", "↑", "^", "↓", "v", "…", ".", "·", ".", "›", ">", "–", "-", "—", "-", "’", "'", "‘", "'", "“", "\"", "”", "\"").Replace(view)
}

var (
	violet = lipgloss.AdaptiveColor{Light: "#6042A6", Dark: "#BBA7F5"}
	teal   = lipgloss.AdaptiveColor{Light: "#146B66", Dark: "#72DACC"}
	amber  = lipgloss.AdaptiveColor{Light: "#855014", Dark: "#F1BE78"}
	green  = lipgloss.AdaptiveColor{Light: "#256941", Dark: "#8DDCA9"}
	muted  = lipgloss.AdaptiveColor{Light: "#526175", Dark: "#AAB9CC"}
	border = lipgloss.AdaptiveColor{Light: "#8593AC", Dark: "#526583"}
	active = lipgloss.AdaptiveColor{Light: "#E8DDF8", Dark: "#352956"}
)

func (m model) paint(color lipgloss.TerminalColor, label string, bold bool) string {
	if !m.useColor {
		return label
	}
	return lipgloss.NewStyle().Foreground(color).Bold(bold).Render(label)
}

func (m model) activeTab(label string) string {
	if !m.useColor {
		return "[" + label + "]"
	}
	return lipgloss.NewStyle().Foreground(violet).Background(active).Bold(true).Render(" " + label + " ")
}

func (m model) selectedRow(label string, width int) string {
	if !m.useColor {
		return label
	}
	return lipgloss.NewStyle().Foreground(teal).Background(active).Bold(true).Width(width).Render(label)
}
