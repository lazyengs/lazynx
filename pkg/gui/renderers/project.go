package renderers

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/gantoreno/lazynx/pkg/files"
)

func Project(g *gocui.Gui, v *gocui.View) {
	packageJson := files.GetPackageJson()

	fmt.Fprintln(v, packageJson["name"])
}
