package handler

import (
	collgrpc "help-save-a-life/proto/collection"
	dregrpc "help-save-a-life/proto/dailyReport"
	settgrpc "help-save-a-life/proto/settings"
	"log"
	"net/http"
)

type DashBoardData struct {
	TargetAmount    int32
	CollectedAmount int32
	RemainingAmount int32
	URLs            map[string]string
	CurrentPageURL  string
	Title           string
}

func (h *Handler) viewDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	template := h.templates.Lookup("dashboard.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	sett, err := h.sc.GetSettings(ctx, &settgrpc.GetSettingsRequest{})
	if err != nil {
		log.Println("unable to get settings: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	collstat, err := h.cc.CollectionStats(ctx, &collgrpc.CollectionStatsRequest{
		Filter: &collgrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get collection stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	drstat, err := h.drc.DailyReportStats(ctx, &dregrpc.DailyReportStatsRequest{
		Filter: &dregrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get daily report stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	var totalCollection int32
	switch sett.Sett.CalculateCollection {
	case 0:
		totalCollection = collstat.Stats.TotalAmount
	case 1:
		totalCollection = drstat.Stats.TotalAmount
	case 2:
		totalCollection = sett.Sett.TotalAmount
	}

	targetAmount := sett.Sett.TargetAmount

	data := DashBoardData{
		TargetAmount:    targetAmount,
		CollectedAmount: totalCollection,
		RemainingAmount: targetAmount - totalCollection,
		URLs:            listOfURLs(),
		CurrentPageURL:  dashboardPath,
		Title:           h.getSettingsTitle(w, r),
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}
