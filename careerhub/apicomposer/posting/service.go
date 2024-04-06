package posting

import (
	"context"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostingService interface {
	JobPostings(ctx context.Context, req *JobPostingsRequest) (*restapi_grpc.JobPostingsResponse, error)
	JobPostingDetail(ctx context.Context, req *restapi_grpc.JobPostingDetailRequest) (*restapi_grpc.JobPostingDetailResponse, error)
	Categories(ctx context.Context) (*restapi_grpc.CategoriesResponse, error)
	Skills(ctx context.Context) (*restapi_grpc.SkillsResponse, error)
}

type PostingServiceImpl struct {
	postingClient restapi_grpc.RestApiGrpcClient
}

func NewPostingService(postingClient restapi_grpc.RestApiGrpcClient) PostingService {
	return &PostingServiceImpl{
		postingClient: postingClient,
	}
}

func (s *PostingServiceImpl) JobPostings(ctx context.Context, req *JobPostingsRequest) (*restapi_grpc.JobPostingsResponse, error) {
	pbCategories := make([]*restapi_grpc.CategoryQueryReq, len(req.QueryReq.Categories))
	for i, category := range req.QueryReq.Categories {
		pbCategories[i] = &restapi_grpc.CategoryQueryReq{
			Site:         category.Site,
			CategoryName: category.CategoryName,
		}
	}

	pbSkillNames := make([]*restapi_grpc.SkillQueryReq, len(req.QueryReq.SkillNames))
	for i, skillNames := range req.QueryReq.SkillNames {
		pbSkillNames[i] = &restapi_grpc.SkillQueryReq{
			Or: skillNames,
		}
	}

	return s.postingClient.JobPostings(ctx, &restapi_grpc.JobPostingsRequest{
		Page: req.Page,
		Size: req.Size,
		QueryReq: &restapi_grpc.QueryReq{
			Categories: pbCategories,
			SkillNames: pbSkillNames,
			MinCareer:  req.QueryReq.MinCareer,
			MaxCareer:  req.QueryReq.MaxCareer,
		},
	})
}

func (s *PostingServiceImpl) JobPostingDetail(ctx context.Context, req *restapi_grpc.JobPostingDetailRequest) (*restapi_grpc.JobPostingDetailResponse, error) {
	return s.postingClient.JobPostingDetail(ctx, req)
}

func (s *PostingServiceImpl) Categories(ctx context.Context) (*restapi_grpc.CategoriesResponse, error) {
	return s.postingClient.Categories(ctx, &emptypb.Empty{})
}

func (s *PostingServiceImpl) Skills(ctx context.Context) (*restapi_grpc.SkillsResponse, error) {
	return s.postingClient.Skills(ctx, &emptypb.Empty{})
}
