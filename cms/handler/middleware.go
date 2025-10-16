package handler

import (
	"log"
	"net/http"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.session.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}

		authUserID := session.Values["authUserId"]

		if authUserID != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, unauthorizedPath, http.StatusTemporaryRedirect)
		}
	})
}

func (h *Handler) restrictMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.session.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}

		authUserID := session.Values["authUserId"]

		if authUserID != nil {
			http.Redirect(w, r, homePath, http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
