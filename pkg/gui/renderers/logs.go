package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func Logs(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintln(v, "(e) Logs")
}
