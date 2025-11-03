package handler

import (
	"encoding/json"
	"help-save-a-life/cms/paginator"
	collgrpc "help-save-a-life/proto/collection"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Collection struct {
	CollectionID  string
	SerialNumber  int32
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
	Coll           Collection
	List           []Collection
	Currencies     []Currency
	AccountTypes   []AccountType
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
	Title          string
}

type CollectionHome struct {
	Date          string
	AccountNumber string
	AccountType   string
	Amount        int32
	Currency      string
}

type CollJSONData struct {
	CollHome   []CollectionHome
	Paginator  paginator.Paginator
	FilterData Filter
}

func (c Collection) Validate(h *Handler) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AccountType,
			validation.Required.Error("Account type field can not be empty"),
		),
		validation.Field(&c.AccountNumber,
			validation.Required.Error("Account number field can not be empty"),
		),
		validation.Field(&c.Date,
			validation.Required.Error("Date field can not be empty"),
		),
		validation.Field(&c.Amount,
			validation.Required.Error("Amount field can not be empty"),
		),
		validation.Field(&c.Currency,
			validation.Required.Error("Currency field can not be empty"),
		),
	)
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	data := CollTemplateData{
		Coll:           Collection{},
		Currencies:     h.getCurrencyList(w, r),
		AccountTypes:   h.getAccountTypeList(w, r),
		URLs:           listOfURLs(),
		CurrentPageURL: collectionListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	h.loadCollectionCreateForm(w, data)
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

	if err := coll.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := CollTemplateData{
			Coll:           coll,
			Currencies:     h.getCurrencyList(w, r),
			AccountTypes:   h.getAccountTypeList(w, r),
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: collectionListPath,
			Title:          h.getSettingsTitle(w, r),
		}
		h.loadCollectionCreateForm(w, data)
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
			CreatedBy:     h.getLoggedUser(r),
			UpdatedBy:     h.getLoggedUser(r),
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
	res, err := h.cc.GetCollection(r.Context(), &collgrpc.GetCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: id,
		},
	})
	if err != nil {
		log.Println("unable to get collection info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadCollectionEditForm(w, CollTemplateData{
		Coll: Collection{
			CollectionID:  res.Coll.CollectionID,
			AccountType:   res.Coll.AccountType,
			AccountNumber: res.Coll.AccountNumber,
			Sender:        res.Coll.Sender,
			Date:          res.Coll.Date,
			Amount:        res.Coll.Amount,
			Currency:      res.Coll.Currency,
		},
		Currencies:     h.getCurrencyList(w, r),
		AccountTypes:   h.getAccountTypeList(w, r),
		URLs:           listOfURLs(),
		CurrentPageURL: collectionListPath,
		Title:          h.getSettingsTitle(w, r),
	})
}

func (h *Handler) updateCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["collection_id"]
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
	coll.CollectionID = id

	if err := coll.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := CollTemplateData{
			Coll:           coll,
			Currencies:     h.getCurrencyList(w, r),
			AccountTypes:   h.getAccountTypeList(w, r),
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: collectionListPath,
			Title:          h.getSettingsTitle(w, r),
		}
		h.loadCollectionEditForm(w, data)
		return
	}

	if _, err := h.cc.UpdateCollection(ctx, &collgrpc.UpdateCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID:  id,
			AccountType:   coll.AccountType,
			AccountNumber: coll.AccountNumber,
			Sender:        coll.Sender,
			Date:          coll.Date,
			Amount:        coll.Amount,
			Currency:      coll.Currency,
			UpdatedBy:     h.getLoggedUser(r),
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

	currencyList := h.getCurrencyListMap(w, r)
	acctList := h.getAccountTypeListMap(w, r)
	collList := make([]Collection, 0, len(clst.GetColl()))
	for _, item := range clst.GetColl() {
		cData := Collection{
			CollectionID:  item.CollectionID,
			SerialNumber:  item.SerialNumber,
			AccountType:   acctList[item.AccountType],
			AccountNumber: item.AccountNumber,
			Sender:        item.Sender,
			Date:          item.Date,
			Amount:        item.Amount,
			Currency:      currencyList[item.Currency],
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
		FilterData:     *filterData,
		List:           collList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: collectionListPath,
		Title:          h.getSettingsTitle(w, r),
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
	res, err := h.cc.GetCollection(r.Context(), &collgrpc.GetCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: id,
		},
	})
	if err != nil {
		log.Println("unable to get collection info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := h.getUserListMap(w, r)
	data := CollTemplateData{
		Coll: Collection{
			CollectionID:  res.Coll.CollectionID,
			SerialNumber:  res.Coll.SerialNumber,
			AccountType:   h.getAccountTypeListMap(w, r)[res.Coll.AccountType],
			AccountNumber: res.Coll.AccountNumber,
			Sender:        res.Coll.Sender,
			Date:          res.Coll.Date,
			Amount:        res.Coll.Amount,
			Currency:      h.getCurrencyListMap(w, r)[res.Coll.Currency],
			CreatedAt:     res.Coll.CreatedAt.AsTime(),
			CreatedBy:     users[res.Coll.CreatedBy],
			UpdatedAt:     res.Coll.UpdatedAt.AsTime(),
			UpdatedBy:     users[res.Coll.UpdatedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: collectionListPath,
		Title:          h.getSettingsTitle(w, r),
	}

	err = h.templates.ExecuteTemplate(w, "coll-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["collection_id"]
	if _, err := h.cc.DeleteCollection(r.Context(), &collgrpc.DeleteCollectionRequest{
		Coll: &collgrpc.Collection{
			CollectionID: id,
			DeletedBy:    h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete collection: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, collectionListPath, http.StatusSeeOther)
}

func (h *Handler) loadCollectionCreateForm(w http.ResponseWriter, data CollTemplateData) {
	err := h.templates.ExecuteTemplate(w, "coll-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCollectionEditForm(w http.ResponseWriter, data CollTemplateData) {
	err := h.templates.ExecuteTemplate(w, "coll-edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) viewCollectionHome(w http.ResponseWriter, r *http.Request) {
	acctList := h.getAccountTypeListMap(w, r)
	currencyList := h.getCurrencyListMap(w, r)
	filterData := GetFilterData(r)
	clst, err := h.cc.ListCollection(r.Context(), &collgrpc.ListCollectionRequest{
		Filter: &collgrpc.Filter{
			Offset: filterData.Offset,
			Limit:  limitPerPage,
			SortBy: filterData.SortBy,
			Order:  filterData.Order,
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	collList := make([]CollectionHome, 0, len(clst.GetColl()))

	for _, item := range clst.GetColl() {
		cData := CollectionHome{
			AccountType:   acctList[item.AccountType],
			AccountNumber: hideDigits(item.AccountNumber),
			Date:          item.Date,
			Amount:        item.Amount,
			Currency:      currencyList[item.Currency],
		}
		collList = append(collList, cData)
	}

	collstat, err := h.cc.CollectionStats(r.Context(), &collgrpc.CollectionStatsRequest{
		Filter: &collgrpc.Filter{
			Offset: filterData.Offset,
			Limit:  limitPerPage,
			SortBy: filterData.SortBy,
			Order:  filterData.Order,
		},
	})
	if err != nil {
		log.Println("unable to get stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	data := CollJSONData{
		FilterData: *filterData,
		CollHome:   collList,
	}

	if len(collList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, collstat.Stats.Count, r)
	}

	json.NewEncoder(w).Encode(data)
}
