package userBind

type FuzzyQueryJobs struct {
	FuzzyName string `json:"fuzzy_name"`
	QueryType string `json:"query_type"`
	PageNum   int    `json:"page_num"`
}

type FuzzyQueryComs struct {
	FuzzyName string `json:"fuzzy_name"`
	PageNum   int    `json:"page_num"`
}
