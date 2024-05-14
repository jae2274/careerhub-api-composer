package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/httputils"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
)

type ScrapJobController struct {
	scrapJobSvc *userinfo.ScrapJobService
}

func NewScrapJobController(scrapJobSvc *userinfo.ScrapJobService) *ScrapJobController {
	return &ScrapJobController{
		scrapJobSvc: scrapJobSvc,
	}
}

// RegisterRoutes registers the routes for the ScrapJobController
func (c *ScrapJobController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/scrap-job", c.GetScrapJobs).Methods("GET")
	router.HandleFunc("/scrap-job", c.AddScrapJob).Methods("POST")
	router.HandleFunc("/scrap-job", c.RemoveScrapJob).Methods("DELETE")
	router.HandleFunc("/scrap-job/untagged", c.GetUntaggedScrapJobs).Methods("GET")
	router.HandleFunc("/scrap-job/tags", c.GetScrapTags).Methods("GET")
	router.HandleFunc("/scrap-job/tags", c.AddTag).Methods("POST")
	router.HandleFunc("/scrap-job/tags", c.RemoveTag).Methods("DELETE")
}

func (c *ScrapJobController) GetScrapJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	tag := r.URL.Query().Get("tag")
	var tagPtr *string
	if tag != "" {
		tagPtr = &tag
	}

	jobPostings, err := c.scrapJobSvc.GetScrapJobs(ctx, claims.UserId, tagPtr)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(jobPostings)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

func (c *ScrapJobController) AddScrapJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var reqBody restapi_grpc.AddScrapJobRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}
	reqBody.UserId = claims.UserId
	err = c.scrapJobSvc.AddScrapJob(ctx, &reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ScrapJobController) RemoveScrapJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var reqBody restapi_grpc.RemoveScrapJobRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}
	reqBody.UserId = claims.UserId
	isExisted, err := c.scrapJobSvc.RemoveScrapJob(ctx, &reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "Scrap job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ScrapJobController) GetScrapTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	res, err := c.scrapJobSvc.GetScrapTags(ctx, claims.UserId)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

func (c *ScrapJobController) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var reqBody restapi_grpc.AddTagRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}
	reqBody.UserId = claims.UserId
	isExisted, err := c.scrapJobSvc.AddTag(ctx, &reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "Scrap job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ScrapJobController) RemoveTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var reqBody restapi_grpc.RemoveTagRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}
	reqBody.UserId = claims.UserId
	isExisted, err := c.scrapJobSvc.RemoveTag(ctx, &reqBody)
	if httputils.IsError(ctx, w, err) {
		return
	}

	if !isExisted {
		http.Error(w, "Scrap job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ScrapJobController) GetUntaggedScrapJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	res, err := c.scrapJobSvc.GetUntaggedScrapJobs(ctx, claims.UserId)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if httputils.IsError(ctx, w, err) {
		return
	}
}
