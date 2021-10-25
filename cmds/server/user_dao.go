package server

import (
	"database/sql"
	"log"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) getUsers() ([]user, *custonError) {
	query := `select * FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error getting user", err)
		return nil, checkError(err)
	}
	defer rows.Close()
	var users []user
	for rows.Next() {
		var e user
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Phone, &e.Age, &e.Address); err != nil {
			log.Println("Error getting user", err)
			return nil, checkError(err)
		}
		users = append(users, e)
	}
	if err := rows.Close(); err != nil {
		log.Println("Error getting user", err)
		return nil, checkError(err)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error getting user", err)
		return nil, checkError(err)
	}
	return users, nil
}

func (r *userRepository) getUser(id int) (user, *custonError) {
	query := `select * FROM users where id = $1`
	row := r.db.QueryRow(query, id)
	var u user
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Age, &u.Address); err != nil {
		log.Println("Error getting user", err)
		return u, checkError(err)
	}
	return u, nil
}

func (r *userRepository) addUser(u user) (user, *custonError) {
	query := `insert into users 
	values($1, $2, $3, $4, $5, $6) 
	returning id, name, email, phone, age, address`
	row := r.db.QueryRow(query, u.ID, u.Name, u.Email, u.Phone, u.Age, u.Address)
	var e user
	if err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Phone, &e.Age, &e.Address); err != nil {
		log.Println("Error inserting user", err)
		return e, checkError(err)
	}
	return e, nil
}

func (r *userRepository) updateUser(id int, e user) (user, *custonError) {
	query := `UPDATE users SET
	name = $1,
	email = $2,
	phone = $3,
	age = $4,
	address = $5
	WHERE id = $6
	returning id, name, email, phone, age, address`
	row := r.db.QueryRow(query, e.Name, e.Email, e.Phone, e.Age, e.Address, id)
	var ue user
	if err := row.Scan(&ue.ID, &ue.Name, &ue.Email, &ue.Phone, &ue.Age, &ue.Address); err != nil {
		log.Println("Error updating user", err)
		return ue, checkError(err)
	}
	return ue, nil
}

func (r *userRepository) deleteUser(id int) *custonError {
	query := `delete from users
	where id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting user", err)
		return checkError(err)
	}
	return nil
}
