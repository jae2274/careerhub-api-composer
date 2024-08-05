package query

type Query struct {
	Categories []*CategoryQuery   `json:"categories"`
	SkillNames [][]string         `json:"skillNames"`
	MinCareer  *int32             `json:"minCareer"`
	MaxCareer  *int32             `json:"maxCareer"`
	Companies  []SiteCompanyQuery `json:"companies"`
}

type CategoryQuery struct {
	Site         string `json:"site"`
	CategoryName string `json:"categoryName"`
}

type SiteCompanyQuery struct {
	Site        string `json:"site"`
	CompanyName string `json:"companyName"`
}
