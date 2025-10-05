package handler

import (
	"help-save-a-life/cms/paginator"
	currgrpc "help-save-a-life/proto/currency"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Currency struct {
	ID           string
	SerialNumber int32
	Name         string
	ExchangeRate int32
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    time.Time
	DeletedBy    string
}

type CurrencyTemplateData struct {
	Curr           Currency
	List           []Currency
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
}

func (c Currency) Validate(h *Handler) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name,
			validation.Required.Error("Name field can not be empty"),
			validation.Length(1, 20).Error("Name field can not contain more than 20 characters"),
		),
		validation.Field(&c.ExchangeRate,
			validation.Required.Error("Exchange rate field can not be empty"),
		),
	)
}

func (h *Handler) createCurrency(w http.ResponseWriter, r *http.Request) {
	data := CurrencyTemplateData{
		Curr:           Currency{},
		URLs:           listOfURLs(),
		CurrentPageURL: currencyListPath,
	}
	h.loadCurrencyCreateForm(w, data)
}

func (h *Handler) storeCurrency(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var curr Currency
	err = h.decoder.Decode(&curr, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := curr.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := CurrencyTemplateData{
			Curr:           curr,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: collectionListPath,
		}
		h.loadCurrencyCreateForm(w, data)
		return
	}

	_, err = h.cuc.CreateCurrency(r.Context(), &currgrpc.CreateCurrencyRequest{
		Curr: &currgrpc.Currency{
			Name:         curr.Name,
			ExchangeRate: curr.ExchangeRate,
			CreatedBy:    h.getLoggedUser(r),
			UpdatedBy:    h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, currencyListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.cuc.GetCurrency(r.Context(), &currgrpc.GetCurrencyRequest{
		Curr: &currgrpc.Currency{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get currency info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadCurrencyEditForm(w, CurrencyTemplateData{
		Curr: Currency{
			ID:           res.Curr.ID,
			Name:         res.Curr.Name,
			ExchangeRate: res.Curr.ExchangeRate,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: currencyListPath,
	})
}

func (h *Handler) updateCurrency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var curr Currency
	if err := h.decoder.Decode(&curr, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	curr.ID = id

	if err := curr.Validate(h); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for k, v := range e {
					vErrs[k] = v.Error()
				}
			}
		}
		h.loadCurrencyEditForm(w, CurrencyTemplateData{
			Curr:           curr,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: currencyListPath,
		})
		return
	}

	if _, err := h.cuc.UpdateCurrency(ctx, &currgrpc.UpdateCurrencyRequest{
		Curr: &currgrpc.Currency{
			ID:           id,
			Name:         curr.Name,
			ExchangeRate: curr.ExchangeRate,
			UpdatedBy:    h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, currencyListPath, http.StatusSeeOther)
}

func (h *Handler) listCurrency(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("curr-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	clst, err := h.cuc.ListCurrency(r.Context(), &currgrpc.ListCurrencyRequest{
		Filter: &currgrpc.Filter{
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

	currList := make([]Currency, 0, len(clst.GetCurr()))
	for _, item := range clst.GetCurr() {
		cData := Currency{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			Name:         item.Name,
			ExchangeRate: item.ExchangeRate,
			CreatedAt:    item.CreatedAt.AsTime(),
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt.AsTime(),
			UpdatedBy:    item.UpdatedBy,
		}
		currList = append(currList, cData)
	}

	currstat, err := h.cuc.CurrencyStats(r.Context(), &currgrpc.CurrencyStatsRequest{
		Filter: &currgrpc.Filter{
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
	if filterData.SearchTerm != "" && len(clst.GetCurr()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(clst.GetCurr()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := CurrencyTemplateData{
		FilterData:     *filterData,
		List:           currList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: currencyListPath,
	}
	if len(currList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, currstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.cuc.GetCurrency(r.Context(), &currgrpc.GetCurrencyRequest{
		Curr: &currgrpc.Currency{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get currency info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := h.getUserListMap(w, r)
	data := CurrencyTemplateData{
		Curr: Currency{
			ID:           res.Curr.ID,
			SerialNumber: res.Curr.SerialNumber,
			Name:         res.Curr.Name,
			ExchangeRate: res.Curr.ExchangeRate,
			CreatedAt:    res.Curr.CreatedAt.AsTime(),
			CreatedBy:    users[res.Curr.CreatedBy],
			UpdatedAt:    res.Curr.UpdatedAt.AsTime(),
			UpdatedBy:    users[res.Curr.UpdatedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: currencyListPath,
	}

	err = h.templates.ExecuteTemplate(w, "curr-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, err := h.cuc.DeleteCurrency(r.Context(), &currgrpc.DeleteCurrencyRequest{
		Curr: &currgrpc.Currency{
			ID:        id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete currency: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, currencyListPath, http.StatusSeeOther)
}

func (h *Handler) loadCurrencyCreateForm(w http.ResponseWriter, data CurrencyTemplateData) {
	err := h.templates.ExecuteTemplate(w, "curr-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCurrencyEditForm(w http.ResponseWriter, data CurrencyTemplateData) {
	err := h.templates.ExecuteTemplate(w, "curr-edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getCurrencyList(w http.ResponseWriter, r *http.Request) []Currency {
	clst, err := h.cuc.ListCurrency(r.Context(), &currgrpc.ListCurrencyRequest{
		Filter: &currgrpc.Filter{
			SortBy: "name",
			Order:  "ASC",
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	currList := make([]Currency, 0, len(clst.GetCurr()))
	for _, item := range clst.GetCurr() {
		cData := Currency{
			ID:   item.ID,
			Name: item.Name,
		}
		currList = append(currList, cData)
	}
	return currList
}

func (h *Handler) getCurrencyListMap(w http.ResponseWriter, r *http.Request) map[string]string {
	clst, err := h.cuc.ListCurrency(r.Context(), &currgrpc.ListCurrencyRequest{
		Filter: &currgrpc.Filter{
			SortBy: "name",
			Order:  "ASC",
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	currList := make(map[string]string, len(clst.GetCurr()))
	for _, item := range clst.GetCurr() {
		currList[item.ID] = item.Name
	}
	return currList
}
