package userinfo

type GetMatchJobResponse struct {
	Conditions  []*Condition `json:"conditions"`
	AgreeToMail bool         `json:"agreeToMail"`
}

type Condition struct {
	ConditionId   string `json:"conditionId"`
	ConditionName string `json:"conditionName"`
	Query         *Query `json:"query"`
}

type Query struct {
	Categories []*Category `json:"categories"`
	SkillNames []*Skill    `json:"skillNames"`
	MinCareer  *int32      `json:"minCareer"`
	MaxCareer  *int32      `json:"maxCareer"`
}

type Category struct {
	Site         string `json:"site"`
	CategoryName string `json:"categoryName"`
}

type Skill struct {
	Or []string `json:"Or"`
}
