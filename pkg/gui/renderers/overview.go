package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func Overview(g *gocui.Gui, v *gocui.View) {
	// packageJson := files.GetPackageJson()
	// nxJson := files.GetNxJson()

	fmt.Fprintln(v, "(b) Overview")
	// fmt.Fprintln(v, icons.Icons["package"], packageJson["name"])
}
