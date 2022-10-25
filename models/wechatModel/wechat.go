package wechatModel

// WXLoginResp 微信登录申请url结构体
type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// RawData 微信获取rowData结构体
type RawData struct {
	NickName  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
}
