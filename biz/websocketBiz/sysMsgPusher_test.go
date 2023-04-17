package websocketBiz

import (
	"testing"
	"time"
)

func TestSysMsgPusher(t *testing.T) {
	SysMsgPusher("oZ65W5UKcACPZCW3AqsFf8LEs7tM", "<h1>系统通知：您的升级已经审核成功</h1>")
	time.Sleep(time.Second)
}
