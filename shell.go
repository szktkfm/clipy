package clipy

import (
	_ "embed"
	"fmt"
)

//go:embed clipy.zsh
var zsh string

func Shell(sh string) {
	fmt.Println(zsh)
}
