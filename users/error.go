package users

import (
	"errors"
	"fmt"
	"strings"
)

type CustomErr struct {
	Err  error
	Code int
}

func (e *CustomErr) Error() string {
	return fmt.Sprintf("%s,%d", e.Err.Error(), e.Code)
}

func CheckError(err error) error {
	ce := CustomErr{}
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

func Validate(u User) error {
	if u.ID <= 0 {
		return &CustomErr{errors.New("missing id"), 400}
	} else if u.Name == "" {
		return &CustomErr{errors.New("missing name"), 400}
	} else if u.Password == "" {
		return &CustomErr{errors.New("missing password"), 400}
	} else if u.Email == "" {
		return &CustomErr{errors.New("missing email"), 400}
	}
	return nil
}
