package users

import (
	"context"

	usergrpc "help-save-a-life/proto/users"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListUser(ctx context.Context, req *usergrpc.ListUserRequest) (*usergrpc.ListUserResponse, error) {
	users, err := s.ust.ListUser(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user found")
	}

	list := make([]*usergrpc.User, len(users))
	for i, u := range users {
		list[i] = &usergrpc.User{
			UserID:       u.UserID,
			SerialNumber: u.SerialNumber,
			Name:         u.Name,
			Batch:        u.Batch,
			Email:        u.Email,
			CreatedAt:    tspb.New(u.CreatedAt),
			CreatedBy:    u.CreatedBy,
			UpdatedAt:    tspb.New(u.UpdatedAt),
			UpdatedBy:    u.UpdatedBy,
		}
	}

	return &usergrpc.ListUserResponse{
		User: list,
	}, nil
}
