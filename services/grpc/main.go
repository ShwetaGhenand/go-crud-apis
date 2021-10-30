package main

import (
	"log"
	"os"

	server "github.com/go-crud-apis/services/grpc/cmd/server"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "user-apis",
		Commands: []*cli.Command{
			server.Cmd(),
		},
	}
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
