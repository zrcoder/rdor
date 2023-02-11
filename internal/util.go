package internal

import (
	"fmt"
	"strings"
)

func printError(err error) {
	fmt.Println(errorStyle.Render(err.Error()))
}

func blanks(n int) string {
	return strings.Repeat(" ", n)
}
