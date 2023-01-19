#sanHeRecruitment
##项目结构
阅读时不要打开预览模式，以编辑模式阅读

-sanHeRecruitment
--config  //配置层
--controller //接口层
--dao //data的初始化层
--logs  //日志存放，路由日志已分割，错误日志输出到systemLogOut.txt,直接log.Println即可 
--middleware //中间件
--models // 各层所引用的结构体和未分割出来的websocket
--- BindModel //control层接受参数用来绑定
--- CreatorSql //创建数据库的sql
--- exportModel //excel到处的struct
--- loginModel //用来登录的
--- moduleModel //module层需要绑定的
--- mysqlModel //mysql每张表及service层查询数据库时需要的
--- websocketModel //websocket业务，最初没设计好导致封装在这里了
--- wechatModel //微信业务相关的结构体
--module //逻辑层，但是由于当初经验不足导致逻辑被control层执行大部分
--router //路由
--service //类似dao层或data层，与服务器交互，每个service代表每张表的交互
--ssl //https证书升级
--test //测试保存文件使用的，无实际意义
--timeTask //时间任务
--uploadPic //测试保存图片文件使用的，无实际意义
--util //工具层 但是忘记了很多次导致删除本地文件封装了多个函数
分别位于saveUtil，osUtil

--vendor //第三方依赖
--wechatPubAcc //微信公众号推送
--main.go //主函数

# 项目运行
需要先运行nsq
将models-CreatorSql中的sql文件运行创建数据库
再config-conf中进行配置
再主目录下运行go run main.go

# 开发
直接再control对应的逻辑下进行开发即可
control-module（简单逻辑可在control直接实现）-service
接口和逻辑  复杂逻辑                         数据库交互
注意websocketmodel 和 wsModule 不要相互引用结构体，会导致循环引用
目前及时通讯已调制到最佳状态，尽量不要修改主逻辑，可添加其他逻辑

main层中关闭了调试模式，将main主项中的文件注释即可开启调试模式

