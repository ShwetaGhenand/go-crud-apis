package main

import (
	"log"
	"os"

	"github.com/go-crud-apis/cmds/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "user-apis",
		Commands: []*cli.Command{
			server.Cmd(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
