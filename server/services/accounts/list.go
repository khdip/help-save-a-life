package accounts

import (
	"context"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListAccounts(ctx context.Context, req *acntgrpc.ListAccountsRequest) (*acntgrpc.ListAccountsResponse, error) {
	acnt, err := s.accst.ListAccounts(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no accounts found")
	}

	list := make([]*acntgrpc.Accounts, len(acnt))
	for i, a := range acnt {
		list[i] = &acntgrpc.Accounts{
			ID:           a.ID,
			SerialNumber: a.SerialNumber,
			AccountType:  a.AccountType,
			Details:      a.Details,
			CreatedAt:    tspb.New(a.CreatedAt),
			CreatedBy:    a.CreatedBy,
			UpdatedAt:    tspb.New(a.UpdatedAt),
			UpdatedBy:    a.UpdatedBy,
		}
	}

	return &acntgrpc.ListAccountsResponse{
		Acnt: list,
	}, nil
}
