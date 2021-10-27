package server

import (
	"errors"
	"fmt"
	"strings"
)

type customErr struct {
	Err  error
	Code int
}

func (e *customErr) Error() string {
	return fmt.Sprintf("%s,%d", e.Err.Error(), e.Code)
}

func checkError(err error) error {
	ce := customErr{}
	switch e := err.Error(); {
	case strings.Contains(e, "no rows"):
		ce.Err = errors.New("user not found")
		ce.Code = 404
	case strings.Contains(e, "users_pkey"):
		ce.Err = errors.New("duplicate user id")
		ce.Code = 400
	default:
		ce.Err = errors.New("database error")
		ce.Code = 500
	}
	return &ce
}

func validate(u userDto) error {
	if u.ID <= 0 {
		return &customErr{errors.New("missing id"), 400}
	} else if u.Name == "" {
		return &customErr{errors.New("missing name"), 400}
	} else if u.Email == "" {
		return &customErr{errors.New("missing email"), 400}
	}
	return nil
}
