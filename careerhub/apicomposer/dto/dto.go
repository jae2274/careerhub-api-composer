package dto

import "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"

type JobPostingsResponse struct {
	JobPostings []*JobPostingRes `json:"jobPostings"`
}

type JobPostingRes struct {
	Site        string   `json:"site"`
	PostingId   string   `json:"postingId"`
	Title       string   `json:"title"`
	CompanyName string   `json:"companyName"`
	Skills      []string `json:"skills"`
	Categories  []string `json:"categories"`
	ImageUrl    *string  `json:"imageUrl"`
	Addresses   []string `json:"addresses"`
	MinCareer   *int32   `json:"minCareer"`
	MaxCareer   *int32   `json:"maxCareer"`
	IsScrapped  bool     `json:"isScrapped"`
}

func ConvertGrpcToJobPostingRes(jobPosting *restapi_grpc.JobPostingRes) *JobPostingRes {
	return &JobPostingRes{
		Site:        jobPosting.Site,
		PostingId:   jobPosting.PostingId,
		Title:       jobPosting.Title,
		CompanyName: jobPosting.CompanyName,
		Skills:      jobPosting.Skills,
		Categories:  jobPosting.Categories,
		ImageUrl:    jobPosting.ImageUrl,
		Addresses:   jobPosting.Addresses,
		MinCareer:   jobPosting.MinCareer,
		MaxCareer:   jobPosting.MaxCareer,
		IsScrapped:  true,
	}
}

func ConvertGrpcToJobPostingResList(jobPostings []*restapi_grpc.JobPostingRes) []*JobPostingRes {
	jobPostingResList := make([]*JobPostingRes, len(jobPostings))
	for i, jobPosting := range jobPostings {
		jobPostingResList[i] = ConvertGrpcToJobPostingRes(jobPosting)
	}
	return jobPostingResList
}
