package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	"github.com/jae2274/goutils/llog"
)

const (
	initPage                   = 0
	messageInternalServerError = "Internal Server Error"
)

type JobPostingController struct {
	postingService posting.PostingService
}

func NewJobPostingController(service posting.PostingService) *JobPostingController {
	return &JobPostingController{
		postingService: service,
	}
}

func (c *JobPostingController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/job_postings", c.JobPostings).Methods("GET")
	router.HandleFunc("/job_postings/{site}/{postingId}", c.JobPostingDetail).Methods("GET")
	router.HandleFunc("/categories", c.Categories).Methods("GET")
	router.HandleFunc("/skills", c.Skills).Methods("GET")
	// c.router.HandleFunc(rootPath + "/match_job", c.).Methods("GET")
}

func (c *JobPostingController) JobPostings(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingsRequest(r, initPage)
	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId *string
	claims, ok := middleware.GetClaims(r.Context())
	if ok {
		userId = &claims.UserId
	}

	resp, err := c.postingService.JobPostings(reqCtx, userId, req)
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

func (c *JobPostingController) JobPostingDetail(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingDetailRequest(r)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId *string
	claims, isExisted := middleware.GetClaims(r.Context())
	if isExisted {
		userId = &claims.UserId
	}

	resp, err := c.postingService.JobPostingDetail(reqCtx, userId, req)

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

func (c *JobPostingController) Categories(w http.ResponseWriter, r *http.Request) {
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

func (c *JobPostingController) Skills(w http.ResponseWriter, r *http.Request) {
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
