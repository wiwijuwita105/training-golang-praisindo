package grpc

import (
	"context"
	"fmt"
	"log"
	"user_svc/entity"
	pb "user_svc/proto/user_service/v1"
	"user_svc/service"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandler is used to implement UnimplementedUserServiceServer
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var usersProto []*pb.User
	for _, user := range users {
		usersProto = append(usersProto, &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.GetUsersResponse{
		Users: usersProto,
	}, nil
}

func (u *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := u.userService.GetUserByID(ctx, int(req.GetId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return res, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.MutationResponse, error) {
	createdUser, err := u.userService.CreateUser(ctx, &entity.User{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created user with ID %d", createdUser.ID),
		Id:      int32(createdUser.ID),
	}, nil
}
