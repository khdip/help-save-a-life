package handler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/text/message"
)

const (
	limitPerPage = 10
	targetamount = 7500000
	sessionName  = "tamal-session"
)

type Filter struct {
	SearchTerm  string
	PageNumber  int32
	CurrentPage int32
	Offset      int32
	Limit       int32
	SortBy      string
	Order       string
}

func hideDigits(s string) string {
	modifiedStr := s[:len(s)-4]
	return modifiedStr + "****"
}

func formatWithCommas(number int32) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%d", number)
}

func GetFilterData(r *http.Request) *Filter {
	var data Filter
	queryParams := r.URL.Query()
	var err error
	data.SearchTerm, err = url.PathUnescape(queryParams.Get("SearchTerm"))
	if err != nil {
		data.SearchTerm = ""
	}
	data.SortBy = "created_at"
	data.SortBy, err = url.PathUnescape(queryParams.Get("SortBy"))
	if err != nil {
		data.SortBy = "created_at"
	}
	data.Order = "ASC"
	data.Order, err = url.PathUnescape(queryParams.Get("Order"))
	if err != nil {
		data.Order = "ASC"
	}
	page, err := url.PathUnescape(queryParams.Get("page"))
	if err != nil {
		page = "1"
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		pageNumber = 1
	}
	data.PageNumber = int32(pageNumber)
	var offset int32 = 0
	currentPage := pageNumber
	if currentPage <= 0 {
		currentPage = 1
	} else {
		offset = limitPerPage*int32(currentPage) - limitPerPage
	}
	data.CurrentPage = int32(currentPage)
	data.Offset = offset
	return &data
}

func (h *Handler) getLoggedUser(w http.ResponseWriter, r *http.Request) string {
	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	authUserID := session.Values["authUserId"]
	if authUserID != nil {
		return authUserID.(string)
	} else {
		return ""
	}
}
