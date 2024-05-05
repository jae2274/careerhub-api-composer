package dto

type JobPostingsResponse struct {
	JobPostings []*JobPostingRes `json:"jobPostings"`
}

type JobPostingRes struct {
	Site        string   `json:"site"`
	PostingId   string   `json:"postingId"`
	Title       string   `json:"title"`
	CompanyName string   `json:"companyName"`
	Skills      []string `json:"skills"`
	Categories  []string `json:"categories"`
	ImageUrl    *string  `json:"imageUrl"`
	Addresses   []string `json:"addresses"`
	MinCareer   *int32   `json:"minCareer"`
	MaxCareer   *int32   `json:"maxCareer"`
	IsScrapped  bool     `json:"isScrapped"`
}
