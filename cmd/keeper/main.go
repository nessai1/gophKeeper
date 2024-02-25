package main

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper"
)

func main() {
	if err := keeper.Listen(); err != nil {
		fmt.Printf("\033[31mKeeper was crashed: %s\033[0m", err.Error())
	}
}
