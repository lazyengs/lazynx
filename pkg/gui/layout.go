package gui

import (
	"log"
	"math"

	"github.com/awesome-gocui/gocui"
	"github.com/gantoreno/lazynx/pkg/gui/renderers"
)

type Panel struct {
	name     string
	title    string
	area     rune
	frame    bool
	renderer func(g *gocui.Gui, v *gocui.View)
}

var panels = []Panel{
	{name: "project", title: "Project", area: 'a', frame: true, renderer: renderers.Project},
	{name: "overview", title: "Overview", area: 'b', frame: true, renderer: renderers.Overview},
	{name: "cloud", title: "NX Cloud :: CI Optimizations", area: 'c', frame: true, renderer: renderers.Cloud},
	{name: "commands", title: "Common NX Commands", area: 'd', frame: true, renderer: renderers.Commands},
	{name: "logs", title: "Logs", area: 'e', frame: true, renderer: renderers.Logs},
	{name: "footer", title: "Footer", area: 'f', frame: false, renderer: renderers.Footer},
}

func getAreaDimensions(g *gocui.Gui, area rune) (int, int, int, int) {
	maxX, maxY := g.Size()

	leftHalfWidth := int(math.Floor(float64(maxX) / 3))
	rightPanelThirdHeight := int(math.Floor((float64(maxY) - 4) / 3))

	switch area {
	case 'a':
		return 0, 0, leftHalfWidth, 2
	case 'b':
		_, _, _, y1, _ := g.ViewPosition("project")
		return 0, y1 + 1, leftHalfWidth, y1 + rightPanelThirdHeight
	case 'c':
		_, _, _, y1, _ := g.ViewPosition("overview")
		return 0, y1 + 1, leftHalfWidth, y1 + rightPanelThirdHeight
	case 'd':
		_, _, _, y1, _ := g.ViewPosition("cloud")
		return 0, y1 + 1, leftHalfWidth, maxY - 2
	case 'e':
		return leftHalfWidth + 1, 0, maxX - 1, maxY - 2
	case 'f':
		return -1, maxY - 2, maxX + 1, maxY
	}

	return 0, 0, 0, 0

}

func LayoutManager(g *gocui.Gui) error {
	for _, panel := range panels {
		x0, y0, x1, y1 := getAreaDimensions(g, panel.area)

		if v, err := g.SetView(panel.name, x0, y0, x1, y1, 0); err != nil {
			if err != gocui.ErrUnknownView {
				log.Panicln(err)
			}

			v.Title = panel.title
			v.Frame = panel.frame
			v.FrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

			panel.renderer(g, v)
		}
	}

	return nil
}
