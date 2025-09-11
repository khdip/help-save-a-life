package handler

import (
	dregrpc "help-save-a-life/proto/dailyReport"
	"log"
	"net/http"
)

type DashBoardData struct {
	TargetAmount    int32
	CollectedAmount int32
	RemainingAmount int32
	URLs            map[string]string
}

func (h *Handler) viewDashboard(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("dashboard.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	drstat, err := h.drc.DailyReportStats(r.Context(), &dregrpc.DailyReportStatsRequest{
		Filter: &dregrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	data := DashBoardData{
		TargetAmount:    targetamount,
		CollectedAmount: drstat.Stats.TotalAmount,
		RemainingAmount: targetamount - drstat.Stats.TotalAmount,
		URLs:            listOfURLs(),
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}
