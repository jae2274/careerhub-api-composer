package userinfo

import (
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
)

type AddConditionRequest struct {
	ConditionName string       `json:"conditionName"`
	Query         *query.Query `json:"query"`
}

type UpdateConditionRequest struct {
	ConditionId   string       `json:"conditionId"`
	ConditionName string       `json:"conditionName"`
	Query         *query.Query `json:"query"`
}

type IsSuccessResponse struct {
	IsSuccess bool `json:"isSuccess"`
}
