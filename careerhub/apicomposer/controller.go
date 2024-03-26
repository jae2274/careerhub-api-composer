package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"github.com/jae2274/goutils/llog"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	initPage                   = 0
	messageInternalServerError = "Internal Server Error"
)

type Controller struct {
	postingClient restapi_grpc.RestApiGrpcClient
	router        *mux.Router
}

func NewController(postingClient restapi_grpc.RestApiGrpcClient, router *mux.Router) *Controller {
	return &Controller{
		postingClient: postingClient,
		router:        router,
	}
}

func (c *Controller) RegisterRoutes(rootPath string) {
	c.router.HandleFunc(rootPath+"/job_postings", c.JobPostings).Methods("GET")
	c.router.HandleFunc(rootPath+"/job_postings/{site}/{postingId}", c.JobPostingDetail).Methods("GET")
	c.router.HandleFunc(rootPath+"/categories", c.Categories).Methods("GET")
	c.router.HandleFunc(rootPath+"/skills", c.Skills).Methods("GET")
}

func (c *Controller) JobPostings(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	req, err := posting.ExtractJobPostingsRequest(r, initPage)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.postingClient.JobPostings(reqCtx, req)
	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// jobPostings를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: 이후 세부적으로 설정 필요
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

	resp, err := c.postingClient.JobPostingDetail(reqCtx, req)

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// jobPostingDetail을 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: 이후 세부적으로 설정 필요
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) Categories(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	categories, err := c.postingClient.Categories(reqCtx, &emptypb.Empty{})
	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// categories를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: 이후 세부적으로 설정 필요
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (c *Controller) Skills(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	skills, err := c.postingClient.Skills(reqCtx, &emptypb.Empty{})

	if err != nil {
		llog.LogErr(reqCtx, err)
		http.Error(w, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	// skills를 JSON으로 변환하여 응답
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: 이후 세부적으로 설정 필요
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(skills)
}
