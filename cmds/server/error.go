package server

import (
	"strings"
)

type custonError struct {
	Message string
	Code    int
}

func (e *custonError) Error(msg string, code int) *custonError {
	return e
}

func checkError(err error) *custonError {
	dbErr := custonError{}
	switch e := err.Error(); {
	case strings.Contains(e, "no rows"):
		dbErr.Message = "User not found!"
		dbErr.Code = 404
	case strings.Contains(e, "users_pkey"):
		dbErr.Message = "Duplicate user id!"
		dbErr.Code = 400
	default:
		dbErr.Message = "Database error!"
		dbErr.Code = 500
	}
	return &dbErr
}

func validate(u user) *custonError {
	if u.ID <= 0 {
		return &custonError{"Missing id!", 400}
	} else if u.Name == "" {
		return &custonError{"Missing name!", 400}
	} else if u.Email == "" {
		return &custonError{"Missing email!", 400}
	}
	return nil
}
