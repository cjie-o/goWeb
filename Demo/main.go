package main

import (
	"embed"
	"fmt"

	goweb "github.com/cjie9759/goWeb"
	"github.com/cjie9759/goWeb/Demo/controller"
)

//go:embed public/*
var FS embed.FS

func main() {
	fmt.Println(goweb.NewApp(&FS).Get(&controller.Index{}).SetMiddle(goweb.MWLog).Run(":8080"))
}
