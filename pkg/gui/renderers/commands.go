package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func Commands(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintln(v, "(d) Commands")
}
