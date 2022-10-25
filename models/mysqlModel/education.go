package mysqlModel

type Education struct {
	ID          int    `json:"id" map:"id" gorm:"primary_key"`
	Username    string `json:"-" map:"username"`
	School      string `json:"school" map:"school"`
	Major       string `json:"major" map:"major"`
	StartTime   string `json:"start_time" map:"start_time"`
	EndTime     string `json:"end_time" map:"end_time"`
	Degree      string `json:"degree" map:"degree"`
	DegreeLevel int    `json:"degree_level"`
}

type EmployEduInfo struct {
	Username string `json:"username" map:"username"`
	School   string `json:"school" map:"school"`
	Major    string `json:"major" map:"major"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Degree   string `json:"degree" map:"degree"`
}

var DegreeWeight = map[string]int{
	"小学": 1,
	"初中": 2,
	"高中": 3,
	"专科": 4,
	"本科": 5,
	"硕士": 6,
	"博士": 7,
}
