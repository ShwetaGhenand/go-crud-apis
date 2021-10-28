package server

import (
	"fmt"
	"log"
	"net/http"

	u "github.com/go-crud-apis/users"
	"github.com/go-crud-apis/users/auth"
	c "github.com/go-crud-apis/users/config"
	"github.com/gorilla/mux"
	"github.com/sebnyberg/flagtags"
	"github.com/urfave/cli/v2"
)

// server declartion
type server struct {
	router  *mux.Router
	service *u.Service
	secret  string
}

func newServer(conf c.Config) *server {
	db, err := c.InitDB(conf.DBConfig)
	if err != nil {
		log.Fatalf("Database Error %v", err)
	}
	s := &server{
		router:  mux.NewRouter(),
		service: u.NewService(db),
		secret:  conf.Secret,
	}
	s.router.HandleFunc("/health", getHealth).Methods("GET")
	s.router.HandleFunc("/signin", s.createUser).Methods("POST")
	s.router.HandleFunc("/login", s.loginUser).Methods("POST")

	sr := s.router.PathPrefix("/users").Subrouter()
	sr.Use(auth.JWTMiddleware(conf.Secret))
	sr.HandleFunc("", s.listUsers).Methods("GET")
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
	var conf c.Config
	return &cli.Command{
		Name:  "server",
		Usage: "users rest apis",
		Flags: flagtags.MustParseFlags(&conf),
		Action: func(c *cli.Context) error {
			s := newServer(conf)
			log.Println("Server is listening on port : ", conf.Port)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), s); err != nil {
				log.Fatalln(err)
			}
			return nil
		},
	}
}
