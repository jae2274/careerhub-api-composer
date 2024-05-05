package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/httputils"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/service"
	"github.com/jae2274/goutils/llog"
)

type MatchJobController struct {
	matchJobSvc service.MatchJobService
}

func NewMatchJobController(matchJobSvc service.MatchJobService) *MatchJobController {
	return &MatchJobController{
		matchJobSvc: matchJobSvc,
	}
}

func (c *MatchJobController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/match-job", c.FindMatchJob).Methods("GET")
	router.HandleFunc("/match-job/condition", c.AddCondition).Methods("POST")
	router.HandleFunc("/match-job/condition", c.UpdateCondition).Methods("PUT")
	router.HandleFunc("/match-job/condition", c.DeleteCondition).Methods("DELETE")
	router.HandleFunc("/match-job/agree-to-mail", c.UpdateAgreeToMail).Methods("PUT")
}

func (c *MatchJobController) FindMatchJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	matchJob, err := c.matchJobSvc.FindMatchJob(ctx, claims.UserId)
	if httputils.IsError(ctx, w, err) {
		return
	}

	str, err := json.Marshal(matchJob)
	if httputils.IsError(ctx, w, err) {
		return
	}
	llog.Info(ctx, string(str))

	err = json.NewEncoder(w).Encode(matchJob)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

const limitCount = 2

func (c *MatchJobController) AddCondition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var req userinfo.AddConditionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if httputils.IsError(ctx, w, err) {
		return
	}

	isSuccess, err := c.matchJobSvc.AddCondition(ctx, claims.UserId, limitCount, &req)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(isSuccess)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

func (c *MatchJobController) UpdateCondition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	var req userinfo.UpdateConditionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if httputils.IsError(ctx, w, err) {
		return
	}

	isSuccess, err := c.matchJobSvc.UpdateCondition(ctx, claims.UserId, &req)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(isSuccess)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

func (c *MatchJobController) DeleteCondition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	conditionIdStruct := struct {
		ConditionId string `json:"conditionId"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&conditionIdStruct)
	if httputils.IsError(ctx, w, err) {
		return
	}

	isSuccess, err := c.matchJobSvc.DeleteCondition(ctx, claims.UserId, conditionIdStruct.ConditionId)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(isSuccess)
	if httputils.IsError(ctx, w, err) {
		return
	}
}

func (c *MatchJobController) UpdateAgreeToMail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetClaims(ctx)
	if httputils.IsNotLoggedIn(ctx, w, ok) {
		return
	}

	agreeToMailStruct := struct {
		AgreeToMail bool `json:"agreeToMail"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&agreeToMailStruct)
	if httputils.IsError(ctx, w, err) {
		return
	}

	isSuccess, err := c.matchJobSvc.UpdateAgreeToMail(ctx, claims.UserId, agreeToMailStruct.AgreeToMail)
	if httputils.IsError(ctx, w, err) {
		return
	}

	err = json.NewEncoder(w).Encode(isSuccess)
	if httputils.IsError(ctx, w, err) {
		return
	}
}
