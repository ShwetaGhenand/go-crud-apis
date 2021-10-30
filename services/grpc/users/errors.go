package users

import (
	"fmt"
	"strings"

	pb "github.com/go-crud-apis/services/grpc/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkError(err error) error {
	switch e := err.Error(); {
	case strings.Contains(e, "no rows"):
		return status.Error(codes.NotFound, "record not found")
	case strings.Contains(e, "users_pkey"):
		return status.Error(codes.AlreadyExists, "user already exists")
	default:
		return status.Error(codes.Internal, fmt.Errorf("database error %w", err).Error())
	}
}

func validate(user *pb.User) error {
	if user.GetId() <= 0 {
		return status.Error(codes.InvalidArgument, "id required")
	} else if user.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name required")
	} else if user.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	} else if user.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	return nil
}
