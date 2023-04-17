package recommendBiz

import (
	"fmt"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/service/mysqlService"
	"strconv"
	"time"
)

var JobService *mysqlService.JobService
var ArticleService *mysqlService.ArticleService

// CountArcHot 计算权重
// 招聘 =(观看*0.4+收藏*0.3-投递量*2)*提权 {（最新发布（一天内发布）+公司规模）提权[总值*1.5*公司规模权重]}
func CountArcHot(aw mysqlModel.ArtWeightCount) (artWeight float64) {
	artWeight = (float64(aw.View)*0.4 +
		float64(aw.Collect)*0.3 -
		float64(aw.DeliveryNum)*2) *
		companyWeight[aw.PersonScale]
	if aw.UpdateTime.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		artWeight = artWeight * 1.5
	}
	artWeight, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", artWeight), 64)
	return
}

// DealArtRecommendWeight 处理招聘权重
func DealArtRecommendWeight(artId string) {
	artIdInt, _ := strconv.Atoi(artId)
	aw := JobService.QueryArtWeight(artId)
	artWeight := CountArcHot(aw)
	//fmt.Println(artWeight)
	ArticleService.UpdateArtWeight(artIdInt, artWeight)
}
