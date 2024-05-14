package dto

import (
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/domain"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
)

type JobPostingsResponse struct {
	JobPostings []*domain.JobPosting `json:"jobPostings"`
}

type JobPostingsRequest struct {
	Page     int32        `json:"page"`
	Size     int32        `json:"size"`
	QueryReq *query.Query `json:"queryReq"`
}
