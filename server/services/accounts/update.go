package accounts

import (
	"context"

	acntgrpc "help-save-a-life/proto/accounts"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateAccounts(ctx context.Context, req *acntgrpc.UpdateAccountsRequest) (*acntgrpc.UpdateAccountsResponse, error) {
	res, err := s.accst.UpdateAccounts(ctx, storage.Accounts{
		ID:          req.Acnt.ID,
		AccountType: req.Acnt.AccountType,
		Details:     req.Acnt.Details,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Acnt.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &acntgrpc.UpdateAccountsResponse{
		Acnt: &acntgrpc.Accounts{
			ID:          res.ID,
			AccountType: res.AccountType,
			Details:     res.Details,
		},
	}, nil
}
