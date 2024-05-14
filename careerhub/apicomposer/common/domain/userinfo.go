package domain

import "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"

type ScrapJob struct {
	Site      string   `json:"site"`
	PostingId string   `json:"postingId"`
	Tags      []string `json:"tags"`
}

func ConvertGrpcToScrapJob(scrapJob *restapi_grpc.ScrapJob) *ScrapJob {
	return &ScrapJob{
		Site:      scrapJob.Site,
		PostingId: scrapJob.PostingId,
		Tags:      scrapJob.Tags,
	}
}

func ConvertGrpcToScrapJobs(scrapJobs []*restapi_grpc.ScrapJob) []*ScrapJob {
	result := make([]*ScrapJob, len(scrapJobs))
	for i, scrapJob := range scrapJobs {
		result[i] = ConvertGrpcToScrapJob(scrapJob)
	}
	return result
}
