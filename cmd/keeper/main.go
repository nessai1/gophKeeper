package main

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper"
)

func main() {
	err := keeper.Listen()
	if err != nil {
		fmt.Printf("Error while work keeper app: %s", err.Error())
	}
}
