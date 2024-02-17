package server

import (
	"context"
	"database/sql"
	"errors"
)

type service struct {
	db *Queries
}

func reqToModel(req JSONUser) User {
	m := User{
		ID:       int32(req.ID),
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Address:  sql.NullString{String: req.Address, Valid: req.Address != ""},
		Age:      sql.NullInt32{Int32: int32(req.Age), Valid: req.Age != 0},
	}
	return m
}

func modelToRes(user User) JSONUser {
	res := JSONUser{
		ID:       int(user.ID),
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}

	if user.Phone.Valid {
		res.Phone = user.Phone.String
	}

	if user.Address.Valid {
		res.Address = user.Address.String
	}

	if user.Age.Valid {
		res.Age = int(user.Age.Int32)
	}

	return res
}

func (s *service) listUsers() ([]JSONUser, error) {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return nil, checkError(err)
	}
	if len(users) == 0 {
		return nil, &customErr{errors.New("users not found"), 404}
	}
	var res []JSONUser
	for _, user := range users {
		res = append(res, modelToRes(user))
	}
	return res, nil
}

func (s *service) getUser(id int32) (JSONUser, error) {
	user, err := s.db.GetUser(context.Background(), id)
	dto := modelToRes(user)
	if err != nil {
		return dto, checkError(err)
	}
	return dto, nil
}

func (s *service) createUser(req JSONUser) error {
	arg := CreateUserParams(reqToModel(req))
	err := s.db.CreateUser(context.Background(), arg)
	if err != nil {
		return checkError(err)
	}
	return nil
}

func (s *service) updateUser(req JSONUser) error {
	m := reqToModel(req)
	arg := UpdateUserParams{m.Name, m.Email, m.Phone, m.Age, m.Address, m.ID}
	err := s.db.UpdateUser(context.Background(), arg)
	if err != nil {
		return checkError(err)
	}
	return nil
}

func (s *service) deleteUser(id int32) error {
	err := s.db.DeleteUser(context.Background(), id)
	if err != nil {
		return checkError(err)
	}
	return nil
}

func (s *service) UserExists(name, password string) error {
	arg := UserExistsParams{name, password}
	_, err := s.db.UserExists(context.Background(), arg)
	if err != nil {
		return checkError(err)
	}
	return nil
}
