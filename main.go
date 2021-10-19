package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	user "github.com/go-crud-apis/handler"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

func main() {
	var port int
	r := mux.NewRouter()
	r.HandleFunc("/health", user.GetHealth).Methods("GET")
	s := r.PathPrefix("/users").Subrouter()
	s.HandleFunc("", user.AddUser).Methods("POST")
	s.HandleFunc("", user.GetUsers).Methods("GET")
	s.HandleFunc("/{id}", user.GetUser).Methods("GET")
	s.HandleFunc("/{id}", user.UpdateUser).Methods("PUT")
	s.HandleFunc("/{id}", user.DeleteUser).Methods("DELETE")

	app := &cli.App{
		Name: "user-apis",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "users rest apis",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Value:       8081,
						Usage:       "port to listen server",
						Destination: &port,
					},
				},
				Action: func(c *cli.Context) error {
					log.Println("Server is listening on port : ", port)
					if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
						log.Println("Error in creating server", err)
						return err
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
