package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func writeError(w http.ResponseWriter, err error) {
	log.Println("Error in sending response", err)
	http.Error(w, "Error in sending response", http.StatusInternalServerError)
}

// GetHealth : returns status of the service
func getHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("Get health endpoint called.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Status: Ok!")); err != nil {
		writeError(w, err)
	}
}

// GetUsers : returns list of users
func (s *server) getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all users endpoint called.")
	users, dbErr := s.userRepo.getUsers()
	if dbErr != nil {
		http.Error(w, dbErr.Message, dbErr.Code)
		return
	}
	if users == nil {
		http.Error(w, "Users not found!", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error in encoding object", err)
		http.Error(w, "Error encoding object!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetUser : returns single user by id
func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	u, dbErr := s.userRepo.getUser(id)
	if dbErr != nil {
		http.Error(w, dbErr.Message, dbErr.Code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(u)
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
func (s *server) addUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Add user endpoint called.")
	userBody := user{}
	if err := json.NewDecoder(r.Body).Decode(&userBody); err != nil {
		log.Println("Error in decoding request body", err)
		http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
		return
	}
	if err := validate(userBody); err != nil {
		log.Println("Validatio failed", err)
		http.Error(w, err.Message, err.Code)
		return
	}
	u, dbErr := s.userRepo.addUser(userBody)
	if dbErr != nil {
		http.Error(w, dbErr.Message, dbErr.Code)
		return
	}
	res, err := json.Marshal(u)
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
func (s *server) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userBody := user{}
	if err := json.NewDecoder(r.Body).Decode(&userBody); err != nil {
		log.Println("Error in decoding request body", err)
		http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
		return
	}
	updatedUser, dbErr := s.userRepo.updateUser(id, userBody)
	if dbErr != nil {
		http.Error(w, dbErr.Message, dbErr.Code)
		return
	}
	res, err := json.Marshal(updatedUser)
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
func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if dbErr := s.userRepo.deleteUser(id); dbErr != nil {
		http.Error(w, dbErr.Message, dbErr.Code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User deleted successfully!")); err != nil {
		writeError(w, err)
		return
	}
}
