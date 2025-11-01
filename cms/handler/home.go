package handler

import (
	"fmt"
	"help-save-a-life/cms/paginator"
	"log"
	"net/http"
)

type HomeTemplateData struct {
	CollList        []Collection
	DreList         []DailyReport
	Sett            SettingsHome
	Paginator       paginator.Paginator
	TargetAmount    string
	CollectedAmount string
	RemainingAmount string
	Percentage      string
	URLs            map[string]string
}

func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("index.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	sett := h.getSettingsHome(w, r)
	totalCollection, targetAmount := h.getTotalAndTargetAmount(w, r, sett)

	data := HomeTemplateData{
		Sett:            sett,
		URLs:            listOfURLs(),
		TargetAmount:    formatWithCommas(targetAmount),
		CollectedAmount: formatWithCommas(totalCollection),
		RemainingAmount: formatWithCommas(targetAmount - totalCollection),
		Percentage:      fmt.Sprintf("%.2f", ((float64(totalCollection) / float64(targetAmount)) * 100)),
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
		return
	}
}
