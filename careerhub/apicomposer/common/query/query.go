package query

type Query struct {
	Categories []*CategoryQuery `json:"categories"`
	SkillNames [][]string       `json:"skillNames"`
	MinCareer  *int32           `json:"minCareer"`
	MaxCareer  *int32           `json:"maxCareer"`
}

type CategoryQuery struct {
	Site         string `json:"site"`
	CategoryName string `json:"categoryName"`
}
