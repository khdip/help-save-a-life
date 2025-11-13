package accounts

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) CreateAccounts(ctx context.Context, req *acntgrpc.CreateAccountsRequest) (*acntgrpc.CreateAccountsResponse, error) {
	res, err := s.accst.CreateAccounts(ctx, storage.Accounts{
		AccountType: req.Acnt.AccountType,
		Details:     req.Acnt.Details,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.Acnt.CreatedBy,
			UpdatedBy: req.Acnt.UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create account")
	}

	return &acntgrpc.CreateAccountsResponse{
		ID: res,
	}, nil
}
