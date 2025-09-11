package handler

import (
	"help-save-a-life/cms/paginator"
	commgrpc "help-save-a-life/proto/comments"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Comment struct {
	CommentID int32
	Name      string
	Email     string
	Comment   string
	CreatedAt time.Time
}

type CommTemplateData struct {
	Comm       Comment
	List       []Comment
	Paginator  paginator.Paginator
	FilterData Filter
	URLs       map[string]string
	Message    map[string]string
}

func (h *Handler) storeComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var comm Comment
	err = h.decoder.Decode(&comm, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	http.Redirect(w, r, homePath, http.StatusTemporaryRedirect)
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
			CommentID: item.CommentID,
			Name:      item.Name,
			Email:     item.Email,
			Comment:   item.Comment,

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
		FilterData: *filterData,
		List:       commList,
		Message:    msg,
		URLs:       listOfURLs(),
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
	cid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.cmc.GetComment(r.Context(), &commgrpc.GetCommentRequest{
		Comm: &commgrpc.Comment{
			CommentID: int32(cid),
		},
	})
	if err != nil {
		log.Println("unable to get comment info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CommTemplateData{
		Comm: Comment{
			CommentID: res.Comm.CommentID,
			Name:      res.Comm.Name,
			Email:     res.Comm.Email,
			Comment:   res.Comm.Comment,
			CreatedAt: res.Comm.CreatedAt.AsTime(),
		},
		URLs: listOfURLs(),
	}

	err = h.templates.ExecuteTemplate(w, "comm-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
