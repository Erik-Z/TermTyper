package cmd

type Theme struct {
	Name  string
	Color string
}

var AvailableThemes = []Theme{
	{Name: "Magenta", Color: "#FF00FF"},
	{Name: "Red", Color: "#FF0000"},
	{Name: "Green", Color: "#00FF00"},
	{Name: "Blue", Color: "#0000FF"},
	{Name: "Yellow", Color: "#FFFF00"},
	{Name: "Cyan", Color: "#00FFFF"},
	{Name: "White", Color: "#FFFFFF"},
	{Name: "Orange", Color: "#FFA500"},
	{Name: "Pink", Color: "#FFC0CB"},
	{Name: "Lime", Color: "#32CD32"},
}

func GetThemeColor(themeName string) string {
	for _, theme := range AvailableThemes {
		if theme.Name == themeName {
			return theme.Color
		}
	}
	return "#FF00FF" // Default to magenta
}

func GetThemeIndex(themeName string) int {
	for i, theme := range AvailableThemes {
		if theme.Name == themeName {
			return i
		}
	}
	return 0
}
