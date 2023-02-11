package main

import (
	"os"

	"github.com/zrcoder/hanoi/internal"
)

func main() {
	internal.Run(os.Args[1:])
}
