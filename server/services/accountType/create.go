package accountType

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acctgrpc "help-save-a-life/proto/accountType"
	"help-save-a-life/server/storage"
)

func (s *Svc) CreateAccountType(ctx context.Context, req *acctgrpc.CreateAccountTypeRequest) (*acctgrpc.CreateAccountTypeResponse, error) {
	res, err := s.atst.CreateAccountType(ctx, storage.AccountType{
		Title: req.Acct.Title,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.Acct.CreatedBy,
			UpdatedBy: req.Acct.UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create account type")
	}

	return &acctgrpc.CreateAccountTypeResponse{
		ID: res,
	}, nil
}
