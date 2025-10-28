package handler

import (
	"help-save-a-life/cms/paginator"
	acctgrpc "help-save-a-life/proto/accountType"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type AccountType struct {
	ID           string
	SerialNumber int32
	Title        string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    time.Time
	DeletedBy    string
}

type AccountTypeTemplateData struct {
	Acct           AccountType
	List           []AccountType
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
}

func (a AccountType) Validate(h *Handler) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title,
			validation.Required.Error("Title field can not be empty"),
			validation.Length(1, 100).Error("Title field can not contain more than 100 characters"),
		),
	)
}

func (h *Handler) createAccountType(w http.ResponseWriter, r *http.Request) {
	data := AccountTypeTemplateData{
		Acct:           AccountType{},
		URLs:           listOfURLs(),
		CurrentPageURL: accountTypeListPath,
	}
	h.loadAccountTypeCreateForm(w, data)
}

func (h *Handler) storeAccountType(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var acct AccountType
	err = h.decoder.Decode(&acct, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := acct.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := AccountTypeTemplateData{
			Acct:           acct,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: accountTypeListPath,
		}
		h.loadAccountTypeCreateForm(w, data)
		return
	}

	_, err = h.at.CreateAccountType(r.Context(), &acctgrpc.CreateAccountTypeRequest{
		Acct: &acctgrpc.AccountType{
			Title:     acct.Title,
			CreatedBy: h.getLoggedUser(r),
			UpdatedBy: h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, accountTypeListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editAccountType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.at.GetAccountType(r.Context(), &acctgrpc.GetAccountTypeRequest{
		Acct: &acctgrpc.AccountType{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get account type info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadAccountTypeEditForm(w, AccountTypeTemplateData{
		Acct: AccountType{
			ID:    res.Acct.ID,
			Title: res.Acct.Title,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: accountTypeListPath,
	})
}

func (h *Handler) updateAccountType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var acct AccountType
	if err := h.decoder.Decode(&acct, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acct.ID = id

	if err := acct.Validate(h); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for k, v := range e {
					vErrs[k] = v.Error()
				}
			}
		}
		h.loadAccountTypeEditForm(w, AccountTypeTemplateData{
			Acct:           acct,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: accountTypeListPath,
		})
		return
	}

	if _, err := h.at.UpdateAccountType(ctx, &acctgrpc.UpdateAccountTypeRequest{
		Acct: &acctgrpc.AccountType{
			ID:        id,
			Title:     acct.Title,
			UpdatedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, accountTypeListPath, http.StatusSeeOther)
}

func (h *Handler) listAccountType(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("acct-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	alst, err := h.at.ListAccountType(r.Context(), &acctgrpc.ListAccountTypeRequest{
		Filter: &acctgrpc.Filter{
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

	acctList := make([]AccountType, 0, len(alst.GetAcct()))
	for _, item := range alst.GetAcct() {
		aData := AccountType{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			Title:        item.Title,
			CreatedAt:    item.CreatedAt.AsTime(),
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt.AsTime(),
			UpdatedBy:    item.UpdatedBy,
		}
		acctList = append(acctList, aData)
	}

	acctstat, err := h.at.AccountTypeStats(r.Context(), &acctgrpc.AccountTypeStatsRequest{
		Filter: &acctgrpc.Filter{
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
	if filterData.SearchTerm != "" && len(alst.GetAcct()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(alst.GetAcct()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := AccountTypeTemplateData{
		FilterData:     *filterData,
		List:           acctList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: accountTypeListPath,
	}
	if len(acctList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, acctstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewAccountType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.at.GetAccountType(r.Context(), &acctgrpc.GetAccountTypeRequest{
		Acct: &acctgrpc.AccountType{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get account type info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := h.getUserListMap(w, r)
	data := AccountTypeTemplateData{
		Acct: AccountType{
			ID:           res.Acct.ID,
			SerialNumber: res.Acct.SerialNumber,
			Title:        res.Acct.Title,
			CreatedAt:    res.Acct.CreatedAt.AsTime(),
			CreatedBy:    users[res.Acct.CreatedBy],
			UpdatedAt:    res.Acct.UpdatedAt.AsTime(),
			UpdatedBy:    users[res.Acct.UpdatedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: accountTypeListPath,
	}

	err = h.templates.ExecuteTemplate(w, "acct-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteAccountType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, err := h.at.DeleteAccountType(r.Context(), &acctgrpc.DeleteAccountTypeRequest{
		Acct: &acctgrpc.AccountType{
			ID:        id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete account type: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, accountTypeListPath, http.StatusSeeOther)
}

func (h *Handler) loadAccountTypeCreateForm(w http.ResponseWriter, data AccountTypeTemplateData) {
	err := h.templates.ExecuteTemplate(w, "acct-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadAccountTypeEditForm(w http.ResponseWriter, data AccountTypeTemplateData) {
	err := h.templates.ExecuteTemplate(w, "acct-edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAccountTypeList(w http.ResponseWriter, r *http.Request) []AccountType {
	alst, err := h.at.ListAccountType(r.Context(), &acctgrpc.ListAccountTypeRequest{
		Filter: &acctgrpc.Filter{
			SortBy: "title",
			Order:  "ASC",
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	acctList := make([]AccountType, 0, len(alst.GetAcct()))
	for _, item := range alst.GetAcct() {
		aData := AccountType{
			ID:    item.ID,
			Title: item.Title,
		}
		acctList = append(acctList, aData)
	}
	return acctList
}

func (h *Handler) getAccountTypeListMap(w http.ResponseWriter, r *http.Request) map[string]string {
	alst, err := h.at.ListAccountType(r.Context(), &acctgrpc.ListAccountTypeRequest{
		Filter: &acctgrpc.Filter{
			SortBy: "title",
			Order:  "ASC",
		},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	acctList := make(map[string]string, len(alst.GetAcct()))
	for _, item := range alst.GetAcct() {
		acctList[item.ID] = item.Title
	}
	return acctList
}
