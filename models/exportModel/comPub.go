package exportModel

type ComPubExport struct {
	CompanyName string `json:"company_name"`
	PersonScale string `json:"personal_scale"`
	Title       string `json:"title"`
	View        int    `json:"view"`
	Nickname    string `json:"nickname"`
	Name        string `json:"name"`
	Tags        string `json:"tags"`
	SalaryMin   string `json:"salary_min"`
	SalaryMax   string `json:"salary_max"`
	Content     string `json:"content"`
	ArtType     string `json:"art_type"`
	Status      string `json:"status"`
}
