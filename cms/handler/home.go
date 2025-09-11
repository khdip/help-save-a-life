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
	List            []Collection
	Paginator       paginator.Paginator
	FilterData      Filter
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filterData := GetFilterData(r)
	clst, err := h.cc.ListCollection(r.Context(), &collgrpc.ListCollectionRequest{
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
	for _, item := range clst.GetColl() {
		cData := Collection{
			CollectionID:  item.CollectionID,
			AccountType:   item.AccountType,
			AccountNumber: hideDigits(item.AccountNumber),
			Sender:        item.Sender,
			Date:          item.Date,
			Amount:        item.Amount,
			Currency:      item.Currency,
			CreatedAt:     item.CreatedAt.AsTime(),
			CreatedBy:     item.CreatedBy,
			UpdatedAt:     item.UpdatedAt.AsTime(),
			UpdatedBy:     item.UpdatedBy,
		}
		collList = append(collList, cData)
	}

	collstat, err := h.cc.CollectionStats(r.Context(), &collgrpc.CollectionStatsRequest{
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

	totalCollection := drstat.Stats.TotalAmount
	if totalCollection > targetamount {
		totalCollection = targetamount
	}

	data := HomeTemplateData{
		List:            collList,
		FilterData:      *filterData,
		URLs:            listOfURLs(),
		TargetAmount:    formatWithCommas(targetamount),
		CollectedAmount: formatWithCommas(totalCollection),
		RemainingAmount: formatWithCommas(targetamount - totalCollection),
		Percentage:      fmt.Sprintf("%.2f", ((float64(totalCollection) / float64(targetamount)) * 100)),
	}
	if len(collList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, collstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}
