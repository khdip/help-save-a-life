package handler

const (
	homePath              = "/"
	notFoundPath          = "/404"
	unauthorizedPath      = "/unauthorized"
	dashboardPath         = "/dashboard"
	loginPath             = "/login"
	loginAuthPath         = "/login/auth"
	logoutPath            = "/logout"
	userListPath          = "/users"
	userCreatePath        = "/users/create"
	userStorePath         = "/users/store"
	userEditPath          = "/users/edit/{user_id}"
	userUpdatePath        = "/users/update/{user_id}"
	userDeletePath        = "/users/delete/{user_id}"
	userViewPath          = "/users/view/{user_id}"
	collectionListPath    = "/collection"
	collectionCreatePath  = "/collection/create"
	collectionStorePath   = "/collection/store"
	collectionEditPath    = "/collection/edit/{collection_id}"
	collectionUpdatePath  = "/collection/update/{collection_id}"
	collectionDeletePath  = "/collection/delete/{collection_id}"
	collectionViewPath    = "/collection/view/{collection_id}"
	commentListPath       = "/comments"
	commentCreatePath     = "/comments/create"
	commentStorePath      = "/comments/store"
	commentViewPath       = "/comments/view/{comment_id}"
	dailyReportListPath   = "/daily_report"
	dailyReportCreatePath = "/daily_report/create"
	dailyReportStorePath  = "/daily_report/store"
	dailyReportEditPath   = "/daily_report/edit/{report_id}"
	dailyReportUpdatePath = "/daily_report/update/{report_id}"
	dailyReportDeletePath = "/daily_report/delete/{report_id}"
	dailyReportViewPath   = "/daily_report/view/{report_id}"
	currencyListPath      = "/currencies"
	currencyCreatePath    = "/currency/create"
	currencyStorePath     = "/currency/store"
	currencyEditPath      = "/currency/edit/{id}"
	currencyUpdatePath    = "/currency/update/{id}"
	currencyDeletePath    = "/currency/delete/{id}"
	currencyViewPath      = "/currency/view/{id}"
	accountTypeListPath   = "/account_types"
	accountTypeCreatePath = "/account_type/create"
	accountTypeStorePath  = "/account_type/store"
	accountTypeEditPath   = "/account_type/edit/{id}"
	accountTypeUpdatePath = "/account_type/update/{id}"
	accountTypeDeletePath = "/account_type/delete/{id}"
	accountTypeViewPath   = "/account_type/view/{id}"
	accountsListPath      = "/accounts"
	accountsCreatePath    = "/account/create"
	accountsStorePath     = "/account/store"
	accountsEditPath      = "/account/edit/{id}"
	accountsUpdatePath    = "/account/update/{id}"
	accountsDeletePath    = "/account/delete/{id}"
	accountsViewPath      = "/account/view/{id}"
	settingsEditPath      = "/settings"
	settingsUpdatePath    = "/settings/update"
)

func listOfURLs() map[string]string {
	return map[string]string{
		"home":         homePath,
		"unauthorized": unauthorizedPath,
		"login":        loginPath,
		"logout":       logoutPath,
		"dashboard":    dashboardPath,
		"userList":     userListPath,
		"userCreate":   userCreatePath,
		"userStore":    userStorePath,
		"userEdit":     userEditPath,
		"userUpdate":   userUpdatePath,
		"userDelete":   userDeletePath,
		"userView":     userViewPath,
		"collList":     collectionListPath,
		"collCreate":   collectionCreatePath,
		"collStore":    collectionStorePath,
		"collEdit":     collectionEditPath,
		"collUpdate":   collectionUpdatePath,
		"collDelete":   collectionDeletePath,
		"collView":     collectionViewPath,
		"commCreate":   commentCreatePath,
		"commStore":    commentStorePath,
		"commList":     commentListPath,
		"commView":     commentViewPath,
		"dreList":      dailyReportListPath,
		"dreCreate":    dailyReportCreatePath,
		"dreStore":     dailyReportStorePath,
		"dreEdit":      dailyReportEditPath,
		"dreUpdate":    dailyReportUpdatePath,
		"dreDelete":    dailyReportDeletePath,
		"dreView":      currencyViewPath,
		"currList":     currencyListPath,
		"currCreate":   currencyCreatePath,
		"currStore":    currencyStorePath,
		"currEdit":     currencyEditPath,
		"currUpdate":   currencyUpdatePath,
		"currDelete":   currencyDeletePath,
		"currView":     currencyViewPath,
		"acctList":     accountTypeListPath,
		"acctCreate":   accountTypeCreatePath,
		"acctStore":    accountTypeStorePath,
		"acctEdit":     accountTypeEditPath,
		"acctUpdate":   accountTypeUpdatePath,
		"acctDelete":   accountTypeDeletePath,
		"acctView":     accountTypeViewPath,
		"acntList":     accountsListPath,
		"acntCreate":   accountsCreatePath,
		"acntStore":    accountsStorePath,
		"acntEdit":     accountsEditPath,
		"acntUpdate":   accountsUpdatePath,
		"acntDelete":   accountsDeletePath,
		"acntView":     currencyViewPath,
		"settEdit":     settingsEditPath,
		"settUpdate":   settingsUpdatePath,
	}
}
