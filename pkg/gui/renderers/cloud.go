package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func Cloud(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintln(v, "(c) NX Cloud")
}
