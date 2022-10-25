package mysqlModel

type Collection struct {
	Id             string `json:"id" gorm:"primary_key"`
	Username       string `json:"username"`
	Art_id         int    `json:"art_id"`
	Classification string `json:"classification"`
}
