package main

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper"
)

func main() {
	if err := keeper.Listen(); err != nil {
		fmt.Printf("Keeper was crashed: %s", err.Error())
	}
}
