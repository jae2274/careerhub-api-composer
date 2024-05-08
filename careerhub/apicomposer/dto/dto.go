package dto

import (
	postingGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	scrapJobGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
)

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
	Tags        []string `json:"tags"`
}

func ConvertGrpcToJobPostingRes(jobPosting *postingGrpc.JobPostingRes) *JobPostingRes {
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
	}
}

func ConvertGrpcToJobPostingResList(jobPostings []*postingGrpc.JobPostingRes) []*JobPostingRes {
	jobPostingResList := make([]*JobPostingRes, len(jobPostings))
	for i, jobPosting := range jobPostings {
		jobPostingResList[i] = ConvertGrpcToJobPostingRes(jobPosting)
	}
	return jobPostingResList
}

func AttachScrapped(jobPostings []*JobPostingRes, scrapJobs []*scrapJobGrpc.ScrapJob) {
	scrapJobMap := make(map[string]*scrapJobGrpc.ScrapJob)
	for _, scrapJob := range scrapJobs {
		scrapJobMap[scrapJob.Site+scrapJob.PostingId] = scrapJob
	}

	for _, jobPosting := range jobPostings {
		if scrapJob, ok := scrapJobMap[jobPosting.Site+jobPosting.PostingId]; ok {
			jobPosting.IsScrapped = true
			jobPosting.Tags = scrapJob.Tags
		}
	}
}
