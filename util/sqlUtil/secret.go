package sqlUtil

import (
	"crypto/md5"
	"fmt"
	"sanHeRecruitment/config"
)

func GenMD5Password(str string) string {
	data := []byte(str + config.ProducerUsername) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
