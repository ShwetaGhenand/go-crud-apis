package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//User declartion
type User struct {
	ID      int
	Name    string
	Email   string
	Phone   string
	Age     int
	Address string
}

var users = []User{}

func isExist(id int) (bool, int) {
	for index, user := range users {
		if user.ID == id {
			return true, index
		}
	}
	return false, 0
}

//GetHealth : returns status of ther sevice
func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("Get health endpoint called.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Ok!"))
}

//GetUsers : returns list of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all users endpoint called.")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error in encoding object", err)
		http.Error(w, "Error enconding object!", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

//GetUser : returns single user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := isExist(id)
	if found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res, err := json.Marshal(&users[index])
		if err != nil {
			log.Println("Error in encoding response object", err)
			http.Error(w, "Error enconding response object!", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found!"))
	}
}

//AddUser : add single user
func AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Add user endpoint called.")
	user := User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Error in decoding request body", err)
		http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
		return
	}
	found, _ := isExist(user.ID)
	if found {
		log.Println("User alredy exists, duplicate user id!")
		http.Error(w, "User alredy exists, duplicate user id!", http.StatusBadRequest)
		return
	}
	users = append(users, user)
	res, err := json.Marshal(&user)
	if err != nil {
		log.Println("Error in encoding response object", err)
		http.Error(w, "Error enconding response object!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

//UpdateUser : update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := isExist(id)
	if found {
		user := User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println("Error in decoding request body", err)
			http.Error(w, "Error in decoding request body!", http.StatusBadRequest)
			return
		}
		d, _ := isExist(user.ID)
		if d {
			log.Println("User alredy exists, duplicate user id!")
			http.Error(w, "User alredy exists, duplicate user id!", http.StatusBadRequest)
			return
		}
		users[index] = user
		res, err := json.Marshal(&user)
		if err != nil {
			log.Println("Error in encoding response object", err)
			http.Error(w, "Error enconding response object!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found!"))
	}
}

//DeleteUser : delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete user endpoint called.")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	found, index := isExist(id)
	if found {
		users = append(users[:index], users[index+1:]...)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully!"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found!"))
	}
}
