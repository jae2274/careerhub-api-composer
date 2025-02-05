package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/domain"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/httputils"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	"github.com/jae2274/goutils/llog"
)

const (
	initPage                   = 1
	messageInternalServerError = "Internal Server Error"
)

type JobPostingController struct {
	postingService *posting.PostingService
}

func NewJobPostingController(postingService *posting.PostingService) *JobPostingController {
	return &JobPostingController{
		postingService: postingService,
	}
}

func (c *JobPostingController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/job_postings", c.JobPostings).Methods("GET")
	router.HandleFunc("/job_postings/count", c.CountJobPostings).Methods("GET")
	router.HandleFunc("/job_postings/{site}/{postingId}", c.JobPostingDetail).Methods("GET")
	router.HandleFunc("/categories", c.Categories).Methods("GET")
	router.HandleFunc("/skills", c.Skills).Methods("GET")
	router.HandleFunc("/companies", c.Companies).Methods("GET")
	// c.router.HandleFunc(rootPath + "/match_job", c.).Methods("GET")
}

func IsInvalidQuery(ctx context.Context, w http.ResponseWriter, queryReq *query.Query, err error) bool {
	if err != nil {
		llog.LogErr(ctx, err) // 파싱 에러 발생으로 별도의 로깅 필요
		http.Error(w, "bad request", http.StatusBadRequest)
		return true
	} else if err := posting.IsValid(queryReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return true
	}

	return false
}

func (c *JobPostingController) JobPostings(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingsRequest(r, initPage)

	if IsInvalidQuery(reqCtx, w, req.QueryReq, err) {
		return
	}

	var jobPostings []*domain.JobPosting

	claims, isExisted := middleware.GetClaims(reqCtx)
	if isExisted {
		jobPostings, err = c.postingService.JobPostingsWithClaims(reqCtx, req, claims)
	} else {
		jobPostings, err = c.postingService.JobPostings(reqCtx, req)
	}

	if httputils.IsError(reqCtx, w, err) {
		return
	}

	// jobPostings를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(jobPostings)
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}

func (c *JobPostingController) CountJobPostings(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	queryReq, err := posting.ExtractJobPostingQuery(r)

	if IsInvalidQuery(reqCtx, w, queryReq, err) {
		return
	}

	count, err := c.postingService.CountJobPostings(reqCtx, queryReq)

	if httputils.IsError(reqCtx, w, err) {
		return
	}

	// count를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&struct {
		Count int64 `json:"count"`
	}{Count: count})
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}

func (c *JobPostingController) JobPostingDetail(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingDetailRequest(r)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, isExisted := middleware.GetClaims(r.Context())

	var resp *domain.JobPostingDetail
	var isExist bool
	if isExisted {
		resp, isExist, err = c.postingService.JobPostingDetailWithClaims(reqCtx, req, claims)
	} else {
		resp, isExist, err = c.postingService.JobPostingDetail(reqCtx, req)
	}

	if httputils.IsError(reqCtx, w, err) {
		return
	} else if !isExist {
		http.Error(w, "job posting not found", http.StatusNotFound)
		return
	}

	// jobPostingDetail을 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}

func (c *JobPostingController) Categories(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	categories, err := c.postingService.Categories(reqCtx)
	if httputils.IsError(reqCtx, w, err) {
		return
	}

	// categories를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=86400, public")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(categories)
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}

func (c *JobPostingController) Skills(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	skills, err := c.postingService.Skills(reqCtx)

	if httputils.IsError(reqCtx, w, err) {
		return
	}

	// skills를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=86400, public")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(skills)
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}

func (c *JobPostingController) Companies(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}

	companies, err := c.postingService.Companies(reqCtx, keyword)

	if httputils.IsError(reqCtx, w, err) {
		return
	}

	// companies를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=43200, public")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(companies)
	if httputils.IsError(reqCtx, w, err) {
		return
	}
}
