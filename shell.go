package clipy

import (
	_ "embed"
	"fmt"
)

//go:embed shell/clipy.zsh
var zsh string

//go:embed shell/clipy.bash
var bash string

func GenerateShellScript(sh string) {
	switch sh {
	case "zsh":
		fmt.Println(zsh)

	case "bash":
		fmt.Println(bash)
	}
}
