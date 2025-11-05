package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/message"

	collgrpc "help-save-a-life/proto/collection"
	dregrpc "help-save-a-life/proto/dailyReport"
)

const (
	limitPerPage = 10
	sessionName  = "save-life"
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
	modifiedStr := s
	if len(s) > 5 {
		modifiedStr = s[:len(s)-4] + "****"
	}
	return modifiedStr
}

func formatWithCommas(number float32) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", number)
}

func GetFilterData(r *http.Request) *Filter {
	var data Filter
	queryParams := r.URL.Query()
	var err error
	data.SearchTerm, err = url.PathUnescape(queryParams.Get("SearchTerm"))
	if err != nil {
		data.SearchTerm = ""
	}
	data.SortBy = "serial_number"
	data.SortBy, err = url.PathUnescape(queryParams.Get("SortBy"))
	if err != nil {
		data.SortBy = "serial_number"
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

func (h *Handler) getLoggedUser(r *http.Request) string {
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

func (h *Handler) saveImage(file multipart.File, fileHeader *multipart.FileHeader, imagePath string) (string, error) {
	if fileHeader != nil {
		fileName := strings.ReplaceAll(fileHeader.Filename, " ", "")
		if len(fileName) > 90 {
			fileName = fileHeader.Filename[:90] + ".jpg"
		}
		dest, err := os.Create(fmt.Sprintf(imagePath+"%s", fileName))
		if err != nil {
			return "", err
		}
		defer dest.Close()
		if _, err := io.Copy(dest, file); err != nil {
			fmt.Println(err.Error())
		}
		return fileName, err
	}
	return "", nil
}

func highlightSubstring(text, keyword string, padding int32) template.HTML {
	if keyword == "" || text == "" {
		return template.HTML(text)
	}
	highlighted := fmt.Sprintf(`<span class="bg-success text-light p-%d rounded-3">%s</span>`, padding, keyword)
	return template.HTML(strings.ReplaceAll(text, keyword, highlighted))
}

func makeHTMLTemplate(s string) template.HTML {
	return template.HTML(s)
}

func (h *Handler) getTotalAndTargetAmount(w http.ResponseWriter, r *http.Request, sett SettingsHome) (float32, int32) {
	ctx := r.Context()

	collstat, err := h.cc.CollectionStats(ctx, &collgrpc.CollectionStatsRequest{
		Filter: &collgrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	drstat, err := h.drc.DailyReportStats(ctx, &dregrpc.DailyReportStatsRequest{
		Filter: &dregrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get stats: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	var totalCollection float32
	switch sett.CalculateCollection {
	case 0:
		totalCollection = collstat.Stats.TotalAmount
	case 1:
		totalCollection = drstat.Stats.TotalAmount
	case 2:
		totalCollection = sett.TotalAmount
	}

	targetAmount := sett.TargetAmount
	if totalCollection > float32(targetAmount) {
		totalCollection = float32(targetAmount)
	}

	return totalCollection, targetAmount
}

var theme_1 = []string{
	"cms/assets/templates/theme_1/layout/header.html",
	"cms/assets/templates/theme_1/layout/banner.html",
	"cms/assets/templates/theme_1/layout/about.html",
	"cms/assets/templates/theme_1/layout/med_docs.html",
	"cms/assets/templates/theme_1/layout/collection.html",
	"cms/assets/templates/theme_1/layout/dailyReport.html",
	"cms/assets/templates/theme_1/layout/importantLinks.html",
	"cms/assets/templates/theme_1/layout/footer.html",
	"cms/assets/templates/theme_1/layout/admin-header.html",
	"cms/assets/templates/theme_1/base/index.html",
	"cms/assets/templates/theme_1/base/dashboard.html",
	"cms/assets/templates/theme_1/users/user-list.html",
	"cms/assets/templates/theme_1/users/user-create.html",
	"cms/assets/templates/theme_1/users/user-edit.html",
	"cms/assets/templates/theme_1/users/user-view.html",
	"cms/assets/templates/theme_1/collection/coll-list.html",
	"cms/assets/templates/theme_1/collection/coll-create.html",
	"cms/assets/templates/theme_1/collection/coll-edit.html",
	"cms/assets/templates/theme_1/collection/coll-view.html",
	"cms/assets/templates/theme_1/dailyReport/dre-list.html",
	"cms/assets/templates/theme_1/dailyReport/dre-create.html",
	"cms/assets/templates/theme_1/dailyReport/dre-edit.html",
	"cms/assets/templates/theme_1/dailyReport/dre-view.html",
	"cms/assets/templates/theme_1/comments/comm-list.html",
	"cms/assets/templates/theme_1/comments/comm-view.html",
	"cms/assets/templates/theme_1/currency/curr-list.html",
	"cms/assets/templates/theme_1/currency/curr-create.html",
	"cms/assets/templates/theme_1/currency/curr-edit.html",
	"cms/assets/templates/theme_1/currency/curr-view.html",
	"cms/assets/templates/theme_1/accountType/acct-list.html",
	"cms/assets/templates/theme_1/accountType/acct-create.html",
	"cms/assets/templates/theme_1/accountType/acct-edit.html",
	"cms/assets/templates/theme_1/accountType/acct-view.html",
	"cms/assets/templates/theme_1/accounts/acnt-list.html",
	"cms/assets/templates/theme_1/accounts/acnt-create.html",
	"cms/assets/templates/theme_1/accounts/acnt-edit.html",
	"cms/assets/templates/theme_1/accounts/acnt-view.html",
	"cms/assets/templates/theme_1/links/link-list.html",
	"cms/assets/templates/theme_1/links/link-create.html",
	"cms/assets/templates/theme_1/links/link-edit.html",
	"cms/assets/templates/theme_1/links/link-view.html",
	"cms/assets/templates/theme_1/base/404.html",
	"cms/assets/templates/theme_1/base/unauthorized.html",
	"cms/assets/templates/theme_1/base/login.html",
	"cms/assets/templates/theme_1/settings/settings.html",
}
