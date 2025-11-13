package settings

import (
	"context"

	settgrpc "github.com/khdip/help-save-a-life/proto/settings"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetSettings(ctx context.Context, req *settgrpc.GetSettingsRequest) (*settgrpc.GetSettingsResponse, error) {
	r, err := s.sst.GetSettings(ctx, storage.Settings{})
	if err != nil {
		return nil, status.Error(codes.NotFound, "settings doesn't exist")
	}
	return &settgrpc.GetSettingsResponse{
		Sett: &settgrpc.Settings{
			PatientName:                  r.PatientName,
			Title:                        r.Title,
			BannerTitle:                  r.BannerTitle,
			HighlightedBannerTitle:       r.HighlightedBannerTitle,
			BannerDescription:            r.BannerDescription,
			HighlightedBannerDescription: r.HighlightedBannerDescription,
			BannerImage:                  r.BannerImage,
			AboutPatient:                 r.AboutPatient,
			TargetAmount:                 r.TargetAmount,
			ShowMedicalDocuments:         r.ShowMedicalDocuments,
			ShowCollection:               r.ShowCollection,
			ShowDailyReport:              r.ShowDailyReport,
			ShowFundUpdates:              r.ShowFundUpdates,
			CalculateCollection:          r.CalculateCollection,
			TotalAmount:                  r.TotalAmount,
			UpdatedBy:                    r.UpdatedBy,
		},
	}, nil
}
