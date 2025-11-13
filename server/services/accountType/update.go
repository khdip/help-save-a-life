package accountType

import (
	"context"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateAccountType(ctx context.Context, req *acctgrpc.UpdateAccountTypeRequest) (*acctgrpc.UpdateAccountTypeResponse, error) {
	res, err := s.atst.UpdateAccountType(ctx, storage.AccountType{
		ID:    req.Acct.ID,
		Title: req.Acct.Title,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Acct.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &acctgrpc.UpdateAccountTypeResponse{
		Acct: &acctgrpc.AccountType{
			ID:    res.ID,
			Title: res.Title,
		},
	}, nil
}
