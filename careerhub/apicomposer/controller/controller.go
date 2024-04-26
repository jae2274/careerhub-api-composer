package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	"github.com/jae2274/goutils/llog"
)

const (
	initPage                   = 0
	messageInternalServerError = "Internal Server Error"
)

type Controller struct {
	postingService posting.PostingService
}

func NewController(service posting.PostingService, router *mux.Router) *Controller {
	return &Controller{
		postingService: service,
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/job_postings", c.JobPostings).Methods("GET")
	router.HandleFunc("/job_postings/{site}/{postingId}", c.JobPostingDetail).Methods("GET")
	router.HandleFunc("/categories", c.Categories).Methods("GET")
	router.HandleFunc("/skills", c.Skills).Methods("GET")
	// c.router.HandleFunc(rootPath + "/match_job", c.).Methods("GET")
}

func (c *Controller) JobPostings(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingsRequest(r, initPage)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.postingService.JobPostings(reqCtx, req)
	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// jobPostings를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) JobPostingDetail(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingDetailRequest(r)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.postingService.JobPostingDetail(reqCtx, req)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// jobPostingDetail을 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) Categories(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	categories, err := c.postingService.Categories(reqCtx)
	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// categories를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (c *Controller) Skills(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	skills, err := c.postingService.Skills(reqCtx)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// skills를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(skills)
}
