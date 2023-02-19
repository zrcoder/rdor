package print

import (
	"fmt"

	"github.com/zrcoder/rdor/pkg/style"
)

func Errorln(err error) {
	fmt.Println(style.Error.Render(err.Error()))
}
