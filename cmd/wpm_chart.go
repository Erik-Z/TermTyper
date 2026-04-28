package cmd

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type WPMChartBubble struct {
	data       []float64
	width      int
	height     int
	style      lipgloss.Style
	showGrid   bool
	showLabels bool
	showYAxis  bool
}

func NewWPMChartBubble(width, height int) *WPMChartBubble {
	return &WPMChartBubble{
		data:       []float64{},
		width:      width,
		height:     height,
		style:      lipgloss.NewStyle().Foreground(lipgloss.Color("46")),
		showGrid:   true,
		showLabels: true,
		showYAxis:  true,
	}
}

func (wc *WPMChartBubble) UpdateData(wpmData []float64) {
	wc.data = wpmData
}

func (wc *WPMChartBubble) View() string {
	if len(wc.data) == 0 {
		return wc.style.Render("No WPM data available")
	}

	var result strings.Builder

	// Title
	title := wc.style.Copy().Bold(true).Render("WPM Progress Chart")
	result.WriteString(title + "\n")

	result.WriteString(strings.Repeat("─", wc.width+8) + "\n") // +8 for Y-axis labels

	// Find min/max for scaling
	min, max := wc.data[0], wc.data[0]
	for _, wpm := range wc.data {
		if wpm < min {
			min = wpm
		}
		if wpm > max {
			max = wpm
		}
	}

	if max == min {
		max = min + 1
	}

	for y := wc.height - 1; y >= 0; y-- {
		// Y-axis label
		if wc.showYAxis {
			yValue := min + (float64(y)/float64(wc.height-1))*(max-min)
			yLabel := fmt.Sprintf("%5.0f", yValue)
			result.WriteString(yLabel + " │")
		} else {
			result.WriteString("     │")
		}

		// Chart content
		for x := 0; x < wc.width; x++ {
			if x < len(wc.data) {
				// Scale the WPM value to chart height
				normalizedY := (wc.data[x] - min) / (max - min)
				chartY := int(normalizedY * float64(wc.height-1))

				if y == chartY {
					result.WriteString(wc.style.Render("●")) // Data point
				} else if y < chartY {
					if wc.showGrid {
						result.WriteString(wc.style.Faint(true).Render("│")) // Grid line
					} else {
						result.WriteString(" ")
					}
				} else {
					result.WriteString(" ")
				}
			} else {
				result.WriteString(" ")
			}
		}
		result.WriteString("│\n")
	}

	result.WriteString("     └" + strings.Repeat("─", wc.width) + "┘\n")

	if wc.showLabels && len(wc.data) > 1 {
		result.WriteString("     0s")
		result.WriteString(strings.Repeat(" ", wc.width-4))
		result.WriteString(fmt.Sprintf("%ds\n", len(wc.data)-1))
	}

	return result.String()
}

func (wc *WPMChartBubble) averageWPM() float64 {
	if len(wc.data) == 0 {
		return 0
	}

	sum := 0.0
	for _, wpm := range wc.data {
		sum += wpm
	}
	return sum / float64(len(wc.data))
}

func (wc *WPMChartBubble) SetStyle(style lipgloss.Style) {
	wc.style = style
}

func (wc *WPMChartBubble) ToggleGrid() {
	wc.showGrid = !wc.showGrid
}

func (wc *WPMChartBubble) ToggleLabels() {
	wc.showLabels = !wc.showLabels
}

func (wc *WPMChartBubble) ToggleYAxis() {
	wc.showYAxis = !wc.showYAxis
}
