package main

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/gantoreno/lazynx/pkg/gui"
	"github.com/gantoreno/lazynx/pkg/hooks"
)

func main() {
	hooks.EnsureNxProject()

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(gui.LayoutManager)

	gui.AddKeybindings(g)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
