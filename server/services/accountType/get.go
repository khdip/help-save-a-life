package accountType

import (
	"context"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetAccountType(ctx context.Context, req *acctgrpc.GetAccountTypeRequest) (*acctgrpc.GetAccountTypeResponse, error) {
	r, err := s.atst.GetAccountType(ctx, storage.AccountType{
		ID: req.Acct.ID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "account type doesn't exist")
	}
	return &acctgrpc.GetAccountTypeResponse{
		Acct: &acctgrpc.AccountType{
			ID:           r.ID,
			SerialNumber: r.SerialNumber,
			Title:        r.Title,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
