package users_test

import (
	"context"
	"testing"

	"github.com/go-crud-apis/services/grpc/cmd/server"
	userspb "github.com/go-crud-apis/services/grpc/proto/gen"
	users "github.com/go-crud-apis/services/grpc/users"
	"github.com/go-crud-apis/services/grpc/users/sql/dbgen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func startDB(t *testing.T) *dbgen.Queries {
	t.Helper()
	conf := server.DBConfig{PGUser: "postgres",
		PGPassword: "zyxwv",
		PGDatabase: "userdb",
		PGHost:     "localhost",
		PGPort:     5432,
		PGSSLMode:  "disable",
	}
	db, err := server.InitDB(conf)
	if err != nil {
		t.Fatalf("Failed to initialize database connection: %s", err)
	}
	return db
}

func TestCreateUserEndpoints(t *testing.T) {
	var service *users.UserService
	service, err := users.NewUserService(startDB(t))
	if err != nil {
		t.Fatalf("Failed to create a new user service: %s", err)
	}
	// Arrange
	var id, age int32 = 10101, 32
	name := "Foo"
	password := "xxx"
	email := "foo@abc.com"
	phone := "3432423122"
	address := "abc"

	t.Run("Succeed if valid user passed", func(t *testing.T) {
		// Act
		user, err := service.CreateUser(context.Background(), &userspb.User{
			Id:       id,
			Name:     name,
			Password: password,
			Email:    email,
			Phone:    phone,
			Age:      age,
			Address:  address,
		})
		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, id, user.GetId())
		require.Equal(t, name, user.GetName())
		require.Equal(t, password, user.GetPassword())
		require.Equal(t, phone, user.GetPhone())
		require.Equal(t, email, user.GetEmail())
		require.Equal(t, age, user.GetAge())
		require.Equal(t, address, user.GetAddress())
	})

	t.Run("Failed if duplicate user passed", func(t *testing.T) {
		// Act
		user, err := service.CreateUser(context.Background(), &userspb.User{
			Id:       id,
			Name:     name,
			Password: password,
			Email:    email,
			Phone:    phone,
			Age:      age,
			Address:  address,
		})
		// Assert
		require.Error(t, err)
		require.Nil(t, user)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.AlreadyExists, st.Code())
	})
	t.Run("Failed if invalid user passed", func(t *testing.T) {
		// Arrange
		fixtures := []struct {
			name string
			user *userspb.User
			code codes.Code
			err  string
		}{
			{
				name: "failed_without_id",
				user: &userspb.User{
					Name:     name,
					Password: password,
					Email:    email,
					Phone:    phone,
					Age:      age,
					Address:  address,
				},
				code: codes.InvalidArgument,
				err:  "id required",
			},
			{
				name: "failed_without_name",
				user: &userspb.User{
					Id:       101012,
					Password: password,
					Email:    email,
					Phone:    phone,
					Age:      age,
					Address:  address,
				},
				code: codes.InvalidArgument,
				err:  "name required",
			},
			{
				name: "failed_without_password",
				user: &userspb.User{
					Id:      101013,
					Name:    name,
					Email:   email,
					Phone:   phone,
					Age:     age,
					Address: address,
				},
				code: codes.InvalidArgument,
				err:  "password required",
			},
			{
				name: "failed_without_email",
				user: &userspb.User{
					Id:       101014,
					Name:     name,
					Password: password,
					Phone:    phone,
					Age:      age,
					Address:  address,
				},
				code: codes.InvalidArgument,
				err:  "email required",
			},
		}

		// Act
		for i := range fixtures {
			f := fixtures[i]
			t.Run(f.name, func(t *testing.T) {
				t.Parallel()
				res, err := service.CreateUser(context.Background(), f.user)
				// Assert
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, f.code, st.Code())
				require.Equal(t, f.err, st.Message())
			})
		}
	})
}

func TestGetUsersEndpoint(t *testing.T) {
	var service *users.UserService
	service, err := users.NewUserService(startDB(t))
	if err != nil {
		t.Fatalf("Failed to create a new user service: %s", err)
	}

	t.Run("Succeed if user exists", func(t *testing.T) {
		// Act
		user, err := service.GetUser(context.Background(), &userspb.GetUserRequest{Id: 10101})
		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user.GetId())
		require.NotEmpty(t, user.GetName())
		require.NotEmpty(t, user.GetPassword())
		require.NotEmpty(t, user.GetPhone())
		require.NotEmpty(t, user.GetEmail())
		require.NotEmpty(t, user.GetAge())
		require.NotEmpty(t, user.GetAddress())
	})

	t.Run("Failed if user does not exits", func(t *testing.T) {
		// Act
		user, err := service.GetUser(context.Background(), &userspb.GetUserRequest{Id: 2})
		// Assert
		require.Error(t, err)
		require.Nil(t, user)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
		require.Equal(t, "record not found", st.Message())
	})
}

func TestListUserEndpoint(t *testing.T) {
	var service *users.UserService
	service, err := users.NewUserService(startDB(t))
	if err != nil {
		t.Fatalf("Failed to create a new user service: %s", err)
	}
	t.Run("Success", func(t *testing.T) {
		// Act
		res, err := service.ListUser(context.Background(), &emptypb.Empty{})
		// Assert
		require.NoError(t, err)
		require.NotNil(t, res)
		for _, user := range res.GetUsers() {
			require.NotEmpty(t, user.GetId())
			require.NotEmpty(t, user.GetName())
			require.NotEmpty(t, user.GetPassword())
			require.NotEmpty(t, user.GetPhone())
			require.NotEmpty(t, user.GetEmail())
			require.NotEmpty(t, user.GetAge())
			require.NotEmpty(t, user.GetAddress())
		}
	})
}

func TestUpdateUserEndpoints(t *testing.T) {
	var service *users.UserService
	service, err := users.NewUserService(startDB(t))
	if err != nil {
		t.Fatalf("Failed to create a new user service: %s", err)
	}

	t.Run("Succeed if valid user passed", func(t *testing.T) {
		// Act
		user, err := service.UpdateUser(context.Background(), &userspb.UpdateUserRequest{
			Id:   10101,
			User: &userspb.User{Phone: "890989796"},
		})
		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "890989796", user.GetPhone())
	})
}

func TestDeleteUserEndpoints(t *testing.T) {
	var service *users.UserService
	service, err := users.NewUserService(startDB(t))
	if err != nil {
		t.Fatalf("Failed to create a new user service: %s", err)
	}

	t.Run("Succeed if valid user passed", func(t *testing.T) {
		// Act
		_, err := service.DeleteUser(context.Background(), &userspb.DeleteUserRequest{Id: 10101})
		// Assert
		require.NoError(t, err)
	})
}
