package handler

import (
	"fmt"
	"help-save-a-life/cms/paginator"
	collgrpc "help-save-a-life/proto/collection"
	dregrpc "help-save-a-life/proto/dailyReport"
	"log"
	"net/http"
)

type HomeTemplateData struct {
	CollList        []Collection
	DreList         []DailyReport
	Sett            SettingsHome
	Paginator       paginator.Paginator
	FilterData      Filter
	TargetAmount    string
	CollectedAmount string
	RemainingAmount string
	Percentage      string
	URLs            map[string]string
}

func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	template := h.templates.Lookup("index.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	currencyList := h.getCurrencyListMap(w, r)
	acctList := h.getAccountTypeListMap(w, r)
	sett := h.getSettingsHome(w, r)
	filterData := GetFilterData(r)
	clst, err := h.cc.ListCollection(ctx, &collgrpc.ListCollectionRequest{
		Filter: &collgrpc.Filter{
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

	collList := make([]Collection, 0, len(clst.GetColl()))
	if sett.ShowCollection {
		for _, item := range clst.GetColl() {
			cData := Collection{
				AccountType:   acctList[item.AccountType],
				AccountNumber: hideDigits(item.AccountNumber),
				Date:          item.Date,
				Amount:        item.Amount,
				Currency:      currencyList[item.Currency],
			}
			collList = append(collList, cData)
		}
	}

	collstat, err := h.cc.CollectionStats(ctx, &collgrpc.CollectionStatsRequest{
		Filter: &collgrpc.Filter{
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

	drlst, err := h.drc.ListDailyReport(ctx, &dregrpc.ListDailyReportRequest{
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
	if sett.ShowDailyReport {
		for _, item := range drlst.GetDre() {
			drData := DailyReport{
				ReportID:     item.ReportID,
				SerialNumber: item.SerialNumber,
				Date:         item.Date,
				Amount:       item.Amount,
				Currency:     currencyList[item.Currency],
				CreatedAt:    item.CreatedAt.AsTime(),
				CreatedBy:    item.CreatedBy,
				UpdatedAt:    item.UpdatedAt.AsTime(),
				UpdatedBy:    item.UpdatedBy,
			}
			drList = append(drList, drData)
		}
	}

	drstat, err := h.drc.DailyReportStats(ctx, &dregrpc.DailyReportStatsRequest{
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

	var totalCollection int32
	switch sett.CalculateCollection {
	case 0:
		totalCollection = collstat.Stats.TotalAmount
	case 1:
		totalCollection = drstat.Stats.TotalAmount
	case 2:
		totalCollection = sett.TotalAmount
	}

	targetAmount := sett.TargetAmount
	if totalCollection > targetAmount {
		totalCollection = targetAmount
	}

	data := HomeTemplateData{
		CollList:        collList,
		DreList:         drList,
		Sett:            sett,
		FilterData:      *filterData,
		URLs:            listOfURLs(),
		TargetAmount:    formatWithCommas(targetAmount),
		CollectedAmount: formatWithCommas(totalCollection),
		RemainingAmount: formatWithCommas(targetAmount - totalCollection),
		Percentage:      fmt.Sprintf("%.2f", ((float64(totalCollection) / float64(targetAmount)) * 100)),
	}
	if len(collList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, collstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
		return
	}
}
