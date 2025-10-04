package handler

import (
	"encoding/json"
	"help-save-a-life/cms/paginator"
	commgrpc "help-save-a-life/proto/comments"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

type Comment struct {
	CommentID    string
	SerialNumber int32
	Name         string
	Email        string
	Comment      string
	CreatedAt    time.Time
}

type CommTemplateData struct {
	Comm           Comment
	List           []Comment
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
}

func (c Comment) Validate(h *Handler) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name,
			validation.Required.Error("Name field can not be empty"),
			validation.Length(1, 50).Error("Name field can not contain more than 50 characters"),
		),
		validation.Field(&c.Email,
			validation.Required.Error("Email field can not be empty"),
			is.Email,
		),
		validation.Field(&c.Comment,
			validation.Required.Error("Comment field can not be empty"),
		),
	)
}

func (h *Handler) storeComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var comm Comment
	err = h.decoder.Decode(&comm, r.MultipartForm.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := comm.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := CommTemplateData{
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err = h.cmc.CreateComment(r.Context(), &commgrpc.CreateCommentRequest{
		Comm: &commgrpc.Comment{
			Name:    comm.Name,
			Email:   comm.Email,
			Comment: comm.Comment,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CollTemplateData{
		Message: map[string]string{"SuccessMessage": "Thank you for contacting us. We will reach out to you via email soon."},
	}
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) listComment(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("comm-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	clst, err := h.cmc.ListComment(r.Context(), &commgrpc.ListCommentRequest{
		Filter: &commgrpc.Filter{
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

	commList := make([]Comment, 0, len(clst.GetComm()))
	for _, item := range clst.GetComm() {
		cData := Comment{
			CommentID:    item.CommentID,
			SerialNumber: item.SerialNumber,
			Name:         item.Name,
			Email:        item.Email,
			Comment:      item.Comment,

			CreatedAt: item.CreatedAt.AsTime(),
		}
		commList = append(commList, cData)
	}

	msg := map[string]string{}
	if filterData.SearchTerm != "" && len(clst.GetComm()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(clst.GetComm()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := CommTemplateData{
		FilterData:     *filterData,
		List:           commList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: commentListPath,
	}
	if len(commList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, int32(len(commList)), r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["comment_id"]
	res, err := h.cmc.GetComment(r.Context(), &commgrpc.GetCommentRequest{
		Comm: &commgrpc.Comment{
			CommentID: id,
		},
	})
	if err != nil {
		log.Println("unable to get comment info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CommTemplateData{
		Comm: Comment{
			CommentID:    res.Comm.CommentID,
			SerialNumber: res.Comm.SerialNumber,
			Name:         res.Comm.Name,
			Email:        res.Comm.Email,
			Comment:      res.Comm.Comment,
			CreatedAt:    res.Comm.CreatedAt.AsTime(),
		},
		URLs:           listOfURLs(),
		CurrentPageURL: commentListPath,
	}

	err = h.templates.ExecuteTemplate(w, "comm-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
