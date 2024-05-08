package posting

import (
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
)

type JobPostingsRequest struct {
	Page     int32        `json:"page"`
	Size     int32        `json:"size"`
	QueryReq *query.Query `json:"queryReq"`
}

type JobPostingDetailResponse struct {
	Site           string   `json:"site"`
	PostingId      string   `json:"postingId"`
	Title          string   `json:"title"`
	Skills         []string `json:"skills"`
	Intro          string   `json:"intro"`
	MainTask       string   `json:"mainTask"`
	Qualifications string   `json:"qualifications"`
	Preferred      string   `json:"preferred"`
	Benefits       string   `json:"benefits"`
	RecruitProcess *string  `json:"recruitProcess"`
	CareerMin      *int32   `json:"careerMin"`
	CareerMax      *int32   `json:"careerMax"`
	Addresses      []string `json:"addresses"`
	CompanyId      string   `json:"companyId"`
	CompanyName    string   `json:"companyName"`
	CompanyImages  []string `json:"companyImages"`
	Tags           []string `json:"tags"`
	IsScrapped     bool     `json:"isScrapped"`
	ScrapTags      []string `json:"scrapTags"`
}

func ConvertGrpcToJobPostingDetail(jobPosting *restapi_grpc.JobPostingDetailResponse) *JobPostingDetailResponse {
	return &JobPostingDetailResponse{
		Site:           jobPosting.Site,
		PostingId:      jobPosting.PostingId,
		Title:          jobPosting.Title,
		Skills:         jobPosting.Skills,
		Intro:          jobPosting.Intro,
		MainTask:       jobPosting.MainTask,
		Qualifications: jobPosting.Qualifications,
		Preferred:      jobPosting.Preferred,
		Benefits:       jobPosting.Benefits,
		RecruitProcess: jobPosting.RecruitProcess,
		CareerMin:      jobPosting.CareerMin,
		CareerMax:      jobPosting.CareerMax,
		Addresses:      jobPosting.Addresses,
		CompanyId:      jobPosting.CompanyId,
		CompanyName:    jobPosting.CompanyName,
		CompanyImages:  jobPosting.CompanyImages,
		Tags:           jobPosting.Tags,
		IsScrapped:     false,
		ScrapTags:      []string{},
	}
}
