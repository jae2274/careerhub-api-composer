package userinfo

import (
	"context"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
)

type UserinfoService interface {
	FindMatchJob(ctx context.Context, userId string) (*GetMatchJobResponse, error)
	AddCondition(ctx context.Context, userId string, limitCount uint32, req *AddConditionRequest) (*IsSuccessResponse, error)
	UpdateCondition(ctx context.Context, userId string, req *UpdateConditionRequest) (*IsSuccessResponse, error)
	DeleteCondition(ctx context.Context, userId string, conditionId string) (*IsSuccessResponse, error)
	UpdateAgreeToMail(ctx context.Context, userId string, agreeToMail bool) (*IsSuccessResponse, error)
}

type UserinfoServiceImpl struct {
	userinfoClient restapi_grpc.RestApiGrpcClient
}

func NewUserinfoService(userinfoClient restapi_grpc.RestApiGrpcClient) UserinfoService {
	return &UserinfoServiceImpl{
		userinfoClient: userinfoClient,
	}
}

func (s *UserinfoServiceImpl) FindMatchJob(ctx context.Context, userId string) (*GetMatchJobResponse, error) {
	res, err := s.userinfoClient.FindMatchJob(ctx, &restapi_grpc.FindMatchJobRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return &GetMatchJobResponse{
		AgreeToMail: res.AgreeToMail,
		Conditions:  convertToConditions(res),
	}, nil
}

func convertToConditions(res *restapi_grpc.FindMatchJobResponse) []*Condition {
	conditions := make([]*Condition, len(res.Conditions))
	for i, condition := range res.Conditions {
		skills := make([][]string, len(condition.Query.SkillNames))
		for j, skill := range condition.Query.SkillNames {
			skills[j] = skill.Or
		}

		categories := make([]*query.CategoryQuery, len(condition.Query.Categories))
		for j, category := range condition.Query.Categories {
			categories[j] = &query.CategoryQuery{
				Site:         category.Site,
				CategoryName: category.CategoryName,
			}
		}

		conditions[i] = &Condition{
			ConditionId:   condition.ConditionId,
			ConditionName: condition.ConditionName,
			Query: &query.Query{
				MinCareer:  condition.Query.MinCareer,
				MaxCareer:  condition.Query.MaxCareer,
				SkillNames: skills,
				Categories: categories,
			},
		}
	}

	return conditions
}

func (s *UserinfoServiceImpl) AddCondition(ctx context.Context, userId string, limitCount uint32, req *AddConditionRequest) (*IsSuccessResponse, error) {

	res, err := s.userinfoClient.AddCondition(ctx, &restapi_grpc.AddConditionRequest{
		UserId:     userId,
		LimitCount: limitCount,
		Condition: &restapi_grpc.AddConditionReq{
			ConditionName: req.ConditionName,
			Query:         convertQueryToGrpc(req.Query),
		},
	})

	return convertGrpcToIsSuccessResponse(res), err
}

func (s *UserinfoServiceImpl) UpdateCondition(ctx context.Context, userId string, req *UpdateConditionRequest) (*IsSuccessResponse, error) {
	res, err := s.userinfoClient.UpdateCondition(ctx, &restapi_grpc.UpdateConditionRequest{
		UserId: userId,
		Condition: &restapi_grpc.Condition{
			ConditionId:   req.ConditionId,
			ConditionName: req.ConditionName,
			Query:         convertQueryToGrpc(req.Query),
		},
	})

	return convertGrpcToIsSuccessResponse(res), err
}

func convertQueryToGrpc(query *query.Query) *restapi_grpc.Query {
	categories := make([]*restapi_grpc.Category, len(query.Categories))
	for i, category := range query.Categories {
		categories[i] = &restapi_grpc.Category{
			Site:         category.Site,
			CategoryName: category.CategoryName,
		}
	}

	skillNames := make([]*restapi_grpc.Skill, len(query.SkillNames))
	for i, skillName := range query.SkillNames {
		skillNames[i] = &restapi_grpc.Skill{
			Or: skillName,
		}
	}

	return &restapi_grpc.Query{
		Categories: categories,
		SkillNames: skillNames,
		MinCareer:  query.MinCareer,
		MaxCareer:  query.MaxCareer,
	}

}

func (s *UserinfoServiceImpl) DeleteCondition(ctx context.Context, userId string, conditionId string) (*IsSuccessResponse, error) {
	res, err := s.userinfoClient.DeleteCondition(ctx, &restapi_grpc.DeleteConditionRequest{
		UserId:      userId,
		ConditionId: conditionId,
	})

	return convertGrpcToIsSuccessResponse(res), err
}

func (s *UserinfoServiceImpl) UpdateAgreeToMail(ctx context.Context, userId string, agreeToMail bool) (*IsSuccessResponse, error) {
	res, err := s.userinfoClient.UpdateAgreeToMail(ctx, &restapi_grpc.UpdateAgreeToMailRequest{
		UserId:      userId,
		AgreeToMail: agreeToMail,
	})

	return convertGrpcToIsSuccessResponse(res), err
}

func convertGrpcToIsSuccessResponse(res *restapi_grpc.IsSuccessResponse) *IsSuccessResponse {
	if res == nil {
		return nil
	}

	return &IsSuccessResponse{
		IsSuccess: res.IsSuccess,
	}
}
