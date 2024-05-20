package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	router.HandleFunc("/company-review/reviews", c.GetReviews).Methods("GET")
}

func getPage(urlValues url.Values) (int, error) {
	pageStr := urlValues.Get("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < review.InitPage {
			return 0, fmt.Errorf("invalid page number. %s", pageStr)
		}

		return page, nil
	}

	return review.InitPage, nil
}

func getPageSize(urlValues url.Values) (int, error) {
	sizeStr := urlValues.Get("size")
	if sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 0 {
			return 0, fmt.Errorf("invalid page size. %s", sizeStr)
		}

		return size, nil
	}

	return review.DefaultSize, nil
}

func (c *ReviewController) GetReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, err := getPage(r.URL.Query())
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Invalid page", http.StatusBadRequest)
		return
	}

	size, err := getPageSize(r.URL.Query())
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	companyName := r.URL.Query().Get("companyName")
	llog.Info(ctx, fmt.Sprintf("companyName: %s", companyName))
	if companyName == "" {
		err := fmt.Errorf("companyName is required")
		llog.LogErr(ctx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
