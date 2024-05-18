package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/httputils"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/review"
	"github.com/jae2274/goutils/llog"
)

type ReviewController struct {
	reviewService *review.ReviewService
}

func NewReviewController(reviewService *review.ReviewService) *ReviewController {
	return &ReviewController{
		reviewService: reviewService,
	}
}

func (c *ReviewController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/company-review/{companyName}/reviews", c.GetReviews).Methods("GET")
}

func (c *ReviewController) GetReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	page, size := review.InitPage, review.DefaultSize

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			if err != nil {
				llog.LogErr(ctx, err)
			}

			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 0 {
			if err != nil {
				llog.LogErr(ctx, err)
			}
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}

	vars := mux.Vars(r)
	companyName, ok := vars["companyName"]
	if !ok {
		http.Error(w, "Invalid company name", http.StatusBadRequest)
		return
	}

	reviews, err := c.reviewService.GetReviews(r.Context(), companyName, page, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(reviews)
	if httputils.IsError(ctx, w, err) {
		return
	}
}
