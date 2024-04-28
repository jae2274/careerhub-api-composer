package userinfo

import (
	"context"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
)

type UserinfoService interface {
	FindMatchJob(ctx context.Context, userId string) (*GetMatchJobResponse, error)
	AddCondition(ctx context.Context, userId string, limitCount uint32, req *restapi_grpc.AddConditionReq) (*restapi_grpc.IsSuccessResponse, error)
	UpdateCondition(ctx context.Context, userId string, req *restapi_grpc.Condition) (*restapi_grpc.IsSuccessResponse, error)
	DeleteCondition(ctx context.Context, userId string, conditionId string) (*restapi_grpc.IsSuccessResponse, error)
	UpdateAgreeToMail(ctx context.Context, userId string, agreeToMail bool) (*restapi_grpc.IsSuccessResponse, error)
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
		skills := make([]*Skill, len(condition.Query.SkillNames))
		for j, skill := range condition.Query.SkillNames {
			skills[j] = &Skill{
				Or: skill.Or,
			}
		}

		categories := make([]*Category, len(condition.Query.Categories))
		for j, category := range condition.Query.Categories {
			categories[j] = &Category{
				Site:         category.Site,
				CategoryName: category.CategoryName,
			}
		}

		conditions[i] = &Condition{
			ConditionId:   condition.ConditionId,
			ConditionName: condition.ConditionName,
			Query: &Query{
				MinCareer:  condition.Query.MinCareer,
				MaxCareer:  condition.Query.MaxCareer,
				SkillNames: skills,
				Categories: categories,
			},
		}
	}

	return conditions
}

func (s *UserinfoServiceImpl) AddCondition(ctx context.Context, userId string, limitCount uint32, req *restapi_grpc.AddConditionReq) (*restapi_grpc.IsSuccessResponse, error) {

	return s.userinfoClient.AddCondition(ctx, &restapi_grpc.AddConditionRequest{
		UserId:     userId,
		LimitCount: limitCount,
		Condition:  req,
	})
}

func (s *UserinfoServiceImpl) UpdateCondition(ctx context.Context, userId string, req *restapi_grpc.Condition) (*restapi_grpc.IsSuccessResponse, error) {
	return s.userinfoClient.UpdateCondition(ctx, &restapi_grpc.UpdateConditionRequest{
		UserId:    userId,
		Condition: req,
	})
}

func (s *UserinfoServiceImpl) DeleteCondition(ctx context.Context, userId string, conditionId string) (*restapi_grpc.IsSuccessResponse, error) {
	return s.userinfoClient.DeleteCondition(ctx, &restapi_grpc.DeleteConditionRequest{
		UserId:      userId,
		ConditionId: conditionId,
	})
}

func (s *UserinfoServiceImpl) UpdateAgreeToMail(ctx context.Context, userId string, agreeToMail bool) (*restapi_grpc.IsSuccessResponse, error) {
	return s.userinfoClient.UpdateAgreeToMail(ctx, &restapi_grpc.UpdateAgreeToMailRequest{
		UserId:      userId,
		AgreeToMail: agreeToMail,
	})
}
