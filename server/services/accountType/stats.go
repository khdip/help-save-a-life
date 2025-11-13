package accountType

import (
	"context"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) AccountTypeStats(ctx context.Context, req *acctgrpc.AccountTypeStatsRequest) (*acctgrpc.AccountTypeStatsResponse, error) {
	r, err := s.atst.AccountTypeStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "account type doesn't exist")
	}
	return &acctgrpc.AccountTypeStatsResponse{
		Stats: &acctgrpc.Stats{
			Count: r.Count,
		},
	}, nil
}
