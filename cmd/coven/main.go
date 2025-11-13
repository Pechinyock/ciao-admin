package main

import (
	"ciao-admin/cmd/coven/app"
	"fmt"
)

var (
	version     = "dev"
	gitShortSHA = "unknown"
)

func main() {
	app := app.CovenApplication{
		Version: fmt.Sprintf("%s-%s", version, gitShortSHA),
	}
	if app.Init() {
		app.Run()
	}
}
