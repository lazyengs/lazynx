package gui

import "github.com/awesome-gocui/gocui"

func quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()

	return gocui.ErrQuit
}
