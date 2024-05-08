package posting

import (
	"context"
	"fmt"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/dto"
	postingGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	scrapGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
	"github.com/jae2274/goutils/cchan/async"
	"github.com/jae2274/goutils/terr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostingService interface {
	JobPostings(ctx context.Context, userId *string, req *JobPostingsRequest) (*dto.JobPostingsResponse, error)
	JobPostingDetail(ctx context.Context, userId *string, req *postingGrpc.JobPostingDetailRequest) (*JobPostingDetailResponse, error)
	Categories(ctx context.Context) (*postingGrpc.CategoriesResponse, error)
	Skills(ctx context.Context) (*postingGrpc.SkillsResponse, error)
}

type PostingServiceImpl struct {
	postingClient postingGrpc.RestApiGrpcClient
	scrapClient   scrapGrpc.ScrapJobGrpcClient
}

func NewPostingService(postingClient postingGrpc.RestApiGrpcClient, scrapClient scrapGrpc.ScrapJobGrpcClient) PostingService {
	return &PostingServiceImpl{
		postingClient: postingClient,
		scrapClient:   scrapClient,
	}
}

func (s *PostingServiceImpl) JobPostings(ctx context.Context, userId *string, req *JobPostingsRequest) (*dto.JobPostingsResponse, error) {
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

	jobPostingResList := dto.ConvertGrpcToJobPostingResList(jobPostings.JobPostings)
	if userId != nil {
		err = s.attachIsScrapped(ctx, *userId, jobPostingResList)
		if err != nil {
			return nil, err
		}
	}

	return &dto.JobPostingsResponse{
		JobPostings: jobPostingResList,
	}, nil
}

func (c *PostingServiceImpl) attachIsScrapped(ctx context.Context, userId string, jobPostings []*dto.JobPostingRes) error {
	if len(jobPostings) == 0 {
		return nil
	}

	jobPostingIds := make([]*scrapGrpc.JobPostingId, len(jobPostings))
	for i, jobPosting := range jobPostings {
		jobPostingIds[i] = &scrapGrpc.JobPostingId{
			Site:      jobPosting.Site,
			PostingId: jobPosting.PostingId,
		}
	}

	scrapJobRes, err := c.scrapClient.GetScrapJobsById(ctx, &scrapGrpc.GetScrapJobsByIdRequest{
		UserId:        userId,
		JobPostingIds: jobPostingIds,
	})
	if err != nil {
		return err
	}

	dto.AttachScrapped(jobPostings, scrapJobRes.ScrapJobs)

	return nil
}

func (s *PostingServiceImpl) JobPostingDetail(ctx context.Context, userId *string, req *postingGrpc.JobPostingDetailRequest) (*JobPostingDetailResponse, error) {

	postingChan := async.ExecAsync(func() (*JobPostingDetailResponse, error) {
		res, err := s.postingClient.JobPostingDetail(ctx, req)
		if err != nil {
			return nil, err
		}

		return ConvertGrpcToJobPostingDetail(res), nil
	})

	scrapChan := async.ExecAsync(func() (*scrapGrpc.ScrapJob, error) {
		if userId != nil {
			scrapJobRes, err := s.scrapClient.GetScrapJobsById(ctx, &scrapGrpc.GetScrapJobsByIdRequest{
				UserId: *userId,
				JobPostingIds: []*scrapGrpc.JobPostingId{
					{
						Site:      req.Site,
						PostingId: req.PostingId,
					},
				},
			})
			if err != nil {
				return nil, err
			}

			if len(scrapJobRes.ScrapJobs) > 1 {
				return nil, terr.Wrap(fmt.Errorf("scrap job response has more than one job. site: %s, postingId: %s", req.Site, req.PostingId))
			}

			if len(scrapJobRes.ScrapJobs) == 1 {
				return scrapJobRes.ScrapJobs[0], nil
			}
		}

		return nil, nil
	})

	jobPostingResult := <-postingChan
	if jobPostingResult.Err != nil {
		return nil, jobPostingResult.Err
	}
	response := jobPostingResult.Value

	scrapJobResult := <-scrapChan
	if scrapJobResult.Err != nil {
		return nil, scrapJobResult.Err
	}

	if scrapJobResult.Value != nil {
		response.IsScrapped = true
		response.ScrapTags = scrapJobResult.Value.Tags
	}

	return response, nil
}

func (s *PostingServiceImpl) Categories(ctx context.Context) (*postingGrpc.CategoriesResponse, error) {
	return s.postingClient.Categories(ctx, &emptypb.Empty{})
}

func (s *PostingServiceImpl) Skills(ctx context.Context) (*postingGrpc.SkillsResponse, error) {
	return s.postingClient.Skills(ctx, &emptypb.Empty{})
}
