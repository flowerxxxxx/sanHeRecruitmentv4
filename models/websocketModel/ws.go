package websocketModel

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/util/e"
	"sanHeRecruitment/util/formatUtil"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	Start_time int64       `json:"start_time"`
	Msg        SendSortMsg `json:"msg"`
	From       string      `json:"from"`
	Code       int         `json:"code"`
}

type SendSortMsg struct {
	Content     string `json:"content"`
	Read        int    `json:"read"`
	CreatAt     string `json:"creat_at"`
	MessageType int    `json:"message_type"`
}

// 发送消息
type SendMsg struct {
	Type        int    `json:"type"`
	Content     string `json:"content"`
	MessageType int    `json:"message_type"`
}

// 广播写入conn.send结构体
type BroadcastMsg struct {
	Message     string `json:"message"`
	MessageType int    `json:"message_type"`
}

// 回复消息
type ReplyMsg struct {
	From        string `json:"from"`
	Code        int    `json:"code"`
	Content     string `json:"content"`
	MessageType int    `json:"message_type"`
	PicShow     bool   `json:"pic_show"`
}

// 回复消息2
type ReplyMsg2 struct {
	From string `json:"from"`
	//Msg         interface{} `json:"msg"`
	Code        int    `json:"code"`
	Content     string `json:"content"`
	Read        int    `json:"read"`
	CreatAt     string `json:"creat_at"`
	MessageType int    `json:"message_type"`
	PicShow     bool   `json:"pic_show"`
}

// 用户类
type Client struct {
	ID           string
	SendID       string
	FromUsername string
	ToUsername   string
	Socket       *websocket.Conn
	Send         chan []byte
	SocketMutex  sync.Mutex
}

// 用户登陆后的长连接结构体
type ClientRecMsg struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

// 管理消息推送长连接
type ClientRecMsgManager struct {
	Clients map[string]*ClientRecMsg
	//用户计数器，用来缓存websocket延迟关闭删除用户导致消息推送失败
	ClientCount map[string]int
	Broadcast   chan *Broadcast
	Unregister  chan *ClientRecMsg
	ClientsRWM  sync.RWMutex
	CliCountRWM sync.RWMutex
}

// 广播类
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 管理用户登录登出回复广告等
type ClientManager struct {
	Clients     map[string]*Client
	ClientCount map[string]int
	ClientsRWM  sync.RWMutex
	CliCountRWM sync.RWMutex
	Broadcast   chan *Broadcast
	Reply       chan *Client
	Register    chan *Client
	Unregister  chan *Client
}

// 序列化信息
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// 管理消息推送
var ReceiveMsgManager = ClientRecMsgManager{
	Clients:     make(map[string]*ClientRecMsg),
	ClientCount: make(map[string]int),
	Broadcast:   make(chan *Broadcast, 100),
	Unregister:  make(chan *ClientRecMsg),
}

// 消息推送具体内容
type PublishMsg struct {
	FromUser       string `json:"from_user"`
	MessageContent string `json:"message_content"`
	MessageType    int    `json:"message_type"`
	HeartBeat      int    `json:"HeartBeat"`
}

// 管理
var Manager = ClientManager{
	Clients:     make(map[string]*Client),
	ClientCount: make(map[string]int),
	Broadcast:   make(chan *Broadcast, 100),
	Register:    make(chan *Client),
	Reply:       make(chan *Client),
	Unregister:  make(chan *Client),
}

// websocket用户写入数据
func (c *Client) Read(host string) {
	defer func() {
		close(c.Send)
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		sendMSg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMSg)
		if err != nil {
			//TODO ws conn close reason printer
			//fmt.Println("数据格式不正确", err)
			//Manager.Unregister <- c
			//_ = c.Socket.Close()
			break
		}
		if sendMSg.Type == 1 {
			r1 := dao.Redis.Get(c.ID).Val()     // 1->2
			r2 := dao.Redis.Get(c.SendID).Val() // 2->1
			//fmt.Println("r1", r1, "r2")
			//fmt.Println(r1 == "", r2 == "")
			r1Int := 0
			r1Int, errInt := strconv.Atoi(r1)
			if errInt != nil {
				r1Int = 0
			}
			if r1Int > config.FirstUnreadMsgNum && r2 == "" {
				//1给2发消息，发了大于FirstUnreadMsgNumStr条，但是2没有回，或者没有看到，就停止1的发送
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "消息数达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				c.SocketMutex.Lock()
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				c.SocketMutex.Unlock()
				continue
			} else {
				dao.Redis.Incr(c.ID)
				//3个月进行一次过期，防止过快分手
				_, _ = dao.Redis.Expire(c.ID, time.Hour*24*30*3).Result()
			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMSg.Content), //发送过来的消息
				Type:    0,
			}
		} else if sendMSg.Type == 2 {
			//这里将content内获取的数字当作页码
			pageNum, err := strconv.Atoi(sendMSg.Content)
			if err != nil {
				//如果获取错误默认收到的额content的值是1
				pageNum = 1
			}
			results, _ := FindMany(c.SendID, c.ID, host, pageNum)
			if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				//RWMux.Lock()
				c.SocketMutex.Lock()
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				c.SocketMutex.Unlock()
				//RWMux.Unlock()
				continue
			}
			if pageNum == 1 {
				results = ReverseResults(results)
			}
			for _, result := range results {
				//history msg
				replyMsg := ReplyMsg2{
					From:        result.From,
					Code:        e.WebsocketHistoryMsg,
					Content:     result.Msg.Content,
					Read:        result.Msg.Read,
					CreatAt:     result.Msg.CreatAt,
					MessageType: result.Msg.MessageType,
					//Msg:  result.Msg,
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				//RWMux.Lock()
				c.SocketMutex.Lock()
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				c.SocketMutex.Unlock()
				//RWMux.Unlock()
			}
		}
	}
}

// 将结果倒序
func ReverseResults(s []Result) []Result {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// websocket向用户写入数据
func (c *Client) Write(host string) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.SocketMutex.Lock()
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				c.SocketMutex.Unlock()
				return
			}
			var message2 BroadcastMsg
			_ = json.Unmarshal(message, &message2)
			if message2.MessageType == 1 {
				message2.Message = formatUtil.GetPicHeaderBody(host, message2.Message)
			}
			replyMsg := &ReplyMsg{
				From:        "you",
				Code:        e.WebsocketSuccessMessage,
				Content:     message2.Message,
				MessageType: message2.MessageType,
			}

			msg, _ := json.Marshal(replyMsg)
			//RWMux.Lock()
			c.SocketMutex.Lock()
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			c.SocketMutex.Unlock()
			//RWMux.Unlock()
		}
	}
}

// PushMsg websocket用户消息推送机制
func (r *ClientRecMsg) PushMsg() {
	defer func() {
		//fmt.Println("close succ")
		//_ = r.Socket.Close()
	}()
	for {
		r.Socket.PongHandler()
		select {
		case msg, ok := <-r.Send:
			if !ok {
				_ = r.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				// TODO msgPusher closer printer
				//fmt.Println("msgPush close succ")
				return
			}
			//RWMux.Lock()
			_ = r.Socket.WriteMessage(websocket.TextMessage, msg)
			//RWMux.Unlock()
		}
	}
}

func (c *ClientRecMsg) CheckOnline() {
	defer func() {
		if flag := deleteMsgPusher(c); flag {
			_ = c.Socket.Close()
			//close(c.Send)
			ReceiveMsgManager.Unregister <- c
		} else {
			_ = c.Socket.Close()
			close(c.Send)
		}
	}()

	for {
		PushMsg := struct {
			HeartBeat int
		}{1}
		msg, _ := json.Marshal(PushMsg)
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			//fmt.Println("[CheckLogger]check websocket_biz close")
			break
		}
		//fmt.Println("[HeartbeatLogger]websocket_biz heartbeat")
		time.Sleep(10 * time.Second)
	}
}

// msgPusherCleaner
func deleteMsgPusher(cliRec *ClientRecMsg) (flag bool) {
	//msgCount := ReceiveMsgManager.ClientCount[cliRec.ID]
	msgCount := ReadRecManCliCount(cliRec.ID)
	if msgCount > 1 {
		msgCount--
		ReceiveMsgManager.CliCountRWM.Lock()
		ReceiveMsgManager.ClientCount[cliRec.ID] = msgCount
		ReceiveMsgManager.CliCountRWM.Unlock()
		//log.Println("[CutLogger]",cliRec.ID,"msgCount cut one")
		return false
	} else {
		ReceiveMsgManager.CliCountRWM.Lock()
		delete(ReceiveMsgManager.ClientCount, cliRec.ID)
		ReceiveMsgManager.CliCountRWM.Unlock()
		return true
	}
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

// ReadRecManCliCount read the RecManCliCount where is sync safe
func ReadRecManCliCount(cliRecId string) (msgCount int) {
	ReceiveMsgManager.CliCountRWM.RLock()
	msgCount = ReceiveMsgManager.ClientCount[cliRecId]
	ReceiveMsgManager.CliCountRWM.RUnlock()
	return
}

// EditRecManCliCount edit the RecManCliCount where is sync safe
func EditRecManCliCount(cliRecId string, msgCount int) {
	ReceiveMsgManager.CliCountRWM.Lock()
	ReceiveMsgManager.ClientCount[cliRecId] = msgCount
	ReceiveMsgManager.CliCountRWM.Unlock()
}

func ReadRecManClient(cliId string) (cliMap *ClientRecMsg, ok bool) {
	ReceiveMsgManager.ClientsRWM.RLock()
	cliMap, ok = ReceiveMsgManager.Clients[cliId]
	ReceiveMsgManager.ClientsRWM.RUnlock()
	return
}

func ReadTotalRecManClients() int {
	ReceiveMsgManager.ClientsRWM.RLock()
	totalCounts := len(ReceiveMsgManager.Clients)
	ReceiveMsgManager.ClientsRWM.RUnlock()
	return totalCounts
}

func ReadManClient(cliId string) (cliMap *Client, ok bool) {
	Manager.ClientsRWM.RLock()
	cliMap, ok = Manager.Clients[cliId]
	Manager.ClientsRWM.RUnlock()
	return
}

// DelRecManCli del receiver cli
func DelRecManCli(cliId string) {
	ReceiveMsgManager.ClientsRWM.Lock()
	delete(ReceiveMsgManager.Clients, cliId)
	ReceiveMsgManager.ClientsRWM.Unlock()
}

// DelRecCliCountCli del receiver count cli
func DelRecCliCountCli(cliId string) {
	ReceiveMsgManager.CliCountRWM.Lock()
	delete(ReceiveMsgManager.ClientCount, cliId)
	ReceiveMsgManager.CliCountRWM.Unlock()
}

func ManagerCliCountIncr(connID string) {
	Manager.CliCountRWM.Lock()
	Manager.ClientCount[connID] = Manager.ClientCount[connID] + 1
	Manager.CliCountRWM.Unlock()
}

func ManagerCliCountCutOne(connID string) {
	Manager.CliCountRWM.Lock()
	Manager.ClientCount[connID] = Manager.ClientCount[connID] - 1
	Manager.CliCountRWM.Unlock()
}

func ReadCliCount(connID string) (cliCount int, ok bool) {
	Manager.CliCountRWM.RLock()
	cliCount, ok = Manager.ClientCount[connID]
	Manager.CliCountRWM.RUnlock()
	return
}

// DelManagerCliCount del manager count cli
func DelManagerCliCount(connID string) {
	Manager.CliCountRWM.Lock()
	delete(Manager.ClientCount, connID)
	Manager.CliCountRWM.Unlock()
}
