package accounts

import (
	"context"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) AccountsStats(ctx context.Context, req *acntgrpc.AccountsStatsRequest) (*acntgrpc.AccountsStatsResponse, error) {
	r, err := s.accst.AccountsStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "account doesn't exist")
	}
	return &acntgrpc.AccountsStatsResponse{
		Stats: &acntgrpc.Stats{
			Count: r.Count,
		},
	}, nil
}
