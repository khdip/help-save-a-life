package accounts

import (
	"context"

	acntgrpc "help-save-a-life/proto/accounts"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetAccounts(ctx context.Context, req *acntgrpc.GetAccountsRequest) (*acntgrpc.GetAccountsResponse, error) {
	r, err := s.accst.GetAccounts(ctx, storage.Accounts{
		ID: req.Acnt.ID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "account doesn't exist")
	}
	return &acntgrpc.GetAccountsResponse{
		Acnt: &acntgrpc.Accounts{
			ID:           r.ID,
			SerialNumber: r.SerialNumber,
			AccountType:  r.AccountType,
			Details:      r.Details,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
