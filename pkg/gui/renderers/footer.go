package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func Footer(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintln(v, "PgUp/PgDn: scroll, b: view bulk commands, q: quit, x: menu, ← → ↑ ↓: navigate")
}
