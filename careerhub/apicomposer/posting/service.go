package posting

import (
	"context"
	"fmt"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/domain"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/dto"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
	postingGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	scrapGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostingService struct {
	postingClient  postingGrpc.RestApiGrpcClient
	scrapjobClient scrapGrpc.ScrapJobGrpcClient
}

func NewPostingService(postingClient postingGrpc.RestApiGrpcClient, scrapjobClient scrapGrpc.ScrapJobGrpcClient) *PostingService {
	return &PostingService{
		postingClient:  postingClient,
		scrapjobClient: scrapjobClient,
	}
}

func (s *PostingService) JobPostingsWithClaims(ctx context.Context, req *dto.JobPostingsRequest, claims *jwtresolver.CustomClaims) ([]*domain.JobPosting, error) {
	jobPostings, err := s.JobPostings(ctx, req)
	if err != nil {
		return nil, err
	}

	scrapJobs, err := s.getScrapJobsByPostingIds(ctx, claims.UserId, domain.GetJobPostingIds(jobPostings))

	if err != nil {
		return nil, err
	}

	domain.AttachScrapped(jobPostings, scrapJobs)

	return jobPostings, nil
}

func (s *PostingService) JobPostings(ctx context.Context, req *dto.JobPostingsRequest) ([]*domain.JobPosting, error) {
	pbCategories := make([]*postingGrpc.CategoryQueryReq, len(req.QueryReq.Categories))
	for i, category := range req.QueryReq.Categories {
		pbCategories[i] = &postingGrpc.CategoryQueryReq{
			Site:         category.Site,
			CategoryName: category.CategoryName,
		}
	}

	pbSkillNames := make([]*postingGrpc.SkillQueryReq, len(req.QueryReq.SkillNames))
	for i, skillNames := range req.QueryReq.SkillNames {
		pbSkillNames[i] = &postingGrpc.SkillQueryReq{
			Or: skillNames,
		}
	}

	jobPostings, err := s.postingClient.JobPostings(ctx, &postingGrpc.JobPostingsRequest{
		Page: req.Page,
		Size: req.Size,
		QueryReq: &postingGrpc.QueryReq{
			Categories: pbCategories,
			SkillNames: pbSkillNames,
			MinCareer:  req.QueryReq.MinCareer,
			MaxCareer:  req.QueryReq.MaxCareer,
		},
	})

	if err != nil {
		return nil, err
	}

	jobPostingResList := domain.ConvertGrpcToJobPostingResList(jobPostings.JobPostings)

	return jobPostingResList, nil
}

func (s *PostingService) getScrapJobsByPostingIds(ctx context.Context, userId string, jobPostings []*domain.JobPostingId) ([]*domain.ScrapJob, error) {
	if len(jobPostings) == 0 {
		return []*domain.ScrapJob{}, nil
	}

	jobPostingIds := make([]*scrapGrpc.JobPostingId, len(jobPostings))
	for i, jobPosting := range jobPostings {
		jobPostingIds[i] = &scrapGrpc.JobPostingId{
			Site:      jobPosting.Site,
			PostingId: jobPosting.PostingId,
		}
	}

	res, err := s.scrapjobClient.GetScrapJobsById(ctx, &scrapGrpc.GetScrapJobsByIdRequest{
		UserId:        userId,
		JobPostingIds: jobPostingIds,
	})
	if err != nil {
		return nil, err
	}

	scrapJobs := make([]*domain.ScrapJob, len(res.ScrapJobs))
	for i, scrapJob := range res.ScrapJobs {
		scrapJobs[i] = domain.ConvertGrpcToScrapJob(scrapJob)
	}

	return scrapJobs, nil
}

func (s *PostingService) JobPostingDetailWithClaims(ctx context.Context, req *postingGrpc.JobPostingDetailRequest, claims *jwtresolver.CustomClaims) (*domain.JobPostingDetail, error) {
	res, err := s.JobPostingDetail(ctx, req)
	if err != nil {
		return nil, err
	}

	scrapJobs, err := s.getScrapJobsByPostingIds(ctx, claims.UserId, []*domain.JobPostingId{{Site: res.Site, PostingId: res.PostingId}})
	if err != nil {
		return nil, err
	}

	if len(scrapJobs) == 0 {
		return res, nil
	}

	if len(scrapJobs) > 1 {
		return nil, fmt.Errorf("multiple scrap jobs found for the same job posting id. site: %s, postingId: %s", res.Site, res.PostingId)
	}

	res.IsScrapped = true
	res.ScrapTags = scrapJobs[0].Tags

	return res, nil
}

func (s *PostingService) JobPostingDetail(ctx context.Context, req *postingGrpc.JobPostingDetailRequest) (*domain.JobPostingDetail, error) {
	res, err := s.postingClient.JobPostingDetail(ctx, req)
	if err != nil {
		return nil, err
	}

	return domain.ConvertGrpcToJobPostingDetail(res), nil
}

func (s *PostingService) Categories(ctx context.Context) (*postingGrpc.CategoriesResponse, error) {
	return s.postingClient.Categories(ctx, &emptypb.Empty{})
}

func (s *PostingService) Skills(ctx context.Context) (*postingGrpc.SkillsResponse, error) {
	return s.postingClient.Skills(ctx, &emptypb.Empty{})
}
