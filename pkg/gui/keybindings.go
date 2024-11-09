package gui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type Keybinding struct {
	viewname string
	key      interface{}
	mod      gocui.Modifier
	handler  func(g *gocui.Gui, v *gocui.View) error
}

var keybindings = []Keybinding{
	{viewname: "", key: 'q', mod: gocui.ModNone, handler: quit},
	{viewname: "", key: gocui.KeyEsc, mod: gocui.ModNone, handler: quit},
}

func AddKeybindings(g *gocui.Gui) {
	for _, keybinding := range keybindings {
		if err := g.SetKeybinding(keybinding.viewname, keybinding.key, keybinding.mod, keybinding.handler); err != nil {
			log.Panicln(err)
		}
	}
}
