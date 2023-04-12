package mysqlModel

import "time"

type UserComArticle struct {
	ArtID int    `json:"id"`
	Title string `json:"title"`
	//Content  string `json:"content"`
	View        int      `json:"view"`
	Nickname    string   `json:"nickname"`
	JobLabel    string   `json:"job_label"`
	HeadPic     string   `json:"head_pic"`
	Region      string   `json:"region"`
	SalaryMin   int      `json:"salary_min"`
	SalaryMax   int      `json:"salary_max"`
	CompanyName string   `json:"company_name"`
	PersonScale string   `json:"personal_scale"`
	BossId      int      `json:"boss_id"`
	CareerJobId int      `json:"career_job_id"`
	ComLevel    int      `json:"com_level"`
	Tags        string   `json:"-"`
	TagsOut     []string `json:"tags"`
	Recommend   int      `json:"recommend"`
}

type UserComArticleXls struct {
	Title       string `json:"title"`
	View        int    `json:"view"`
	Nickname    string `json:"nickname"`
	Name        string `json:"name"`
	JobLabel    string `json:"job_label"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
	CompanyName string `json:"company_name"`
	PersonScale string `json:"personal_scale"`
	Tags        string `json:"tags"`
	Content     string `json:"content"`
	ArtType     string `json:"art_type"`
	Status      int    `json:"status"`
}

type UserComArtBoss struct {
	UserComArticle
	Status int `json:"status"`
	Show   int `json:"show"`
}

// OneArticleOut 单个文章详情
type OneArticleOut struct {
	OneArticle
	President     string `json:"president"`
	CreateTimeOut string `json:"create_time"`
}

// OneArticle 单个文章详情
type OneArticle struct {
	UserComArticle
	ComId      int       `json:"com_id"`
	Content    string    `json:"content"`
	PicUrl     string    `json:"pic_url"`
	ArtType    string    `json:"type"`
	CreateTime time.Time `json:"-"`
}

// JobProgress 求职进展页面sql
type JobProgress struct {
	ArtID         int       `json:"art_id"`
	Title         string    `json:"title"`
	Nickname      string    `json:"nickname"`
	HeadPic       string    `json:"head_pic"`
	SalaryMin     int       `json:"salary_min"`
	SalaryMax     int       `json:"salary_max"`
	Qualification int       `json:"qualification"`
	Read          int       `json:"read"`
	President     string    `json:"president"`
	DeliveryTime  time.Time `json:"-"`
}

// JobProgressOut 求职进展页面formatTime
type JobProgressOut struct {
	JobProgress
	DeliveryTimeOut string `json:"delivery_time"`
}
