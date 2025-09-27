package users

import (
	"context"

	usergrpc "help-save-a-life/proto/users"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetUser(ctx context.Context, req *usergrpc.GetUserRequest) (*usergrpc.GetUserResponse, error) {
	r, err := s.ust.GetUser(ctx, storage.User{
		UserID: req.GetUser().UserID,
		Email:  req.GetUser().Email,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	}
	return &usergrpc.GetUserResponse{
		User: &usergrpc.User{
			UserID:       r.UserID,
			SerialNumber: r.SerialNumber,
			Name:         r.Name,
			Batch:        r.Batch,
			Email:        r.Email,
			Password:     r.Password,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
