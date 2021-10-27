package server

import (
	"context"
	"database/sql"
	"errors"
)

type service struct {
	db *Queries
}

func createDtoToModel(dto userDto) User {
	m := User{}
	m.ID = int32(dto.ID)
	m.Name = dto.Name
	m.Password = dto.Password
	m.Email = dto.Email
	if dto.Phone != "" {
		m.Phone = sql.NullString{String: dto.Phone, Valid: true}
	} else {
		m.Phone = sql.NullString{String: dto.Phone, Valid: false}
	}
	if dto.Address != "" {
		m.Address = sql.NullString{String: dto.Address, Valid: true}
	} else {
		m.Address = sql.NullString{String: dto.Address, Valid: false}
	}
	if dto.Age != 0 {
		m.Age = sql.NullInt32{Int32: int32(dto.Age), Valid: true}
	} else {
		m.Age = sql.NullInt32{Int32: int32(dto.Age), Valid: false}
	}
	return m
}

func modelToDto(user User) userDto {
	dto := userDto{}
	if user.ID != 0 {
		dto.ID = int(user.ID)
	}
	if user.Name != "" {
		dto.Name = user.Name
	}
	if user.Password != "" {
		dto.Password = user.Password
	}
	if user.Email != "" {
		dto.Email = user.Email
	}
	if user.Phone.Valid {
		dto.Phone = user.Phone.String
	}
	if user.Address.Valid {
		dto.Address = user.Address.String
	}
	if user.Age.Valid {
		dto.Age = int(user.Age.Int32)
	}
	return dto
}

func (s *service) listUsers() ([]userDto, error) {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return nil, checkError(err)
	}
	if len(users) == 0 {
		return nil, &customErr{errors.New("users not found"), 404}
	}
	var dtos []userDto
	for _, user := range users {
		dtos = append(dtos, modelToDto(user))
	}
	return dtos, nil
}

func (s *service) getUser(id int32) (userDto, error) {
	user, err := s.db.GetUser(context.Background(), id)
	dto := modelToDto(user)
	if err != nil {
		return dto, checkError(err)
	}
	return dto, nil
}

func (s *service) createUser(dto userDto) (userDto, error) {
	arg := CreateUserParams(createDtoToModel(dto))
	user, err := s.db.CreateUser(context.Background(), arg)
	dtoRes := modelToDto(user)
	if err != nil {
		return dtoRes, checkError(err)
	}
	return dtoRes, nil
}

func (s *service) updateUser(dto userDto) (userDto, error) {
	m := createDtoToModel(dto)
	arg := UpdateUserParams{m.Name, m.Email, m.Phone, m.Age, m.Address, m.ID}
	user, err := s.db.UpdateUser(context.Background(), arg)
	dtoRes := modelToDto(user)
	if err != nil {
		return dtoRes, checkError(err)
	}
	return dtoRes, nil
}

func (s *service) deleteUser(id int32) error {
	err := s.db.DeleteUser(context.Background(), id)
	if err != nil {
		return checkError(err)
	}
	return nil
}

func (s *service) UserExists(n, p string) error {
	arg := UserExistsParams{n, p}
	_, err := s.db.UserExists(context.Background(), arg)
	if err != nil {
		return checkError(err)
	}
	return nil
}
