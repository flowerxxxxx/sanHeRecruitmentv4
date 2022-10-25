package mysqlModel

type DailySaver struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ArtId       int    `json:"art_id"`
	DayView     int    `json:"day_view"`
	DayDelivery int    `json:"day_delivery"`
	Date        string `json:"date"`
}

type DailySaverSqlOut struct {
	Id          int    `json:"id" gorm:"primary_key"`
	ArtId       int    `json:"art_id"`
	DayView     int    `json:"day_view"`
	DayDelivery int    `json:"day_delivery"`
	Date        string `json:"date"`
}

// DailyInfo 单日信息查看
type DailyInfo struct {
	DailySaverSqlOut
	Title string `json:"title"`
}

type DailyLabelInfo struct {
	DayView     int    `json:"day_view"`
	DayDelivery int    `json:"day_delivery"`
	Label       string `json:"label"`
}
