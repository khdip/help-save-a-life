package handler

import (
	"help-save-a-life/cms/paginator"
	dregrpc "help-save-a-life/proto/dailyReport"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DailyReport struct {
	ReportID     string
	SerialNumber int32
	Date         string
	Amount       int32
	Currency     string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    time.Time
	DeletedBy    string
}

type DreTemplateData struct {
	Dre            DailyReport
	List           []DailyReport
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	CurrentPageURL string
}

func (h *Handler) createDailyReport(w http.ResponseWriter, r *http.Request) {
	h.loadDailyReportCreateForm(w, DailyReport{})
}

func (h *Handler) storeDailyReport(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dre DailyReport
	err = h.decoder.Decode(&dre, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.drc.CreateDailyReport(r.Context(), &dregrpc.CreateDailyReportRequest{
		Dre: &dregrpc.DailyReport{
			Date:      dre.Date,
			Amount:    dre.Amount,
			Currency:  dre.Currency,
			CreatedBy: h.getLoggedUser(r),
			UpdatedBy: h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, dailyReportListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editDailyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["report_id"]
	res, err := h.drc.GetDailyReport(r.Context(), &dregrpc.GetDailyReportRequest{
		Dre: &dregrpc.DailyReport{
			ReportID: id,
		},
	})
	if err != nil {
		log.Println("unable to get daily report info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadDailyReportEditForm(w, DailyReport{
		ReportID: res.Dre.ReportID,
		Date:     res.Dre.Date,
		Amount:   res.Dre.Amount,
		Currency: res.Dre.Currency,
	})
}

func (h *Handler) updateDailyReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["report_id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var dre DailyReport
	if err := h.decoder.Decode(&dre, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.drc.UpdateDailyReport(ctx, &dregrpc.UpdateDailyReportRequest{
		Dre: &dregrpc.DailyReport{
			ReportID:  id,
			Date:      dre.Date,
			Amount:    dre.Amount,
			Currency:  dre.Currency,
			UpdatedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, dailyReportListPath, http.StatusSeeOther)
}

func (h *Handler) listDailyReport(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("dre-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filterData := GetFilterData(r)
	drlst, err := h.drc.ListDailyReport(r.Context(), &dregrpc.ListDailyReportRequest{
		Filter: &dregrpc.Filter{
			Offset:     filterData.Offset,
			Limit:      limitPerPage,
			SortBy:     filterData.SortBy,
			Order:      filterData.Order,
			SearchTerm: filterData.SearchTerm,
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	drList := make([]DailyReport, 0, len(drlst.GetDre()))
	for _, item := range drlst.GetDre() {
		drData := DailyReport{
			ReportID:     item.ReportID,
			SerialNumber: item.SerialNumber,
			Date:         item.Date,
			Amount:       item.Amount,
			Currency:     item.Currency,
			CreatedAt:    item.CreatedAt.AsTime(),
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt.AsTime(),
			UpdatedBy:    item.UpdatedBy,
		}
		drList = append(drList, drData)
	}

	drstat, err := h.drc.DailyReportStats(r.Context(), &dregrpc.DailyReportStatsRequest{
		Filter: &dregrpc.Filter{
			Offset:     filterData.Offset,
			Limit:      limitPerPage,
			SortBy:     filterData.SortBy,
			Order:      filterData.Order,
			SearchTerm: filterData.SearchTerm,
		},
	})
	if err != nil {
		log.Println("unable to get stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	msg := map[string]string{}
	if filterData.SearchTerm != "" && len(drlst.GetDre()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(drlst.GetDre()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := DreTemplateData{
		FilterData:     *filterData,
		List:           drList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: dailyReportListPath,
	}
	if len(drList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, drstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewDailyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["report_id"]
	res, err := h.drc.GetDailyReport(r.Context(), &dregrpc.GetDailyReportRequest{
		Dre: &dregrpc.DailyReport{
			ReportID: id,
		},
	})
	if err != nil {
		log.Println("unable to get daily report info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DreTemplateData{
		Dre: DailyReport{
			ReportID:     res.Dre.ReportID,
			SerialNumber: res.Dre.SerialNumber,
			Date:         res.Dre.Date,
			Amount:       res.Dre.Amount,
			Currency:     res.Dre.Currency,
			CreatedAt:    res.Dre.CreatedAt.AsTime(),
			CreatedBy:    res.Dre.CreatedBy,
			UpdatedAt:    res.Dre.UpdatedAt.AsTime(),
			UpdatedBy:    res.Dre.UpdatedBy,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: dailyReportListPath,
	}

	err = h.templates.ExecuteTemplate(w, "dre-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteDailyReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["report_id"]
	if _, err := h.drc.DeleteDailyReport(r.Context(), &dregrpc.DeleteDailyReportRequest{
		Dre: &dregrpc.DailyReport{
			ReportID:  id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete collection: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, dailyReportListPath, http.StatusSeeOther)
}

func (h *Handler) loadDailyReportCreateForm(w http.ResponseWriter, dre DailyReport) {
	form := DreTemplateData{
		Dre:            dre,
		URLs:           listOfURLs(),
		CurrentPageURL: dailyReportListPath,
	}

	err := h.templates.ExecuteTemplate(w, "dre-create.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadDailyReportEditForm(w http.ResponseWriter, dre DailyReport) {
	form := DreTemplateData{
		Dre:            dre,
		URLs:           listOfURLs(),
		CurrentPageURL: dailyReportListPath,
	}

	err := h.templates.ExecuteTemplate(w, "dre-edit.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
