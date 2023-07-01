package main

import (
	"fmt"
	"os"

	"github.com/zrcoder/rdor/internal"
)

//go:generate go run ./internal/maze/levels/tool

func main() {
	if err := internal.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
