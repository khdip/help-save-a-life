package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	acctgrpc "help-save-a-life/proto/accountType"
	acntgrpc "help-save-a-life/proto/accounts"
	collgrpc "help-save-a-life/proto/collection"
	commgrpc "help-save-a-life/proto/comments"
	currgrpc "help-save-a-life/proto/currency"
	dregrpc "help-save-a-life/proto/dailyReport"
	settgrpc "help-save-a-life/proto/settings"
	usergrpc "help-save-a-life/proto/users"
)

type Handler struct {
	templates *template.Template
	decoder   *schema.Decoder
	session   *sessions.CookieStore
	assets    fs.FS
	assetFS   *hashfs.FS
	uc        usergrpc.UserServiceClient
	cc        collgrpc.CollectionServiceClient
	cmc       commgrpc.CommentServiceClient
	drc       dregrpc.DailyReportServiceClient
	cuc       currgrpc.CurrencyServiceClient
	sc        settgrpc.SettingsServiceClient
	at        acctgrpc.AccountTypeServiceClient
	acc       acntgrpc.AccountsServiceClient
}

func GetHandler(decoder *schema.Decoder, session *sessions.CookieStore, assets fs.FS, uc usergrpc.UserServiceClient, cc collgrpc.CollectionServiceClient, cmc commgrpc.CommentServiceClient, drc dregrpc.DailyReportServiceClient, cuc currgrpc.CurrencyServiceClient, sc settgrpc.SettingsServiceClient, at acctgrpc.AccountTypeServiceClient, acc acntgrpc.AccountsServiceClient) *mux.Router {
	hand := &Handler{
		decoder: decoder,
		session: session,
		assets:  assets,
		assetFS: hashfs.NewFS(assets),
		uc:      uc,
		cc:      cc,
		cmc:     cmc,
		drc:     drc,
		cuc:     cuc,
		sc:      sc,
		at:      at,
		acc:     acc,
	}
	hand.GetTemplate()

	r := mux.NewRouter()
	r.HandleFunc(homePath, hand.homeHandler)
	r.HandleFunc(commentStorePath, hand.storeComment)
	r.HandleFunc(unauthorizedPath, hand.unauthorizedHandler)
	r.HandleFunc(logoutPath, hand.logout)

	loginRouter := r.NewRoute().Subrouter()
	loginRouter.HandleFunc(loginPath, hand.login)
	loginRouter.HandleFunc(loginAuthPath, hand.loginAuth)
	loginRouter.Use(hand.restrictMiddleware)

	s := r.NewRoute().Subrouter()
	s.HandleFunc(dashboardPath, hand.viewDashboard)
	s.HandleFunc(userCreatePath, hand.createUser)
	s.HandleFunc(userStorePath, hand.storeUser)
	s.HandleFunc(userEditPath, hand.editUser)
	s.HandleFunc(userUpdatePath, hand.updateUser)
	s.HandleFunc(userListPath, hand.listUser)
	s.HandleFunc(userViewPath, hand.viewUser)
	s.HandleFunc(userDeletePath, hand.deleteUser)
	s.HandleFunc(collectionCreatePath, hand.createCollection)
	s.HandleFunc(collectionStorePath, hand.storeCollection)
	s.HandleFunc(collectionEditPath, hand.editCollection)
	s.HandleFunc(collectionUpdatePath, hand.updateCollection)
	s.HandleFunc(collectionListPath, hand.listCollection)
	s.HandleFunc(collectionViewPath, hand.viewCollection)
	s.HandleFunc(collectionDeletePath, hand.deleteCollection)
	s.HandleFunc(commentListPath, hand.listComment)
	s.HandleFunc(commentViewPath, hand.viewComment)
	s.HandleFunc(dailyReportCreatePath, hand.createDailyReport)
	s.HandleFunc(dailyReportStorePath, hand.storeDailyReport)
	s.HandleFunc(dailyReportEditPath, hand.editDailyReport)
	s.HandleFunc(dailyReportUpdatePath, hand.updateDailyReport)
	s.HandleFunc(dailyReportListPath, hand.listDailyReport)
	s.HandleFunc(dailyReportViewPath, hand.viewDailyReport)
	s.HandleFunc(dailyReportDeletePath, hand.deleteDailyReport)
	s.HandleFunc(currencyCreatePath, hand.createCurrency)
	s.HandleFunc(currencyStorePath, hand.storeCurrency)
	s.HandleFunc(currencyEditPath, hand.editCurrency)
	s.HandleFunc(currencyUpdatePath, hand.updateCurrency)
	s.HandleFunc(currencyListPath, hand.listCurrency)
	s.HandleFunc(currencyViewPath, hand.viewCurrency)
	s.HandleFunc(currencyDeletePath, hand.deleteCurrency)
	s.HandleFunc(accountTypeCreatePath, hand.createAccountType)
	s.HandleFunc(accountTypeStorePath, hand.storeAccountType)
	s.HandleFunc(accountTypeEditPath, hand.editAccountType)
	s.HandleFunc(accountTypeUpdatePath, hand.updateAccountType)
	s.HandleFunc(accountTypeListPath, hand.listAccountType)
	s.HandleFunc(accountTypeViewPath, hand.viewAccountType)
	s.HandleFunc(accountTypeDeletePath, hand.deleteAccountType)
	s.HandleFunc(accountsCreatePath, hand.createAccounts)
	s.HandleFunc(accountsStorePath, hand.storeAccounts)
	s.HandleFunc(accountsEditPath, hand.editAccounts)
	s.HandleFunc(accountsUpdatePath, hand.updateAccounts)
	s.HandleFunc(accountsListPath, hand.listAccounts)
	s.HandleFunc(accountsViewPath, hand.viewAccounts)
	s.HandleFunc(accountsDeletePath, hand.deleteAccounts)
	s.HandleFunc(settingsEditPath, hand.editSettings)
	s.HandleFunc(settingsUpdatePath, hand.saveSettings)
	s.Use(hand.authMiddleware)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.FS(hand.assetFS))))

	type NotFoundTempData struct {
		URLs map[string]string
		Sett Settings
	}
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hand.templates.ExecuteTemplate(w, "404.html", NotFoundTempData{
			URLs: listOfURLs(),
			Sett: hand.getSettings(w, r),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) GetTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"cms/assets/templates/layout/header.html",
		"cms/assets/templates/layout/banner.html",
		"cms/assets/templates/layout/admin-header.html",
		"cms/assets/templates/base/index.html",
		"cms/assets/templates/base/dashboard.html",
		"cms/assets/templates/users/user-list.html",
		"cms/assets/templates/users/user-create.html",
		"cms/assets/templates/users/user-edit.html",
		"cms/assets/templates/users/user-view.html",
		"cms/assets/templates/collection/coll-list.html",
		"cms/assets/templates/collection/coll-create.html",
		"cms/assets/templates/collection/coll-edit.html",
		"cms/assets/templates/collection/coll-view.html",
		"cms/assets/templates/dailyReport/dre-list.html",
		"cms/assets/templates/dailyReport/dre-create.html",
		"cms/assets/templates/dailyReport/dre-edit.html",
		"cms/assets/templates/dailyReport/dre-view.html",
		"cms/assets/templates/comments/comm-list.html",
		"cms/assets/templates/comments/comm-view.html",
		"cms/assets/templates/currency/curr-list.html",
		"cms/assets/templates/currency/curr-create.html",
		"cms/assets/templates/currency/curr-edit.html",
		"cms/assets/templates/currency/curr-view.html",
		"cms/assets/templates/accountType/acct-list.html",
		"cms/assets/templates/accountType/acct-create.html",
		"cms/assets/templates/accountType/acct-edit.html",
		"cms/assets/templates/accountType/acct-view.html",
		"cms/assets/templates/accounts/acnt-list.html",
		"cms/assets/templates/accounts/acnt-create.html",
		"cms/assets/templates/accounts/acnt-edit.html",
		"cms/assets/templates/accounts/acnt-view.html",
		"cms/assets/templates/base/404.html",
		"cms/assets/templates/base/unauthorized.html",
		"cms/assets/templates/base/login.html",
		"cms/assets/templates/settings/settings.html",
	))
}
