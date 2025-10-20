package settings

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetSettings(ctx context.Context, sett storage.Settings) (*storage.Settings, error) {
	st, err := s.st.GetSettings(ctx, sett)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return st, nil
}
