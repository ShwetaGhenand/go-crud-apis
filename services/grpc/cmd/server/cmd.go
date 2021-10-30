package server

import (
	"log"

	"github.com/sebnyberg/flagtags"
	"github.com/urfave/cli/v2"
)

// server declartion

// Cmd : command to start the server
func Cmd() *cli.Command {
	var conf config
	return &cli.Command{
		Name:  "server",
		Usage: "users rest apis",
		Flags: flagtags.MustParseFlags(&conf),
		Action: func(c *cli.Context) error {
			if err := conf.DBConfig.validate(); err != nil {
				log.Fatalf("invalid configuration %v", err)
			}
			return run(&conf)
		},
	}
}
