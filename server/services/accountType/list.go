package accountType

import (
	"context"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListAccountType(ctx context.Context, req *acctgrpc.ListAccountTypeRequest) (*acctgrpc.ListAccountTypeResponse, error) {
	acct, err := s.atst.ListAccountType(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no account type found")
	}

	list := make([]*acctgrpc.AccountType, len(acct))
	for i, a := range acct {
		list[i] = &acctgrpc.AccountType{
			ID:           a.ID,
			SerialNumber: a.SerialNumber,
			Title:        a.Title,
			CreatedAt:    tspb.New(a.CreatedAt),
			CreatedBy:    a.CreatedBy,
			UpdatedAt:    tspb.New(a.UpdatedAt),
			UpdatedBy:    a.UpdatedBy,
		}
	}

	return &acctgrpc.ListAccountTypeResponse{
		Acct: list,
	}, nil
}
