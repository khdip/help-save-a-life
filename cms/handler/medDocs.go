package handler

import (
	"fmt"
	"help-save-a-life/cms/paginator"
	docsgrpc "help-save-a-life/proto/medDocs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type MedDocs struct {
	ID           string
	SerialNumber int32
	Name         string
	Type         string
	UploadedAt   time.Time
	UploadedBy   string
	DeletedAt    time.Time
	DeletedBy    string
}

type MedDocsTemplateData struct {
	MedDocs        MedDocs
	List           []MedDocs
	Paginator      paginator.Paginator
	FilterData     Filter
	URLs           map[string]string
	Message        map[string]string
	FormErrors     map[string]string
	CurrentPageURL string
	Title          string
}

func (h *Handler) createMedDocs(w http.ResponseWriter, r *http.Request) {
	data := MedDocsTemplateData{
		MedDocs:        MedDocs{},
		URLs:           listOfURLs(),
		CurrentPageURL: medDocsListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	h.loadMedDocsCreateForm(w, data)
}

func (h *Handler) storeMedDocs(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["Name"]
	data := MedDocsTemplateData{
		URLs:           listOfURLs(),
		CurrentPageURL: medDocsListPath,
		Title:          h.getSettingsTitle(w, r),
	}

	if files == nil {
		data.FormErrors = map[string]string{"Error": "No file selected."}
		h.loadMedDocsCreateForm(w, data)
		return
	}

	for _, fileHeader := range files {
		if fileHeader.Size > 5*1024*1024 {
			data.FormErrors = map[string]string{"Error": "File size too big. Please upload a file less than 5MB."}
			h.loadMedDocsCreateForm(w, data)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "application/pdf" {
			data.FormErrors = map[string]string{"Error": "Invalid file type. Please upload jpg/png/pdf file."}
			h.loadMedDocsCreateForm(w, data)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileName := strings.ReplaceAll(fileHeader.Filename, " ", "")
		if len(fileName) > 90 {
			fileName = fileHeader.Filename[:90] + filepath.Ext(fileName)
		}

		f, err := os.Create(fmt.Sprintf("./cms/assets/files/%s", fileName))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		check, err := h.mds.ListMedDocs(r.Context(), &docsgrpc.ListMedDocsRequest{
			Filter: &docsgrpc.Filter{
				SearchTerm: fileName,
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if check.GetDocs() != nil {
			data.FormErrors = map[string]string{"Error": "Duplicate file name."}
			h.loadMedDocsCreateForm(w, data)
			return
		}

		_, err = h.mds.CreateMedDocs(r.Context(), &docsgrpc.CreateMedDocsRequest{
			Docs: &docsgrpc.MedDocs{
				Name:       fileName,
				Type:       filetype,
				UploadedBy: h.getLoggedUser(r),
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	data.Message = map[string]string{"Message": "File uploaded successfully"}
	h.loadMedDocsCreateForm(w, data)
}

func (h *Handler) listMedDocs(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("docs-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	filterData := GetFilterData(r)
	users := h.getUserListMap(w, r)
	dlst, err := h.mds.ListMedDocs(r.Context(), &docsgrpc.ListMedDocsRequest{
		Filter: &docsgrpc.Filter{
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

	docsList := make([]MedDocs, 0, len(dlst.GetDocs()))
	for _, item := range dlst.GetDocs() {
		dData := MedDocs{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			Name:         item.Name,
			Type:         item.Type,
			UploadedAt:   item.UploadedAt.AsTime(),
			UploadedBy:   users[item.UploadedBy],
		}
		docsList = append(docsList, dData)
	}

	docsstat, err := h.mds.MedDocsStats(r.Context(), &docsgrpc.MedDocsStatsRequest{
		Filter: &docsgrpc.Filter{
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
	if filterData.SearchTerm != "" && len(dlst.GetDocs()) > 0 {
		msg = map[string]string{"FoundMessage": "Data Found"}
	} else if filterData.SearchTerm != "" && len(dlst.GetDocs()) == 0 {
		msg = map[string]string{"NotFoundMessage": "Data Not Found"}
	}
	data := MedDocsTemplateData{
		FilterData:     *filterData,
		List:           docsList,
		Message:        msg,
		URLs:           listOfURLs(),
		CurrentPageURL: medDocsListPath,
		Title:          h.getSettingsTitle(w, r),
	}
	if len(docsList) > 0 {
		data.Paginator = paginator.NewPaginator(int32(filterData.CurrentPage), limitPerPage, docsstat.Stats.Count, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}
}

func (h *Handler) viewMedDocs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res, err := h.mds.GetMedDocs(r.Context(), &docsgrpc.GetMedDocsRequest{
		Docs: &docsgrpc.MedDocs{
			ID: id,
		},
	})
	if err != nil {
		log.Println("unable to get doc info: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := h.getUserListMap(w, r)
	data := MedDocsTemplateData{
		MedDocs: MedDocs{
			ID:           res.Docs.ID,
			SerialNumber: res.Docs.SerialNumber,
			Name:         res.Docs.Name,
			Type:         res.Docs.Type,
			UploadedAt:   res.Docs.UploadedAt.AsTime(),
			UploadedBy:   users[res.Docs.UploadedBy],
		},
		URLs:           listOfURLs(),
		CurrentPageURL: medDocsListPath,
		Title:          h.getSettingsTitle(w, r),
	}

	err = h.templates.ExecuteTemplate(w, "docs-view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteMedDocs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, err := h.mds.DeleteMedDocs(r.Context(), &docsgrpc.DeleteMedDocsRequest{
		Docs: &docsgrpc.MedDocs{
			ID:        id,
			DeletedBy: h.getLoggedUser(r),
		},
	}); err != nil {
		log.Println("unable to delete doc: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, medDocsListPath, http.StatusSeeOther)
}

func (h *Handler) loadMedDocsCreateForm(w http.ResponseWriter, data MedDocsTemplateData) {
	err := h.templates.ExecuteTemplate(w, "docs-create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getMedDocsList(w http.ResponseWriter, r *http.Request) []MedDocs {
	dlst, err := h.mds.ListMedDocs(r.Context(), &docsgrpc.ListMedDocsRequest{
		Filter: &docsgrpc.Filter{},
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, notFoundPath, http.StatusSeeOther)
	}

	docsList := make([]MedDocs, 0, len(dlst.GetDocs()))
	for _, item := range dlst.GetDocs() {
		dData := MedDocs{
			Name: item.Name,
		}
		docsList = append(docsList, dData)
	}
	return docsList
}
