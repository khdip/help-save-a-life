package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/khdip/help-save-a-life/cms/paginator"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Accounts struct {
	ID           string
	SerialNumber int32
	AccountType  string
	Details      string
	DetailsT     template.HTML
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    time.Time
	DeletedBy    string
}

type AccountsTemplateData struct {
	Acnt           Accounts
	List           []Accounts
	AccountTypes   []AccountType
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
	Title          string
}

func (a Accounts) Validate(h *Handler) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AccountType,
			validation.Required.Error("Title field can not be empty"),
			validation.Length(1, 100).Error("Title field can not contain more than 100 characters"),
		),
		validation.Field(&a.Details,
			validation.Required.Error("Title field can not be empty"),
		),
	)
}

func (h *Handler) createAccounts(w http.ResponseWriter, r *http.Request) {
	data := AccountsTemplateData{
		Acnt:           Accounts{},
		AccountTypes:   h.getAccountTypeList(w, r),
		URLs:           listOfURLs(),
		CurrentPageURL: accountsListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	h.loadAccountsCreateForm(w, data)
}

func (h *Handler) storeAccounts(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var acnt Accounts
	err = h.decoder.Decode(&acnt, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := acnt.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := AccountsTemplateData{
			Acnt:           acnt,
			AccountTypes:   h.getAccountTypeList(w, r),
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: accountsListPath,
			Title:          h.getSettingsTitle(w, r),
		}
		h.loadAccountsCreateForm(w, data)
		return
	}

	_, err = h.acc.CreateAccounts(r.Context(), &acntgrpc.CreateAccountsRequest{
		Acnt: &acntgrpc.Accounts{
			AccountType: acnt.AccountType,
			Details:     acnt.Details,
			CreatedBy:   h.getLoggedUser(r),
			UpdatedBy:   h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, accountsListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editAccounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.acc.GetAccounts(r.Context(), &acntgrpc.GetAccountsRequest{
		Acnt: &acntgrpc.Accounts{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get accounts info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadAccountsEditForm(w, AccountsTemplateData{
		Acnt: Accounts{
			ID:          res.Acnt.ID,
			AccountType: res.Acnt.AccountType,
			Details:     res.Acnt.Details,
		},
		AccountTypes:   h.getAccountTypeList(w, r),
		URLs:           listOfURLs(),
		CurrentPageURL: accountsListPath,
		Title:          h.getSettingsTitle(w, r),
	})
}

func (h *Handler) updateAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var acnt Accounts
	if err := h.decoder.Decode(&acnt, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acnt.ID = id

	if err := acnt.Validate(h); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for k, v := range e {
					vErrs[k] = v.Error()
				}
			}
		}
		h.loadAccountsEditForm(w, AccountsTemplateData{
			Acnt:           acnt,
			AccountTypes:   h.getAccountTypeList(w, r),
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: accountsListPath,
			Title:          h.getSettingsTitle(w, r),
		})
		return
	}

	if _, err := h.acc.UpdateAccounts(ctx, &acntgrpc.UpdateAccountsRequest{
		Acnt: &acntgrpc.Accounts{
			ID:          id,
			AccountType: acnt.AccountType,
			Details:     acnt.Details,
			UpdatedBy:   h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, accountsListPath, http.StatusSeeOther)
}

func (h *Handler) listAccounts(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("acnt-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	acct := h.getAccountTypeListMap(w, r)
	alst, err := h.acc.ListAccounts(r.Context(), &acntgrpc.ListAccountsRequest{
		Filter: &acntgrpc.Filter{
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

	acntList := make([]Accounts, 0, len(alst.GetAcnt()))
	for _, item := range alst.GetAcnt() {
		aData := Accounts{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			AccountType:  acct[item.AccountType],
			DetailsT:     makeHTMLTemplate(item.Details),
			CreatedAt:    item.CreatedAt.AsTime(),
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt.AsTime(),
			UpdatedBy:    item.UpdatedBy,
		}
		acntList = append(acntList, aData)
	}

	acntstat, err := h.acc.AccountsStats(r.Context(), &acntgrpc.AccountsStatsRequest{
		Filter: &acntgrpc.Filter{
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
	if filterData.SearchTerm != "" && len(alst.GetAcnt()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(alst.GetAcnt()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := AccountsTemplateData{
		FilterData:     *filterData,
		List:           acntList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: accountsListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	if len(acntList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, acntstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewAccounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.acc.GetAccounts(r.Context(), &acntgrpc.GetAccountsRequest{
		Acnt: &acntgrpc.Accounts{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get accounts info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := h.getUserListMap(w, r)
	data := AccountsTemplateData{
		Acnt: Accounts{
			ID:           res.Acnt.ID,
			SerialNumber: res.Acnt.SerialNumber,
			AccountType:  h.getAccountTypeListMap(w, r)[res.Acnt.AccountType],
			DetailsT:     template.HTML(res.Acnt.Details),
			CreatedAt:    res.Acnt.CreatedAt.AsTime(),
			CreatedBy:    users[res.Acnt.CreatedBy],
			UpdatedAt:    res.Acnt.UpdatedAt.AsTime(),
			UpdatedBy:    users[res.Acnt.UpdatedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: accountsListPath,
		Title:          h.getSettingsTitle(w, r),
	}

	err = h.templates.ExecuteTemplate(w, "acnt-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteAccounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, err := h.acc.DeleteAccounts(r.Context(), &acntgrpc.DeleteAccountsRequest{
		Acnt: &acntgrpc.Accounts{
			ID:        id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete accounts: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, accountsListPath, http.StatusSeeOther)
}

func (h *Handler) loadAccountsCreateForm(w http.ResponseWriter, data AccountsTemplateData) {
	err := h.templates.ExecuteTemplate(w, "acnt-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadAccountsEditForm(w http.ResponseWriter, data AccountsTemplateData) {
	err := h.templates.ExecuteTemplate(w, "acnt-edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAccounts(w http.ResponseWriter, r *http.Request) map[string][]Accounts {
	alst, err := h.acc.ListAccounts(r.Context(), &acntgrpc.ListAccountsRequest{
		Filter: &acntgrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	acctList := h.getAccountTypeListMap(w, r)

	acntList := make(map[string][]Accounts, len(alst.GetAcnt()))

	for _, item := range alst.GetAcnt() {
		acntList[acctList[item.AccountType]] = append(acntList[acctList[item.AccountType]], Accounts{DetailsT: template.HTML(item.Details)})
	}

	return acntList
}
