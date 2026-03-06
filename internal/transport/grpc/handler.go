package grpc

import (
	"context"

	userpb "github.com/AntonRadchenko/project-protos/proto/user"
	"github.com/AntonRadchenko/users-service/internal/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc *user.UserService
	userpb.UnimplementedUserServiceServer
}

func NewHandler(svc *user.UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	// Преобразовать req → user.User
	params := user.CreateUserParams{
		Email:    req.Email,
		Password: req.Password,
	}

	// Вызвать svc.CreateUser
	createdUser, err := h.svc.CreateUser(params)
	if err != nil {
		return nil, err
	}

	// Вернуть &userpb.CreateUserResponse{User: …}
	return &userpb.User{
		Id:    uint32(createdUser.ID),
		Email: createdUser.Email,
	}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	// Вызвать svc.GetUserByID(req.Id)
	user, err := h.svc.GetUser(uint(req.Id))
	if err != nil {
		return nil, err
	}

	// Вернуть &userpb.GetUserResponse{User: …}
	return &userpb.User{
		Id:    req.Id,
		Email: user.Email,
	}, nil
}

func (h *Handler) ListUsers(ctx context.Context, _ *emptypb.Empty) (*userpb.UserList, error) {
	// - Вызвать `svc.GetAllUsers()`
	users, err := h.svc.GetUsers()
	if err != nil {
		return nil, err
	}

	// - Преобразовать срез `user.User` → `[]*userpb.User`
	pbUsers := make([]*userpb.User, len(users))
	for i, u := range users {
		pbUsers[i] = &userpb.User{
			Id:    uint32(u.ID),
            Email: u.Email,
		}
	}

	// - Вернуть `&userpb.ListUsersResponse{Users: …}` (возвращаем nil, когда все отправили)
	return &userpb.UserList{Users: pbUsers}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	// Создаём params из req
	params := user.UpdateUserParams{
		Email:    req.Email, // optional поля приходят как указатели
		Password: req.Password,
	}

	// Вызвать svc.UpdateUserByID(req.Id, user.User{…})
	updatedUser, err := h.svc.UpdateUser(uint(req.Id), params)
	if err != nil {
		return nil, err
	}

	// Вернуть &userpb.UpdateUserResponse{User: …}
	return &userpb.User{
		Id:    uint32(updatedUser.ID),
		Email: updatedUser.Email,
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	// Вызвать svc.DeleteUserByID(req.Id)
	err := h.svc.DeleteUser(uint(req.Id))
	if err != nil {
		return nil, err
	}

	// Вернуть пустой &userpb.DeleteUserResponse{}
	return &emptypb.Empty{}, nil
}

