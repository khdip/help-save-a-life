package users

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usergrpc "help-save-a-life/proto/users"
	"help-save-a-life/server/storage"
)

func (s *Svc) CreateUser(ctx context.Context, req *usergrpc.CreateUserRequest) (*usergrpc.CreateUserResponse, error) {
	res, err := s.ust.CreateUser(ctx, storage.User{
		UserID:   req.GetUser().UserID,
		Name:     req.GetUser().Name,
		Batch:    req.GetUser().Batch,
		Email:    req.GetUser().Email,
		Password: req.GetUser().Password,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.GetUser().CreatedBy,
			UpdatedBy: req.GetUser().UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &usergrpc.CreateUserResponse{
		UserID: res,
	}, nil
}
