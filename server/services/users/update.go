package users

import (
	"context"

	usergrpc "help-save-a-life/proto/users"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateUser(ctx context.Context, req *usergrpc.UpdateUserRequest) (*usergrpc.UpdateUserResponse, error) {
	res, err := s.ust.UpdateUser(ctx, storage.User{
		UserID:   req.GetUser().UserID,
		Name:     req.GetUser().Name,
		Batch:    req.GetUser().Batch,
		Email:    req.GetUser().Email,
		Password: req.GetUser().Password,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.GetUser().UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &usergrpc.UpdateUserResponse{
		User: &usergrpc.User{
			UserID:   res.UserID,
			Name:     res.Name,
			Batch:    res.Batch,
			Email:    res.Email,
			Password: res.Password,
		},
	}, nil
}
