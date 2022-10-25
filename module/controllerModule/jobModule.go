package controllerModule

import (
	"encoding/json"
	"log"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/e"
	"strconv"
	"time"
)

type JobConModule struct {
}

// SaveRecruitInfoToRedis 保存详细页面信息到redis
func (jcm JobConModule) SaveRecruitInfoToRedis(artId string, recInfo mysqlModel.OneArticleOut) {
	recInfoByte, errMar := json.Marshal(recInfo)
	if errMar != nil {
		log.Println("SaveRecruitInfoToRedis Marshal failed,err:", errMar)
		return
	}
	RISaveTime := 120 * time.Second
	hotCount := dao.Redis.Get("RecruitInfo_count_" + artId).Val()
	if hotCount == "" {
		if errReSet := dao.Redis.Set("RecruitInfo_count_"+artId, 1, 600*time.Second).Err(); errReSet != nil {
			log.Println("SaveRecruitInfoToRedis set failed,err:", errReSet)
			return
		}
	} else {
		hotCountInt, _ := strconv.Atoi(hotCount)
		if hotCountInt > 200 {
			RISaveTime = 300 * time.Second
		}
	}
	if errReSet := dao.Redis.Set("RecruitInfo_"+artId, recInfoByte, RISaveTime).Err(); errReSet != nil {
		log.Println("SaveRecruitInfoToRedis set failed,err:", errReSet)
		return
	}
}

// GetRecruitInfoFromRedis 详细页面信息尝试命中redis
func (jcm JobConModule) GetRecruitInfoFromRedis(artId string) (mysqlModel.OneArticleOut, error) {
	recInfo := mysqlModel.OneArticleOut{}
	redisGet := dao.Redis.Get("RecruitInfo_" + artId).Val()
	if redisGet == "" {
		return recInfo, e.RedisNoVal
	}
	dao.Redis.Incr("RecruitInfo_count_" + artId)
	errUmMar := json.Unmarshal([]byte(redisGet), &recInfo)
	if errUmMar != nil {
		log.Println("GetRecruitInfoFromRedis Unmarshal failed,err:", errUmMar)
		return recInfo, e.UnMarshalErr
	}
	return recInfo, nil
}
