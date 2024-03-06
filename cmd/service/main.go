package main

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/service"
)

var (
	Version   string // TODO: use git tags
	BuildTime string
	Commit    string
)

func main() {
	fmt.Println(Commit)
	fmt.Println(BuildTime)
	service.Run()
}
