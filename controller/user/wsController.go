package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sanHeRecruitment/config"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/module/websocketModule"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util"
	"sanHeRecruitment/util/tokenUtil"
	"sanHeRecruitment/util/uploadUtil"
	"strconv"
	"strings"
	"time"
)

type WsController struct {
	*service.UserService
	*service.ChatService
	*websocketModule.WsModule
	*service.MsgObjService
}

func WsControllerRouterToken(router *gin.RouterGroup) {
	w := WsController{}

	router.POST("/postPic", w.PostPic)
	router.POST("/msgPublisherLogout", w.MsgPublisherLogout)
	router.POST("/checkIfOneself", w.CheckIfOneself)
}

func WsControllerRouter(router *gin.RouterGroup) {
	w := WsController{}
	router.GET("/testCli", w.test)
	router.GET("/pushMsg", w.ReceiveMsgWebsocket)
	router.GET("/ws", w.Handler)
}

//var registerMux sync.Mutex

// Handler 聊一聊接口
func (ws *WsController) Handler(c *gin.Context) {
	//获取发送者uid和被发送者uid
	//auth := c.Request.Header.Get("Authorization")
	auth := c.Query("Authorization")
	if auth == "" {
		controller.ErrorResp(c, 201, "token not found")
		return
	}
	claims, _ := tokenUtil.ParseToken(auth)
	uid := claims.User.UserName
	//uid := c.Query("uid")
	toUid := c.Query("toUid")
	toUserInfo, err := ws.UserService.QueryUserInfoByUserId(toUid)
	if err != nil {
		controller.ErrorResp(c, 300, "该用户不存在")
		return
	}
	toUid = toUserInfo.Username
	if uid == toUid {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "自己无法与自己建立沟通",
		})
		return
	}
	//升级websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户实例
	client := &websocketModel.Client{
		ID:           ws.WsModule.CreateID(uid, toUid), //1->2
		SendID:       ws.WsModule.CreateID(toUid, uid), //2->1
		FromUsername: uid,
		ToUsername:   toUid,
		Socket:       conn,
		Send:         make(chan []byte),
	}
	//用户注册到用户管理上
	//registerMux.Lock()
	websocketModel.Manager.Register <- client
	//registerMux.Unlock()
	go client.Read(c.Request.Host)
	go client.Write(c.Request.Host)
	go func() {
		ws.ChatService.BatchRead(toUid, uid)
		ws.MsgObjService.BatchRead(uid, toUid)
	}()
}

// CheckIfOneself 立即沟通前check是不是与自己建立沟通
func (ws *WsController) CheckIfOneself(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	UserBasicInfo, _ := ws.UserService.QueryUserBasicInfo(username, c.Request.Host)
	userDataByte, _ := json.Marshal(&UserBasicInfo)
	userDataMap := make(map[string]interface{})
	_ = json.Unmarshal(userDataByte, &userDataMap)
	for _, v := range userDataMap {
		if v == "" || v == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "基础信息未完善",
			})
			return
		}
	}
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	toUid := recJson["toUid"].(string)
	toUserInfo, _ := ws.UserService.QueryUserInfoByUserId(toUid)
	toUser := toUserInfo.Username
	if username == "" || toUser == "" || username == toUser {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "不能与自己沟通",
		})
		return
	}
	uidOne := ws.WsModule.CreateID(username, toUserInfo.Username)
	uidTwo := ws.WsModule.CreateID(toUserInfo.Username, username)
	var firstFlag bool
	u1 := dao.Redis.Get(uidOne).Val()
	u2 := dao.Redis.Get(uidTwo).Val()
	if u1 == "" || u2 == "" {
		firstFlag = true
	} else {
		firstFlag = false
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "success",
		"firstFlag": firstFlag,
	})
}

// ReceiveMsgWebsocket websocket登陆后建立消息推送长连接(message publisher)
func (ws *WsController) ReceiveMsgWebsocket(c *gin.Context) {
	//auth := c.Request.Header.Get("Authorization")
	auth := c.Query("Authorization")
	if auth == "" {
		controller.ErrorResp(c, 201, "token not found")
		return
	}
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	//username = "3"
	//fmt.Println(username)
	//升级websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户实例
	recClient := &websocketModel.ClientRecMsg{
		ID:     username,
		Socket: conn,
		Send:   make(chan []byte),
	}
	ws.WsModule.AddMsgPusher(username, recClient)
	//models.ReceiveMsgManager.Clients[username] = recClient
	go recClient.PushMsg()
	go recClient.CheckOnline()
}

// MsgPublisherLogout 外部影响关闭websocket msgPublisher
func (ws *WsController) MsgPublisherLogout(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	var msgClient *websocketModel.ClientRecMsg
	if len(websocketModel.ReceiveMsgManager.Clients) != 0 {
		for _, client := range websocketModel.ReceiveMsgManager.Clients {
			if client.ID == username {
				msgClient = client
			}
		}
	}
	if msgClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "websocket msg publisher close failed,user doesn't exist",
		})
		return
	}
	//websocketModel.RecUnregisterMux.Lock()
	websocketModel.ReceiveMsgManager.Unregister <- msgClient
	//websocketModel.RecUnregisterMux.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "websocket msg publisher close success",
	})
	return
}

// PostPic websocket发送图片
func (ws *WsController) PostPic(c *gin.Context) {
	//获取发送者uid和被发送者uid
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	uid := claims.User.UserName
	//uid := c.Query("uid")
	toUid := c.Query("toUid")
	toUserInfo, err := ws.UserService.QueryUserInfoByUserId(toUid)
	if err != nil {
		controller.ErrorResp(c, 201, "该用户不存在")
		return
	}
	toUid = toUserInfo.Username
	//获取文件头
	file, err := c.FormFile("uploadPicture")
	if err != nil {
		//fmt.Println(file)
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	//fmt.Println(file.Filename)

	fileFormat := file.Filename[strings.Index(file.Filename, "."):]
	judgeFlag := uploadUtil.FormatJudge(fileFormat, ".jpg", ".png", ".jpeg", "gif")
	if judgeFlag == false {
		c.String(http.StatusBadRequest, "上传失败，格式不支持")
		return
	}

	uuid := util.GetUUID() + "-" + strconv.Itoa(int(time.Now().Unix()))
	newFileName := uuid + fileFormat
	filePath := config.PicSaverPath + "/" + newFileName
	realPicFormat := "uploadPic/" + newFileName

	userid := uid + "->" + toUid
	//fmt.Println(realPicFormat)

	var client *websocketModel.Client
	for id, conn := range websocketModel.Manager.Clients {
		if userid != id {
			continue
		}
		client = conn
	}
	if client != nil {
		dao.Redis.Incr(client.ID)
		//3个月进行一次过期，防止过快分手
		_, _ = dao.Redis.Expire(client.ID, time.Hour*24*30*3).Result()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 209,
			"msg":    "未连接到服务器，请返回重试",
		})
		//log.Println("PostPic", err)
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		log.Println("Msg pic save err,err:", err)
		return
	}
	websocketModel.Manager.Broadcast <- &websocketModel.Broadcast{
		Client:  client,
		Message: []byte(realPicFormat), //发送过来的消息
		Type:    1,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "发送成功",
	})
}

// websocket测试端口
func (ws *WsController) test(c *gin.Context) {
	//userid := "1->6"
	for id, conn := range websocketModel.Manager.Clients {
		fmt.Println(conn.ID)
		fmt.Println(id)
		if id != conn.ID {
			continue
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "发送成功",
	})
}
