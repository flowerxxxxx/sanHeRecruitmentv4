package userBind

type PubEmployReqBind struct {
	CareerJobId string   `json:"careerJobId"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Region      string   `json:"region"`
	SalaryMin   string   `json:"salaryMin"`
	SalaryMax   string   `json:"salaryMax"`
	TagList     []string `json:"tag_list"`
}
