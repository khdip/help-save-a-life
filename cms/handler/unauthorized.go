package handler

import "net/http"

func (h *Handler) unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "unauthorized.html", HomeTemplateData{
		URLs: listOfURLs(),
		Sett: h.getSettingsHome(w, r),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
