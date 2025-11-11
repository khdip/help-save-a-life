package handler

import (
	"fmt"
	"log"
	"net/http"
)

type HomeTemplateData struct {
	Sett            SettingsHome
	Link            []Link
	MedDocs         []MedDocs
	Accounts        map[string][]Accounts
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
		Link:            h.getLinkList(w, r),
		MedDocs:         h.getMedDocsList(w, r),
		Accounts:        h.getAccounts(w, r),
		URLs:            listOfURLs(),
		TargetAmount:    formatWithCommas(float32(targetAmount)),
		CollectedAmount: formatWithCommas(totalCollection),
		RemainingAmount: formatWithCommas(float32(targetAmount) - totalCollection),
		Percentage:      fmt.Sprintf("%.2f", ((float64(totalCollection) / float64(targetAmount)) * 100)),
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
		return
	}
}
