package service

import (
	"context"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/dto"
	postingGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	scrapJobGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
	"google.golang.org/grpc"
)

type ScrapJobService struct {
	scrapjobClient scrapJobGrpc.ScrapJobGrpcClient
	postingClient  postingGrpc.RestApiGrpcClient
}

func NewScrapJobService(
	scrapjobClient scrapJobGrpc.ScrapJobGrpcClient,
	postingClient postingGrpc.RestApiGrpcClient,
) *ScrapJobService {
	return &ScrapJobService{
		scrapjobClient: scrapjobClient,
		postingClient:  postingClient,
	}
}

func (s *ScrapJobService) GetScrapJobs(ctx context.Context, userId string, tag *string) (*dto.JobPostingsResponse, error) {
	scrapJobRes, err := s.scrapjobClient.GetScrapJobs(ctx, &scrapJobGrpc.GetScrapJobsRequest{
		UserId: userId,
		Tag:    tag,
	})

	if err != nil {
		return nil, err
	}

	if len(scrapJobRes.ScrapJobs) == 0 {
		return &dto.JobPostingsResponse{
			JobPostings: []*dto.JobPostingRes{},
		}, nil
	}

	return s.getJobPostings(ctx, scrapJobRes.ScrapJobs)
}

func (s *ScrapJobService) getJobPostings(ctx context.Context, scrapJobs []*scrapJobGrpc.ScrapJob) (*dto.JobPostingsResponse, error) {
	postingIds := make([]*postingGrpc.JobPostingIdReq, len(scrapJobs))
	for i, scrapJob := range scrapJobs {
		postingIds[i] = &postingGrpc.JobPostingIdReq{
			Site:      scrapJob.Site,
			PostingId: scrapJob.PostingId,
		}
	}

	postings, err := s.postingClient.JobPostingsById(ctx, &postingGrpc.JobPostingsByIdRequest{
		JobPostingIds: postingIds,
	})

	if err != nil {
		return nil, err
	}

	jobPostings := dto.ConvertGrpcToJobPostingResList(postings.JobPostings)
	dto.AttachScrapped(jobPostings, scrapJobs)

	return &dto.JobPostingsResponse{
		JobPostings: jobPostings,
	}, nil
}

func (s *ScrapJobService) AddScrapJob(ctx context.Context, in *scrapJobGrpc.AddScrapJobRequest, opts ...grpc.CallOption) error {
	_, err := s.scrapjobClient.AddScrapJob(ctx, in)
	return err
}
func (s *ScrapJobService) RemoveScrapJob(ctx context.Context, in *scrapJobGrpc.RemoveScrapJobRequest, opts ...grpc.CallOption) (bool, error) {
	isExisted, err := s.scrapjobClient.RemoveScrapJob(ctx, in)

	return isExisted.IsExisted, err
}
func (s *ScrapJobService) AddTag(ctx context.Context, in *scrapJobGrpc.AddTagRequest, opts ...grpc.CallOption) (bool, error) {
	isExisted, err := s.scrapjobClient.AddTag(ctx, in)

	return isExisted.IsExisted, err
}
func (s *ScrapJobService) RemoveTag(ctx context.Context, in *scrapJobGrpc.RemoveTagRequest, opts ...grpc.CallOption) (bool, error) {
	isExisted, err := s.scrapjobClient.RemoveTag(ctx, in)

	return isExisted.IsExisted, err
}
func (s *ScrapJobService) GetScrapTags(ctx context.Context, userId string, opts ...grpc.CallOption) (*scrapJobGrpc.GetScrapTagsResponse, error) {
	return s.scrapjobClient.GetScrapTags(ctx, &scrapJobGrpc.GetScrapTagsRequest{
		UserId: userId,
	})
}

func (s *ScrapJobService) GetUntaggedScrapJobs(ctx context.Context, userId string) (*dto.JobPostingsResponse, error) {
	scrapJobRes, err := s.scrapjobClient.GetUntaggedScrapJobs(ctx, &scrapJobGrpc.GetUntaggedScrapJobsRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	if len(scrapJobRes.ScrapJobs) == 0 {
		return &dto.JobPostingsResponse{
			JobPostings: []*dto.JobPostingRes{},
		}, nil
	}

	return s.getJobPostings(ctx, scrapJobRes.ScrapJobs)
}
