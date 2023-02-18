package util

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

func GetMarkdowdRender() *glamour.TermRenderer {
	styleConfig := glamour.LightStyleConfig
	if termenv.HasDarkBackground() {
		styleConfig = glamour.DarkStyleConfig
	}
	styleConfig.Document.Margin = nil
	render, _ := glamour.NewTermRenderer(glamour.WithStyles(styleConfig))
	return render
}

func RenderedMarkdown(md string) string {
	res, _ := GetMarkdowdRender().Render(md)
	return res
}

func Blanks(n int) string {
	return strings.Repeat(" ", n)
}
