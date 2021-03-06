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
		ce.Err = fmt.Errorf("database error %w", err)
		ce.Code = 500
	}
	return &ce
}

func validate(user JSONUser) error {
	if user.ID <= 0 {
		return &customErr{errors.New("missing id"), 400}
	} else if user.Name == "" {
		return &customErr{errors.New("missing name"), 400}
	} else if user.Password == "" {
		return &customErr{errors.New("missing password"), 400}
	} else if user.Email == "" {
		return &customErr{errors.New("missing email"), 400}
	}
	return nil
}
