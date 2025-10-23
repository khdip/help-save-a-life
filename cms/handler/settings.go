package handler

import (
	"log"
	"net/http"

	settgrpc "help-save-a-life/proto/settings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Settings struct {
	PatientName                  string
	Title                        string
	BannerTitle                  string
	HighlightedBannerTitle       string
	BannerDescription            string
	HighlightedBannerDescription string
	BannerImage                  string
	AboutPatient                 string
	TargetAmount                 int32
	ShowMedicalDocuments         bool
	ShowCollection               bool
	ShowDailyReport              bool
	ShowFundUpdates              bool
	CalculateCollection          int32
	TotalAmount                  int32
	UpdatedBy                    string
}

type SettingsData struct {
	Sett           Settings
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
}

func (s Settings) Validate(h *Handler) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.PatientName,
			validation.Required.Error("Patient name field can not be empty"),
			validation.Length(2, 20).Error("The length must be between 2 to 20 characters"),
		),
		validation.Field(&s.Title,
			validation.Required.Error("Title field can not be empty"),
			validation.Length(2, 30).Error("The length must be between 2 to 30 characters"),
		),
		validation.Field(&s.BannerTitle,
			validation.Required.Error("Banner title field can not be empty"),
			validation.Length(3, 100).Error("The length must be between 3 to 100 characters"),
		),
		validation.Field(&s.HighlightedBannerTitle,
			validation.Length(0, 100).Error("The length must be less than 100 characters"),
		),
		validation.Field(&s.BannerDescription,
			validation.Required.Error("Banner description field can not be empty"),
			validation.Length(10, 1200).Error("The length must be between 10 to 500 characters"),
		),
		validation.Field(&s.HighlightedBannerDescription,
			validation.Length(0, 500).Error("The length must be less than 500 characters"),
		),
		validation.Field(&s.AboutPatient,
			validation.Required.Error("About patient field can not be empty"),
		),
		validation.Field(&s.TargetAmount,
			validation.Required.Error("Target amount field can not be empty"),
		),
	)
}

func (h *Handler) editSettings(w http.ResponseWriter, r *http.Request) {
	res, err := h.sc.GetSettings(r.Context(), &settgrpc.GetSettingsRequest{})
	if err != nil {
		log.Println("unable to get settings info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.loadSettingsForm(w, SettingsData{
		Sett: Settings{
			PatientName:                  res.Sett.PatientName,
			Title:                        res.Sett.Title,
			BannerTitle:                  res.Sett.BannerTitle,
			HighlightedBannerTitle:       res.Sett.HighlightedBannerTitle,
			BannerDescription:            res.Sett.BannerDescription,
			HighlightedBannerDescription: res.Sett.HighlightedBannerDescription,
			BannerImage:                  res.Sett.BannerImage,
			AboutPatient:                 res.Sett.AboutPatient,
			TargetAmount:                 res.Sett.TargetAmount,
			ShowMedicalDocuments:         res.Sett.ShowMedicalDocuments,
			ShowCollection:               res.Sett.ShowCollection,
			ShowDailyReport:              res.Sett.ShowDailyReport,
			ShowFundUpdates:              res.Sett.ShowFundUpdates,
			CalculateCollection:          res.Sett.CalculateCollection,
		},
		URLs:           listOfURLs(),
		CurrentPageURL: settingsEditPath,
	})
}

func (h *Handler) saveSettings(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var sett Settings
	if err := h.decoder.Decode(&sett, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sett.Validate(h); err != nil {
		vErrs := map[string]string{}
		if err, ok := (err).(validation.Errors); ok {
			if len(err) > 0 {
				for k, v := range err {
					vErrs[k] = v.Error()
				}
			}
		}
		data := SettingsData{
			Sett:           sett,
			FormErrors:     vErrs,
			URLs:           listOfURLs(),
			CurrentPageURL: settingsEditPath,
		}
		h.loadSettingsForm(w, data)
		return
	}

	file, fileHeader, _ := r.FormFile("BannerImage")
	image, err := h.saveImage(file, fileHeader, "./cms/assets/images/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.sc.UpdateSettings(r.Context(), &settgrpc.UpdateSettingsRequest{
		Sett: &settgrpc.Settings{
			PatientName:                  sett.PatientName,
			Title:                        sett.Title,
			BannerTitle:                  sett.BannerTitle,
			HighlightedBannerTitle:       sett.HighlightedBannerDescription,
			BannerDescription:            sett.BannerDescription,
			HighlightedBannerDescription: sett.HighlightedBannerDescription,
			BannerImage:                  image,
			AboutPatient:                 sett.AboutPatient,
			TargetAmount:                 sett.TargetAmount,
			ShowMedicalDocuments:         sett.ShowMedicalDocuments,
			ShowCollection:               sett.ShowCollection,
			ShowDailyReport:              sett.ShowDailyReport,
			ShowFundUpdates:              sett.ShowFundUpdates,
			CalculateCollection:          sett.CalculateCollection,
			UpdatedBy:                    h.getLoggedUser(r),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, settingsEditPath, http.StatusTemporaryRedirect)
}

func (h *Handler) loadSettingsForm(w http.ResponseWriter, data SettingsData) {
	err := h.templates.ExecuteTemplate(w, "settings.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
