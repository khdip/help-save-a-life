package handler

const (
	homePath              = "/"
	notFoundPath          = "/404"
	dashboardPath         = "/dashboard"
	loginPath             = "/login"
	loginAuthPath         = "/login/auth"
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
)

func listOfURLs() map[string]string {
	return map[string]string{
		"home":       homePath,
		"dashboard":  dashboardPath,
		"userList":   userListPath,
		"userCreate": userCreatePath,
		"userStore":  userStorePath,
		"userEdit":   userEditPath,
		"userUpdate": userUpdatePath,
		"userDelete": userDeletePath,
		"userView":   userViewPath,
		"collList":   collectionListPath,
		"collCreate": collectionCreatePath,
		"collStore":  collectionStorePath,
		"collEdit":   collectionEditPath,
		"collUpdate": collectionUpdatePath,
		"collDelete": collectionDeletePath,
		"collView":   collectionViewPath,
		"commCreate": commentCreatePath,
		"commStore":  commentStorePath,
		"commList":   commentListPath,
		"commView":   commentViewPath,
		"dreList":    dailyReportListPath,
		"dreCreate":  dailyReportCreatePath,
		"dreStore":   dailyReportStorePath,
		"dreEdit":    dailyReportEditPath,
		"dreUpdate":  dailyReportUpdatePath,
		"dreDelete":  dailyReportDeletePath,
		"dreView":    dailyReportViewPath,
	}
}
