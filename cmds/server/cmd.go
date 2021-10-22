package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

// server declartion
type server struct {
	router   *mux.Router
	userRepo *userRepository
}

func newServer(url string) *server {
	db := initDB(url)
	s := &server{
		router:   mux.NewRouter(),
		userRepo: &userRepository{db: db},
	}
	s.router.HandleFunc("/health", getHealth).Methods("GET")
	sr := s.router.PathPrefix("/users").Subrouter()
	sr.HandleFunc("", s.addUser).Methods("POST")
	sr.HandleFunc("", s.getUsers).Methods("GET")
	sr.HandleFunc("/{id}", s.getUser).Methods("GET")
	sr.HandleFunc("/{id}", s.updateUser).Methods("PUT")
	sr.HandleFunc("/{id}", s.deleteUser).Methods("DELETE")
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Just call the internal router here
	s.router.ServeHTTP(w, r)
}

// Cmd : command to start the server
func Cmd() *cli.Command {
	var port int
	var url string
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
			&cli.StringFlag{
				Name:        "url",
				EnvVars:     []string{"DATABASE_URL"},
				Destination: &url,
			},
		},
		Action: func(c *cli.Context) error {
			s := newServer(url)
			log.Println("Server is listening on port : ", port)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s); err != nil {
				log.Fatalln(err)
			}
			return nil
		},
	}
}
