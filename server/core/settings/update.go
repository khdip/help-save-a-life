package settings

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateSettings(ctx context.Context, sett storage.Settings) (*storage.Settings, error) {
	st, err := s.st.UpdateSettings(ctx, sett)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return st, nil
}
