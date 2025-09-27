package handler

import (
	"help-save-a-life/cms/paginator"
	"log"
	"net/http"
	"time"

	usergrpc "help-save-a-life/proto/users"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    string
	Name      string
	Batch     int32
	Email     string
	Password  string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt time.Time
	DeletedBy string
}

type UserTemplateData struct {
	User           User
	List           []User
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	CurrentPageURL string
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	h.loadUserCreateForm(w, User{})
}

func (h *Handler) storeUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var usr User
	err = h.decoder.Decode(&usr, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	passByte, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = h.uc.CreateUser(r.Context(), &usergrpc.CreateUserRequest{
		User: &usergrpc.User{
			Name:      usr.Name,
			Batch:     usr.Batch,
			Email:     usr.Email,
			Password:  string(passByte),
			CreatedBy: h.getLoggedUser(r),
			UpdatedBy: h.getLoggedUser(r),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, userListPath, http.StatusTemporaryRedirect)
}

func (h *Handler) editUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]
	res, err := h.uc.GetUser(r.Context(), &usergrpc.GetUserRequest{
		User: &usergrpc.User{
			UserID: id,
		},
	})
	if err != nil {
		log.Println("unable to get user info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadUserEditForm(w, User{
		UserID: res.User.UserID,
		Name:   res.User.Name,
		Batch:  res.User.Batch,
		Email:  res.User.Email,
	})
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["user_id"]
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var usr User
	if err := h.decoder.Decode(&usr, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	passByte, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := h.uc.UpdateUser(ctx, &usergrpc.UpdateUserRequest{
		User: &usergrpc.User{
			UserID:    id,
			Name:      usr.Name,
			Batch:     usr.Batch,
			Email:     usr.Email,
			Password:  string(passByte),
			UpdatedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, userListPath, http.StatusSeeOther)
}

func (h *Handler) listUser(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("user-list.html")
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
	usrlst, err := h.uc.ListUser(r.Context(), &usergrpc.ListUserRequest{
		Filter: &usergrpc.Filter{
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

	userList := make([]User, 0, len(usrlst.GetUser()))
	for _, item := range usrlst.GetUser() {
		usrData := User{
			UserID:    item.UserID,
			Name:      item.Name,
			Batch:     item.Batch,
			Email:     item.Email,
			CreatedAt: item.CreatedAt.AsTime(),
			CreatedBy: item.CreatedBy,
			UpdatedAt: item.UpdatedAt.AsTime(),
			UpdatedBy: item.UpdatedBy,
		}
		userList = append(userList, usrData)
	}

	userstat, err := h.uc.UserStats(r.Context(), &usergrpc.UserStatsRequest{
		Filter: &usergrpc.Filter{
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
	if filterData.SearchTerm != "" && len(usrlst.GetUser()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(usrlst.GetUser()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := UserTemplateData{
		FilterData:     *filterData,
		List:           userList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: userListPath,
	}
	if len(userList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, userstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]
	res, err := h.uc.GetUser(r.Context(), &usergrpc.GetUserRequest{
		User: &usergrpc.User{
			UserID: id,
		},
	})
	if err != nil {
		log.Println("unable to get user info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := UserTemplateData{
		User: User{
			UserID:    id,
			Name:      res.User.Name,
			Batch:     res.User.Batch,
			Email:     res.User.Email,
			CreatedAt: res.User.CreatedAt.AsTime(),
			CreatedBy: res.User.CreatedBy,
			UpdatedAt: res.User.UpdatedAt.AsTime(),
			UpdatedBy: res.User.UpdatedBy,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: userListPath,
	}

	err = h.templates.ExecuteTemplate(w, "user-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]
	if _, err := h.uc.DeleteUser(r.Context(), &usergrpc.DeleteUserRequest{
		User: &usergrpc.User{
			UserID:    id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete user: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, userListPath, http.StatusSeeOther)
}

func (h *Handler) loadUserCreateForm(w http.ResponseWriter, usr User) {
	form := UserTemplateData{
		User:           usr,
		URLs:           listOfURLs(),
		CurrentPageURL: userListPath,
	}

	err := h.templates.ExecuteTemplate(w, "user-create.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadUserEditForm(w http.ResponseWriter, usr User) {
	form := UserTemplateData{
		User:           usr,
		URLs:           listOfURLs(),
		CurrentPageURL: userListPath,
	}

	err := h.templates.ExecuteTemplate(w, "user-edit.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getName(w http.ResponseWriter, r *http.Request, userid string) string {
	res, err := h.uc.GetUser(r.Context(), &usergrpc.GetUserRequest{
		User: &usergrpc.User{
			UserID: userid,
		},
	})
	if err != nil {
		log.Println("unable to get user info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return res.User.Name
}
