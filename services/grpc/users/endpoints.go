package users

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/go-crud-apis/services/grpc/gen/proto"
	"github.com/go-crud-apis/services/grpc/users/sql/dbgen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	queries *dbgen.Queries
	pb.UnimplementedUserServiceServer
}

func NewUserService(q *dbgen.Queries) (*UserService, error) {
	return &UserService{queries: q}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if err := validate(req); err != nil {
		log.Printf("Error validating user %v", err)
		return nil, err
	}
	param := dbgen.CreateUserParams{ID: req.GetId(),
		Name:     req.GetName(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Phone:    sql.NullString{String: req.GetPhone(), Valid: true},
		Age:      sql.NullInt32{Int32: req.GetAge(), Valid: true},
		Address:  sql.NullString{String: req.GetAddress(), Valid: true},
	}
	err := s.queries.CreateUser(context.Background(), param)
	if err != nil {
		log.Printf("Error creating user %v", err)
		return nil, checkError(err)
	}
	return req, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.queries.GetUser(context.Background(), req.GetId())
	if err != nil {
		log.Printf("Error getting user %v", err)
		return nil, checkError(err)
	}
	return &pb.User{Id: user.ID,
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
		Phone:    user.Phone.String,
		Age:      user.Age.Int32,
		Address:  user.Address.String,
	}, nil
}

func (s *UserService) ListUser(ctx context.Context, empty *emptypb.Empty) (*pb.ListUserResponse, error) {
	users, err := s.queries.ListUsers(context.Background())
	if err != nil {
		log.Printf("Error listing users %v", err)
		return nil, checkError(err)
	}
	res := &pb.ListUserResponse{Users: make([]*pb.User, 0)}
	for _, user := range users {
		res.Users = append(res.Users, &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Email:    user.Email,
			Phone:    user.Phone.String,
			Age:      user.Age.Int32,
			Address:  user.Address.String,
		})
	}
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	param := dbgen.UpdateUserParams{
		ID:      req.GetId(),
		Name:    req.GetUser().GetName(),
		Email:   req.GetUser().GetEmail(),
		Phone:   sql.NullString{String: req.GetUser().GetPhone(), Valid: true},
		Age:     sql.NullInt32{Int32: req.GetUser().GetAge(), Valid: true},
		Address: sql.NullString{String: req.GetUser().GetAddress(), Valid: true},
	}
	err := s.queries.UpdateUser(context.Background(), param)
	if err != nil {
		log.Printf("Error creating user %v", err)
		return nil, checkError(err)
	}
	return req.GetUser(), nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.queries.DeleteUser(context.Background(), req.GetId())
	if err != nil {
		return nil, checkError(err)
	}
	return &emptypb.Empty{}, err
}
