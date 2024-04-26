package userinfo

import (
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
)

type UserinfoServiceImpl struct {
	userinfoClient restapi_grpc.RestApiGrpcClient
}

func NewUserinfoService(userinfoClient restapi_grpc.RestApiGrpcClient) *UserinfoServiceImpl {
	return &UserinfoServiceImpl{
		userinfoClient: userinfoClient,
	}
}

// func (s *UserinfoServiceImpl) FindMatchJob(ctx context.Context, req *UserinfoRequest) (*UserinfoResponse, error) {
// 	resp, err := s.userinfoClient.FindMatchJob(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &UserinfoResponse{
// 		MatchJob: resp.MatchJob,
// 	}, nil
// }
