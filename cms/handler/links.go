package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/khdip/help-save-a-life/cms/paginator"
	linkgrpc "github.com/khdip/help-save-a-life/proto/links"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Link struct {
	ID           string
	SerialNumber int32
	LinkText     string
	LinkURL      string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    time.Time
	DeletedBy    string
}

type LinkTemplateData struct {
	Link           Link
	List           []Link
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
	Title          string
}

func (l Link) Validate(h *Handler) error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.LinkText,
			validation.Required.Error("Title field can not be empty"),
			validation.Length(1, 100).Error("Title field can not contain more than 100 characters"),
		),
		validation.Field(&l.LinkURL,
			validation.Required.Error("Title field can not be empty"),
		),
	)
}

func (h *Handler) createLink(w http.ResponseWriter, r *http.Request) {
	data := LinkTemplateData{
		Link:           Link{},
		URLs:           listOfURLs(),
		CurrentPageURL: linkListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	h.loadLinkCreateForm(w, data)
}

func (h *Handler) storeLink(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var link Link
	err = h.decoder.Decode(&link, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := link.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := LinkTemplateData{
			Link:           link,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: linkListPath,
			Title:          h.getSettingsTitle(w, r),
		}
		h.loadLinkCreateForm(w, data)
		return
	}

	_, err = h.lnk.CreateLink(r.Context(), &linkgrpc.CreateLinkRequest{
		Link: &linkgrpc.Link{
			LinkText:  link.LinkText,
			LinkURL:   link.LinkURL,
			CreatedBy: h.getLoggedUser(r),
			UpdatedBy: h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, linkListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editLink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.lnk.GetLink(r.Context(), &linkgrpc.GetLinkRequest{
		Link: &linkgrpc.Link{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get link info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadLinkEditForm(w, LinkTemplateData{
		Link: Link{
			ID:       res.Link.ID,
			LinkText: res.Link.LinkText,
			LinkURL:  res.Link.LinkURL,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: linkListPath,
		Title:          h.getSettingsTitle(w, r),
	})
}

func (h *Handler) updateLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var link Link
	if err := h.decoder.Decode(&link, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	link.ID = id

	if err := link.Validate(h); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for k, v := range e {
					vErrs[k] = v.Error()
				}
			}
		}
		h.loadLinkEditForm(w, LinkTemplateData{
			Link:           link,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: linkListPath,
			Title:          h.getSettingsTitle(w, r),
		})
		return
	}

	if _, err := h.lnk.UpdateLink(ctx, &linkgrpc.UpdateLinkRequest{
		Link: &linkgrpc.Link{
			ID:        id,
			LinkText:  link.LinkText,
			LinkURL:   link.LinkURL,
			UpdatedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, linkListPath, http.StatusSeeOther)
}

func (h *Handler) listLink(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("link-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	llst, err := h.lnk.ListLink(r.Context(), &linkgrpc.ListLinkRequest{
		Filter: &linkgrpc.Filter{
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

	linkList := make([]Link, 0, len(llst.GetLink()))
	for _, item := range llst.GetLink() {
		lData := Link{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			LinkText:     item.LinkText,
			LinkURL:      item.LinkURL,
			CreatedAt:    item.CreatedAt.AsTime(),
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt.AsTime(),
			UpdatedBy:    item.UpdatedBy,
		}
		linkList = append(linkList, lData)
	}

	linkstat, err := h.lnk.LinkStats(r.Context(), &linkgrpc.LinkStatsRequest{
		Filter: &linkgrpc.Filter{
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
	if filterData.SearchTerm != "" && len(llst.GetLink()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(llst.GetLink()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := LinkTemplateData{
		FilterData:     *filterData,
		List:           linkList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: linkListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	if len(linkList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, linkstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewLink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.lnk.GetLink(r.Context(), &linkgrpc.GetLinkRequest{
		Link: &linkgrpc.Link{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get link info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := h.getUserListMap(w, r)
	data := LinkTemplateData{
		Link: Link{
			ID:           res.Link.ID,
			SerialNumber: res.Link.SerialNumber,
			LinkText:     res.Link.LinkText,
			LinkURL:      res.Link.LinkURL,
			CreatedAt:    res.Link.CreatedAt.AsTime(),
			CreatedBy:    users[res.Link.CreatedBy],
			UpdatedAt:    res.Link.UpdatedAt.AsTime(),
			UpdatedBy:    users[res.Link.UpdatedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: linkListPath,
		Title:          h.getSettingsTitle(w, r),
	}

	err = h.templates.ExecuteTemplate(w, "link-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteLink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, err := h.lnk.DeleteLink(r.Context(), &linkgrpc.DeleteLinkRequest{
		Link: &linkgrpc.Link{
			ID:        id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete link: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, linkListPath, http.StatusSeeOther)
}

func (h *Handler) loadLinkCreateForm(w http.ResponseWriter, data LinkTemplateData) {
	err := h.templates.ExecuteTemplate(w, "link-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadLinkEditForm(w http.ResponseWriter, data LinkTemplateData) {
	err := h.templates.ExecuteTemplate(w, "link-edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getLinkList(w http.ResponseWriter, r *http.Request) []Link {
	llst, err := h.lnk.ListLink(r.Context(), &linkgrpc.ListLinkRequest{
		Filter: &linkgrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	linkList := make([]Link, 0, len(llst.GetLink()))
	for _, item := range llst.GetLink() {
		lData := Link{
			LinkText: item.LinkText,
			LinkURL:  item.LinkURL,
		}
		linkList = append(linkList, lData)
	}
	return linkList
}
