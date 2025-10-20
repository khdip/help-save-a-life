package settings

import (
	"context"

	settgrpc "help-save-a-life/proto/settings"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateSettings(ctx context.Context, req *settgrpc.UpdateSettingsRequest) (*settgrpc.UpdateSettingsResponse, error) {
	res, err := s.sst.UpdateSettings(ctx, storage.Settings{
		PatientName:                  req.Sett.PatientName,
		Title:                        req.Sett.Title,
		BannerTitle:                  req.Sett.BannerTitle,
		HighlightedBannerTitle:       req.Sett.HighlightedBannerTitle,
		BannerDescription:            req.Sett.BannerDescription,
		HighlightedBannerDescription: req.Sett.HighlightedBannerDescription,
		BannerImage:                  req.Sett.BannerImage,
		AboutPatient:                 req.Sett.AboutPatient,
		TargetAmount:                 req.Sett.TargetAmount,
		ShowMedicalDocuments:         req.Sett.ShowMedicalDocuments,
		ShowCollection:               req.Sett.ShowCollection,
		ShowDailyReport:              req.Sett.ShowDailyReport,
		ShowFundUpdates:              req.Sett.ShowFundUpdates,
		CalculateCollection:          req.Sett.CalculateCollection,
		UpdatedBy:                    req.Sett.UpdatedBy,
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &settgrpc.UpdateSettingsResponse{
		Sett: &settgrpc.Settings{
			PatientName:                  res.PatientName,
			Title:                        res.Title,
			BannerTitle:                  res.BannerTitle,
			HighlightedBannerTitle:       res.HighlightedBannerTitle,
			BannerDescription:            res.BannerDescription,
			HighlightedBannerDescription: res.HighlightedBannerDescription,
			BannerImage:                  res.BannerImage,
			AboutPatient:                 res.AboutPatient,
			TargetAmount:                 res.TargetAmount,
			ShowMedicalDocuments:         res.ShowMedicalDocuments,
			ShowCollection:               res.ShowCollection,
			ShowDailyReport:              res.ShowDailyReport,
			ShowFundUpdates:              res.ShowFundUpdates,
		},
	}, nil
}
