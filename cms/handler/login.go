package handler

import (
	usergrpc "help-save-a-life/proto/users"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]string
	Notices  string
	URLs     map[string]string
}

func (l LoginFormData) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required.Error("The email field is required")),
		validation.Field(&l.Password, validation.Required.Error("The password field is required"), validation.Length(6, 20).Error("Password must be between 6 to 20 characters")),
	)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	form := LoginFormData{
		URLs: listOfURLs(),
	}

	if flashes := session.Flashes(); len(flashes) > 0 {
		if val, ok := flashes[0].(string); ok {
			form.Notices = val
		}
	}

	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	h.loadLoginForm(w, form)
}

func (h *Handler) loginAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var form LoginFormData
	form.URLs = listOfURLs()
	err = h.decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Fatal(err)
	}

	err = form.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for _, value := range vErrors {
				ErrorValue = map[string]string{
					"Error": value.Error(),
				}
			}
			form.Errors = ErrorValue
			h.loadLoginForm(w, form)
			return
		}
	}

	user, err := h.uc.GetUser(r.Context(), &usergrpc.GetUserRequest{
		User: &usergrpc.User{
			Email: form.Email,
		},
	})
	if err != nil {
		log.Println("unable to get user info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.User.UserID == "" {
		form.Errors = map[string]string{"Error": "Invalid User"}
		h.loadLoginForm(w, form)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.User.Password), []byte(form.Password)); err != nil {
		form.Errors = map[string]string{"Error": "Wrong Password"}
		h.loadLoginForm(w, form)
		return
	}

	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Options.HttpOnly = true

	session.Values["authUserId"] = user.User.UserID
	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, dashboardPath, http.StatusTemporaryRedirect)
}

func (h *Handler) loadLoginForm(w http.ResponseWriter, form LoginFormData) {
	err := h.templates.ExecuteTemplate(w, "login.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
