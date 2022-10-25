package timeUtil

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func CreateCaptcha() int64 {
	randStr := fmt.Sprintf("%03v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	randNum, _ := strconv.ParseInt(randStr, 10, 64)
	return randNum
}
