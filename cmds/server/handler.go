package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func writeError(w http.ResponseWriter, err error) {
	log.Printf("Error occurred : %v", err)
	switch err.(type) {
	case *customErr:
		s := strings.Split(err.Error(), ",")
		message := s[0]
		code, e := strconv.Atoi(s[1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, message, code)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getHealth : returns status of the service
func getHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("Get health endpoint called.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Status: Ok!")); err != nil {
		writeError(w, err)
	}
}

// listUsers : returns list of users
func (s *server) listUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all users endpoint called.")
	dtos, err := s.service.listUsers()
	if err != nil {
		writeError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// getUser : returns single user by id
func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single user endpoint called.")
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	dto, err := s.service.getUser(int32(id))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(dto)
	if err != nil {
		writeError(w, err)
		return
	}
	if _, err := w.Write(res); err != nil {
		writeError(w, err)
		return
	}
}

// addUser : add single user
func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Add user endpoint called.")
	dtoReq := userDto{}
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		writeError(w, err)
		return
	}
	if err := validate(dtoReq); err != nil {
		writeError(w, err)
		return
	}
	dtoRes, err := s.service.createUser(dtoReq)
	if err != nil {
		writeError(w, err)
		return
	}
	res, err := json.Marshal(dtoRes)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(res); err != nil {
		writeError(w, err)
		return
	}
}

// updateUser : update user
func (s *server) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user endpoint called.")
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)

	dtoReq := userDto{}
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		writeError(w, err)
		return
	}
	dtoReq.ID = int(id)
	user, err := s.service.updateUser(dtoReq)
	if err != nil {
		writeError(w, err)
		return
	}
	dtoRes, err := json.Marshal(user)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(dtoRes); err != nil {
		writeError(w, err)
		return
	}
}

// deleteUser : delete user
func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete user endpoint called.")
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	if err := s.service.deleteUser(int32(id)); err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User deleted successfully!")); err != nil {
		writeError(w, err)
		return
	}
}
