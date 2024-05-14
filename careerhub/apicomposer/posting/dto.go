package posting

import "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"

type JobPostingsRequest struct {
	Page     int32        `json:"page"`
	Size     int32        `json:"size"`
	QueryReq *query.Query `json:"queryReq"`
}
