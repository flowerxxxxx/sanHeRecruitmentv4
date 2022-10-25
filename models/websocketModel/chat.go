package websocketModel

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/util/formatUtil"
	"strconv"
	"time"
)

//添加message_type 字段，1为文本消息，2为图片，考虑增加可发送的文档（3）

type Trainer struct {
	Id           int    `json:"id" gorm:"primary_key"`
	Userid       string `json:"-"`            //用户名
	Content      string `json:"content"`      // 内容
	Start_time   int64  `json:"start_time"`   // 创建时间
	End_time     int64  `json:"-"`            // 过期时间
	Read         int    `json:"read"`         // 已读
	Message_type int    `json:"message_type"` //消息类型
	FromUsername string `json:"-"`            //发送者
	ToUsername   string `json:"-"`            //接受者
}

func InsertMsg(userid string, content string, read int, expire int64) (err error) {
	comment := Trainer{
		Userid:     userid,
		Content:    content,
		Start_time: time.Now().Unix(),
		End_time:   time.Now().Unix() + expire,
		Read:       read,
	}
	err = dao.DB.Save(&comment).Error
	if err != nil {
		return err
	}
	return
}

func InsertMsg2(msg *InsertMysql) (err error) {
	comment := Trainer{
		Userid:       msg.Id,
		Content:      msg.Content,
		Start_time:   time.Now().Unix(),
		End_time:     time.Now().Unix() + msg.Expire,
		Read:         msg.Read,
		Message_type: msg.MessageType,
		FromUsername: msg.FromUsername,
		ToUsername:   msg.ToUsername,
	}
	err = dao.DB.Save(&comment).Error
	if err != nil {
		return err
	}
	return
}

// sql优化，不用in
// SELECT * FROM `trainers` where userid = 'oZ65W5TklL3gWTCLTllMfiXu97ig->20062111' UNION ALL SELECT * FROM `trainers` where userid = '20062111->oZ65W5TklL3gWTCLTllMfiXu97ig'  ORDER BY id desc LIMIT 0,5
func FindMany(sendID, id, host string, pageNum int) (results []Result, err error) {
	pageSize := 5
	pageSizeStr := strconv.Itoa(pageSize)
	pageNumStr := strconv.Itoa((pageNum - 1) * pageSize)
	var resultAll []Trainer //存放id和sendid的一些信息
	sql := "SELECT * FROM `trainers` where userid = ? UNION ALL SELECT * FROM `trainers` where userid = ?  ORDER BY id desc  LIMIT ?,?"
	//sql := "SELECT * FROM `trainers` where userid in ('" + id + "','" + sendID + "') ORDER BY id desc LIMIT " + pageNumStr + "," + pageSizeStr
	//fmt.Println(sql)
	dao.DB.Raw(sql, id, sendID, pageNumStr, pageSizeStr).Scan(&resultAll)
	for i, m := 0, len(resultAll); i < m; i++ {
		if resultAll[i].Message_type == 1 {
			resultAll[i].Content = formatUtil.GetPicHeaderBody(host, resultAll[i].Content)
		}
	}
	results, _ = AppendAndSort(resultAll, sendID, id)
	return
}

func AppendAndSort(resultAll []Trainer, sendID, id string) (results []Result, err error) {
	for _, r := range resultAll {
		start_time := time.Unix(r.Start_time, 0).Format("2006-01-02 15:04:05")
		sendSort := SendSortMsg{ //构造返回的msg
			Content:     r.Content,
			Read:        r.Read,
			CreatAt:     start_time,
			MessageType: r.Message_type,
		}
		var result Result
		if r.Userid == id {
			result = Result{ //构造返回所有的内容，包括传送者
				Start_time: r.Start_time,
				Msg:        sendSort,
				From:       "me",
			}
		} else {
			result = Result{ //构造返回所有的内容，包括传送者
				Start_time: r.Start_time,
				Msg:        sendSort,
				From:       "you",
			}
		}
		results = append(results, result)
	}
	return
}
