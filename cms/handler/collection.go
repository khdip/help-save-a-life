package handler

import (
	"help-save-a-life/cms/paginator"
	collgrpc "help-save-a-life/proto/collection"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Collection struct {
	CollectionID  int32
	AccountType   string
	AccountNumber string
	Sender        string
	Date          string
	Amount        int32
	Currency      string
	CreatedAt     time.Time
	CreatedBy     string
	UpdatedAt     time.Time
	UpdatedBy     string
	DeletedAt     time.Time
	DeletedBy     string
}

type CollTemplateData struct {
	Coll       Collection
	List       []Collection
	Paginator  paginator.Paginator
	FilterData Filter
	URLs       map[string]string
	Message    map[string]string
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	h.loadCollectionCreateForm(w, Collection{})
}

func (h *Handler) storeCollection(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var coll Collection
	err = h.decoder.Decode(&coll, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.cc.CreateCollection(r.Context(), &collgrpc.CreateCollectionRequest{
		Coll: &collgrpc.Collection{
			AccountType:   coll.AccountType,
			AccountNumber: coll.AccountNumber,
			Sender:        coll.Sender,
			Date:          coll.Date,
			Amount:        coll.Amount,
			Currency:      coll.Currency,
			CreatedBy:     h.getLoggedUser(w, r),
			UpdatedBy:     h.getLoggedUser(w, r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, collectionListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["collection_id"]
	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.cc.GetCollection(r.Context(), &collgrpc.GetCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: int32(cid),
		},
	})
	if err != nil {
		log.Println("unable to get collection info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadCollectionEditForm(w, Collection{
		CollectionID:  res.Coll.CollectionID,
		AccountType:   res.Coll.AccountType,
		AccountNumber: res.Coll.AccountNumber,
		Sender:        res.Coll.Sender,
		Date:          res.Coll.Date,
		Amount:        res.Coll.Amount,
		Currency:      res.Coll.Currency,
	})
}

func (h *Handler) updateCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["collection_id"]
	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var coll Collection
	if err := h.decoder.Decode(&coll, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.cc.UpdateCollection(ctx, &collgrpc.UpdateCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID:  int32(cid),
			AccountType:   coll.AccountType,
			AccountNumber: coll.AccountNumber,
			Sender:        coll.Sender,
			Date:          coll.Date,
			Amount:        coll.Amount,
			Currency:      coll.Currency,
			UpdatedBy:     h.getLoggedUser(w, r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, collectionListPath, http.StatusSeeOther)
}

func (h *Handler) listCollection(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("coll-list.html")
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
			AccountNumber: item.AccountNumber,
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

	msg := map[string]string{}
	if filterData.SearchTerm != "" && len(clst.GetColl()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(clst.GetColl()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := CollTemplateData{
		FilterData: *filterData,
		List:       collList,
		Message:    msg,
		URLs:       listOfURLs(),
	}
	if len(collList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, collstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["collection_id"]
	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.cc.GetCollection(r.Context(), &collgrpc.GetCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: int32(cid),
		},
	})
	if err != nil {
		log.Println("unable to get collection info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CollTemplateData{
		Coll: Collection{
			CollectionID:  res.Coll.CollectionID,
			AccountType:   res.Coll.AccountType,
			AccountNumber: res.Coll.AccountNumber,
			Sender:        res.Coll.Sender,
			Date:          res.Coll.Date,
			Amount:        res.Coll.Amount,
			Currency:      res.Coll.Currency,
			CreatedAt:     res.Coll.CreatedAt.AsTime(),
			CreatedBy:     res.Coll.CreatedBy,
			UpdatedAt:     res.Coll.UpdatedAt.AsTime(),
			UpdatedBy:     res.Coll.UpdatedBy,
		},
		URLs: listOfURLs(),
	}

	err = h.templates.ExecuteTemplate(w, "coll-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["collection_id"]
	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.cc.DeleteCollection(ctx, &collgrpc.DeleteCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: int32(cid),
			DeletedBy:    h.getLoggedUser(w, r),
		},
	}); err != nil {
		log.Println("unable to delete collection: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, collectionListPath, http.StatusSeeOther)
}

func (h *Handler) loadCollectionCreateForm(w http.ResponseWriter, coll Collection) {
	form := CollTemplateData{
		Coll: coll,
		URLs: listOfURLs(),
	}

	err := h.templates.ExecuteTemplate(w, "coll-create.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCollectionEditForm(w http.ResponseWriter, coll Collection) {
	form := CollTemplateData{
		Coll: coll,
		URLs: listOfURLs(),
	}

	err := h.templates.ExecuteTemplate(w, "coll-edit.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
