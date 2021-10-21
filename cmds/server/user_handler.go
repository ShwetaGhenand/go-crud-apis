package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// userExists : to check user exists or not if exist, return found flag and index
func (s *server) userExists(id int) (f bool, i int) {
	for index, user := range s.users {
		if user.ID == id {
			return true, index
		}
	}
	return false, 0
}

func writeError(w http.ResponseWriter, err error) {
	log.Println("Error in sending response", err)
	http.Error(w, "Error in sending response", http.StatusInternalServerError)
}

// GetHealth : returns status of the service
func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("Get health endpoint called.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Status: Ok!")); err != nil {
		writeError(w, err)
	}
}

// GetUsers : returns list of users
func (s *server) GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all users endpoint called.")
	if err := json.NewEncoder(w).Encode(s.users); err != nil {
		log.Println("Error in encoding object", err)
		http.Error(w, "Error encoding object!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetUser : returns single user by id
func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := s.userExists(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("User not found!")); err != nil {
			writeError(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(s.users[index])
	if err != nil {
		log.Println("Error in encoding response object", err)
		http.Error(w, "Error encoding response object!", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(res); err != nil {
		writeError(w, err)
		return
	}
}

// AddUser : add single user
func (s *server) AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Add user endpoint called.")
	u := user{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println("Error in decoding request body", err)
		http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
		return
	}
	found, _ := s.userExists(u.ID)
	if found {
		log.Println("User alredy exists, duplicate user id!")
		http.Error(w, "User alredy exists, duplicate user id!", http.StatusBadRequest)
		return
	}
	s.users = append(s.users, u)
	res, err := json.Marshal(&u)
	if err != nil {
		log.Println("Error in encoding response object", err)
		http.Error(w, "Error encoding response object!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(res); err != nil {
		writeError(w, err)
		return
	}
}

// UpdateUser : update user
func (s *server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := s.userExists(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("User not found!")); err != nil {
			writeError(w, err)
		}
		return
	}
	u := user{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println("Error in decoding request body", err)
		http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
		return
	}
	if u.ID != 0 {
		http.Error(w, "Id can not be updated!", http.StatusBadRequest)
		return
	}
	u.ID = s.users[index].ID
	s.users[index] = u
	res, err := json.Marshal(&u)
	if err != nil {
		log.Println("Error in encoding response object", err)
		http.Error(w, "Error encoding response object!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		writeError(w, err)
		return
	}
}

// DeleteUser : delete user
func (s *server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := s.userExists(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("User not found!")); err != nil {
			writeError(w, err)
		}
		return
	}
	s.users = append(s.users[:index], s.users[index+1:]...)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User deleted successfully!")); err != nil {
		writeError(w, err)
		return
	}
}
