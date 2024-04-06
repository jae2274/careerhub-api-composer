package posting

type JobPostingsRequest struct {
	Page     int32     `json:"page"`
	Size     int32     `json:"size"`
	QueryReq *QueryReq `json:"queryReq"`
}

type QueryReq struct {
	Categories []*CategoryQueryReq
	SkillNames [][]string
	MinCareer  *int32
	MaxCareer  *int32
}

type CategoryQueryReq struct {
	Site         string
	CategoryName string
}
