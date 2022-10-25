package loginModel

type LoginInfo struct {
	AvatarUrl string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	Nickname  string `json:"nickName"`
	Province  string `json:"province"`
	UserLevel int    `json:"user_level"`
}
