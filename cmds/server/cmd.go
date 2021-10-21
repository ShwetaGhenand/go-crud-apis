package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

// user declartion
type user struct {
	ID      int
	Name    string
	Email   string
	Phone   string
	Age     int
	Address string
}

// server declartion
type server struct {
	router *mux.Router
	users  []user
}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
		users:  make([]user, 0),
	}
	s.router.HandleFunc("/health", GetHealth).Methods("GET")
	sr := s.router.PathPrefix("/users").Subrouter()
	sr.HandleFunc("", s.AddUser).Methods("POST")
	sr.HandleFunc("", s.GetUsers).Methods("GET")
	sr.HandleFunc("/{id}", s.GetUser).Methods("GET")
	sr.HandleFunc("/{id}", s.UpdateUser).Methods("PUT")
	sr.HandleFunc("/{id}", s.DeleteUser).Methods("DELETE")
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Just call the internal router here
	s.router.ServeHTTP(w, r)
}

// Cmd : command to start the server
func Cmd() *cli.Command {
	var port int
	return &cli.Command{
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
			s := newServer()
			log.Println("Server is listening on port : ", port)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s); err != nil {
				log.Fatalln(err)
			}
			return nil
		},
	}
}
