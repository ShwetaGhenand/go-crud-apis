package users

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "github.com/go-crud-apis/services/grpc/gen/go/proto"
	"github.com/go-crud-apis/services/grpc/users/sql/dbgen"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	queries *dbgen.Queries
	pb.UnimplementedUserServiceServer
}

func NewUserService(q *dbgen.Queries) (*UserService, error) {
	return &UserService{queries: q}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if err := emptyCheck(req); err != nil {
		log.Printf("Empty field %v", err)
		return nil, err
	}
	if err := validate(req.GetEmail()); err != nil {
		log.Printf("Error validating email %v", err)
		return nil, err
	}
	param := dbgen.CreateUserParams{ID: req.GetId(),
		Name:       req.GetName(),
		Password:   req.GetPassword(),
		Email:      req.GetEmail(),
		Phone:      sql.NullString{String: req.GetPhone(), Valid: true},
		Age:        sql.NullInt32{Int32: req.GetAge(), Valid: true},
		Address:    sql.NullString{String: req.GetAddress(), Valid: true},
		Createtime: time.Now(),
		Updatetime: time.Now(),
	}
	err := s.queries.CreateUser(context.Background(), param)
	if err != nil {
		log.Printf("Error creating user %v", err)
		return nil, checkError(err)
	}
	res := req
	res.CreateTime = timestamppb.New(param.Createtime)
	res.UpdateTime = timestamppb.New(param.Updatetime)
	return res, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.queries.GetUser(context.Background(), req.GetId())
	if err != nil {
		log.Printf("Error getting user %v", err)
		return nil, checkError(err)
	}
	return &pb.User{Id: user.ID,
		Name:       user.Name,
		Password:   user.Password,
		Email:      user.Email,
		Phone:      user.Phone.String,
		Age:        user.Age.Int32,
		Address:    user.Address.String,
		CreateTime: timestamppb.New(user.Createtime),
		UpdateTime: timestamppb.New(user.Updatetime),
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
			Id:         user.ID,
			Name:       user.Name,
			Password:   user.Password,
			Email:      user.Email,
			Phone:      user.Phone.String,
			Age:        user.Age.Int32,
			Address:    user.Address.String,
			CreateTime: timestamppb.New(user.Createtime),
			UpdateTime: timestamppb.New(user.Updatetime),
		})
	}
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req.GetUser().GetEmail() != "" {
		if err := validate(req.GetUser().GetEmail()); err != nil {
			log.Printf("Error validating email %v", err)
			return nil, err
		}
	}
	user, err := s.queries.GetUser(context.Background(), req.GetId())
	if err != nil {
		log.Printf("Error getting user %v", err)
		return nil, checkError(err)
	}
	param := dbgen.UpdateUserParams{
		Name:       req.GetUser().GetName(),
		Email:      req.GetUser().GetEmail(),
		Phone:      sql.NullString{String: req.GetUser().GetPhone(), Valid: true},
		Age:        sql.NullInt32{Int32: req.GetUser().GetAge(), Valid: true},
		Address:    sql.NullString{String: req.GetUser().GetAddress(), Valid: true},
		Updatetime: time.Now(),
	}
	err = s.queries.UpdateUser(context.Background(), param)
	if err != nil {
		log.Printf("Error creating user %v", err)
		return nil, checkError(err)
	}
	res := req.GetUser()
	res.Id = user.ID
	res.CreateTime = timestamppb.New(user.Createtime)
	res.UpdateTime = timestamppb.New(param.Updatetime)
	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := s.queries.GetUser(context.Background(), req.GetId())
	if err != nil {
		log.Printf("Error getting user %v", err)
		return nil, checkError(err)
	}
	err = s.queries.DeleteUser(context.Background(), req.GetId())
	if err != nil {
		return nil, checkError(err)
	}
	return &emptypb.Empty{}, err
}
