package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sanHeRecruitment/biz/nsqBiz"
	"sanHeRecruitment/biz/websocketBiz"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/remote"
	"sanHeRecruitment/router"
	"sanHeRecruitment/timeTask"
	"sanHeRecruitment/util/logUtil"
	"syscall"
	"time"
)

// var MysqlModels = []interface{}{&mysqlModel.User{}}
var ws *websocketBiz.WsModule

func main() {
	//main 主项
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	go remote.StartBeToServer("0.0.0.0" + config.RemoteServer)

	logUtil.LogOutInit()
	//nsq开启生产者
	_ = nsqBiz.InitProducer()
	//nsq开启消费者
	nsqBiz.Consumer()

	//将http写入config
	//logfile, err := os.Create("./gin_http.log")
	//if err != nil {
	//	fmt.Println("Could not create log file")
	//}
	//gin.SetMode(gin.DebugMode)
	//gin.DefaultWriter = io.MultiWriter(logfile)

	//开启定时任务
	timeTask.InitTimer()
	//连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(any(err))
	}
	defer dao.Close() //程序退出关闭数据库
	//初始化redis000000
	err = dao.InitRedis()
	if err != nil {
		panic(any(err))
	}
	//开启监听 双线程监听
	for i := 0; i < config.GoroutineNum; i++ {
		//websocket监听及处理线程
		go ws.WsStart()
		//开启接收消费者动作的处理
		go nsqBiz.ReceiveToInsert()
		go websocketBiz.FromMainToPush()
		//开启消息推送的管理线程
		go ws.RecMsgStart()
	}
	//模型绑定
	//dao.DB.AutoMigrate(MysqlModels...)
	r := router.SetupRouter()
	err = websocketBiz.InitSystemAdminer()
	if err != nil {
		panic(any(err))
	}
	//log.Println("sysAdminer init success")
	//开启的端口号
	//err = r.RunTLS(config.TLSConfig.Addr, config.TLSConfig.CertFile, config.TLSConfig.KeyFile)
	//if err != nil {
	//	panic(err)
	//}

	srv := &http.Server{
		Addr:    config.TLSConfig.Addr,
		Handler: r,
	}

	go func() {
		// 基于https的服务连接开启
		if err := srv.ListenAndServeTLS(config.TLSConfig.CertFile, config.TLSConfig.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 3 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//log.Println(time.Now().Format("2006-01-02 15:04:05")+"Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 3 seconds.
	select {
	case <-ctx.Done():
		//log.Println("timeout of 3 seconds.")
	}
	//log.Println(time.Now().Format("2006-01-02 15:04:05")+" Server Done")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " Server Done")
}
