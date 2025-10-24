package main

import (
	"ciao-admin/cmd/admin/app"
	"fmt"
)

var (
	version     = "dev"
	gitShortSHA = "unknown"
)

func main() {
	app := app.AdminApplication{
		Version: fmt.Sprintf("%s-%s", version, gitShortSHA),
	}
	if app.Init() {
		app.Run()
	}
}
