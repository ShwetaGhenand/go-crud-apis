package user

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/go-crud-apis/users/sql/dbgen"
)

type Service struct {
	db *db.Queries
}

func NewService(db *db.Queries) *Service {
	return &Service{db: db}
}

func reqToModel(req JSONUser) db.User {
	m := db.User{}
	m.ID = int32(req.ID)
	m.Name = req.Name
	m.Password = req.Password
	m.Email = req.Email
	if req.Phone != "" {
		m.Phone = sql.NullString{String: req.Phone, Valid: true}
	} else {
		m.Phone = sql.NullString{String: req.Phone, Valid: false}
	}
	if req.Address != "" {
		m.Address = sql.NullString{String: req.Address, Valid: true}
	} else {
		m.Address = sql.NullString{String: req.Address, Valid: false}
	}
	if req.Age != 0 {
		m.Age = sql.NullInt32{Int32: int32(req.Age), Valid: true}
	} else {
		m.Age = sql.NullInt32{Int32: int32(req.Age), Valid: false}
	}
	return m
}

func modelToRes(user db.User) JSONUser {
	res := JSONUser{}
	if user.ID != 0 {
		res.ID = int(user.ID)
	}
	if user.Name != "" {
		res.Name = user.Name
	}
	if user.Password != "" {
		res.Password = user.Password
	}
	if user.Email != "" {
		res.Email = user.Email
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

func (s *Service) ListUsers() ([]JSONUser, error) {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return nil, CheckError(err)
	}
	if len(users) == 0 {
		return nil, &CustomErr{errors.New("users not found"), 404}
	}
	var res []JSONUser
	for _, user := range users {
		res = append(res, modelToRes(user))
	}
	return res, nil
}

func (s *Service) GetUser(id int32) (JSONUser, error) {
	user, err := s.db.GetUser(context.Background(), id)
	dto := modelToRes(user)
	if err != nil {
		return dto, CheckError(err)
	}
	return dto, nil
}

func (s *Service) CreateUser(req JSONUser) error {
	arg := db.CreateUserParams(reqToModel(req))
	err := s.db.CreateUser(context.Background(), arg)
	if err != nil {
		return CheckError(err)
	}
	return nil
}

func (s *Service) UpdateUser(req JSONUser) error {
	m := reqToModel(req)
	arg := db.UpdateUserParams{Name: m.Name, Email: m.Email, Phone: m.Phone,
		Age: m.Age, Address: m.Address, ID: m.ID}
	err := s.db.UpdateUser(context.Background(), arg)
	if err != nil {
		return CheckError(err)
	}
	return nil
}

func (s *Service) DeleteUser(id int32) error {
	err := s.db.DeleteUser(context.Background(), id)
	if err != nil {
		return CheckError(err)
	}
	return nil
}

func (s *Service) UserExists(n, p string) error {
	arg := db.UserExistsParams{Name: n, Password: p}
	_, err := s.db.UserExists(context.Background(), arg)
	if err != nil {
		return CheckError(err)
	}
	return nil
}
