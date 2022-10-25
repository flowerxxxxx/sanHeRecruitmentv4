package mysqlModel

import "time"

// Article 数据库Articles表
type Article struct {
	ArtId       int       `json:"art_id" gorm:"primary_key"`
	Content     string    `json:"content"`
	JobLabel    string    `json:"job_label"`
	Title       string    `json:"title"`
	CreateTime  time.Time `json:"create_time"`
	View        int       `json:"view"`
	Collect     int       `json:"collect"`
	DeliveryNum int       `json:"delivery_num"`
	Weight      float64   `json:"weight"`
	Status      int       `json:"status"` //管理员申请状态 0待审核 1通过
	UpdateTime  time.Time `json:"update_time"`
	Region      string    `json:"region"`
	SalaryMin   int       `json:"salary_min"`
	SalaryMax   int       `json:"salary_max"`
	Show        int       `json:"show"`
	CareerJobId int       `json:"career_job_id"`
	BossId      int       `json:"boss_id"`
	CompanyId   int       `json:"company_id"`
	Tags        string    `json:"tags"`
	Recommend   int       `json:"recommend"`
	ArtType     string    `json:"art_type"`
	//Types      []Type `gorm:"ForeignKey:ArtId;AssociationForeignKey:Id"`
}

// ArtWeightCount 招聘权值基数结构体
type ArtWeightCount struct {
	ArtID       int       `json:"id"`
	View        int       `json:"view"`
	Collect     int       `json:"collect"`
	DeliveryNum int       `json:"delivery_num"`
	UpdateTime  time.Time `json:"update_time"`
	PersonScale string    `json:"personal_scale"`
}

type FuzzyArtInfo struct {
	ArtId int    `json:"art_id" gorm:"primary_key"`
	Title string `json:"title"`
}

type MaxRecommendCount struct {
	MaxNum int `json:"max_num"`
}

// 多字段查找
//func FindArticles() []Article {
//	var articles []Article
//	//sql := "SELECT * FROM `articles`  WHERE id in (select art_id from types where type in ('前端','后端'))"
//	sql2 := "SELECT articles.id,articles.content,articles.create_time,articles.`view`,users.username,users.head_pic,users.role,users.company,users.vip FROM `articles` LEFT JOIN users on articles.boss_id = users.user_id"
//	dao.DB.Debug().Preload("Types").Raw(sql2).Find(&articles)
//	fmt.Println(articles)
//	return articles
//}
