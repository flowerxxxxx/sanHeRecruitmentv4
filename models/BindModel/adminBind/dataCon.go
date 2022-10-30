package adminBind

type PropagandaContentBinder struct {
	Url     string `json:"url"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type EditPropagandaContent struct {
	ID      int    `json:"id"`
	Url     string `json:"url"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type DeleteStreamBinder struct {
	Url string `json:"url"`
}

type DeleteProContBinder struct {
	ProId int `json:"pro_id"`
}

type NoticeSaveBinder struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

type EditNoticeBinder struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type DeleteNoticeBinder struct {
	Id int `json:"id"`
}

type AddPlatformCon struct {
	DesPerson string `json:"des_person"`
	Connect   string `json:"connect"`
	Type      string `json:"type"`
}

type EditConnectionBinder struct {
	Id        int    `json:"id"`
	DesPerson string `json:"des_person"`
	Connect   string `json:"connect"`
	Type      string `json:"type"`
}

type DeleteConBinder struct {
	Id int `json:"id"`
}

type PlatDescriptionBinder struct {
	Content string `json:"content"`
	Module  string `json:"module"`
}

type EditPlatDesBinder struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Module  string `json:"module"`
}

type DelPlatDesBinder struct {
	Id int `json:"id"`
}

type AddRecoLabelBinder struct {
	Id      int `json:"id"`
	DesReco int `json:"desReco"`
}

type VipStyleAddBinder struct {
	CoverUrl string `json:"cover_url"`
	Content  string `json:"content"`
	Title    string `json:"title"`
}

type VipStyleEditBinder struct {
	Id       int    `json:"id"`
	CoverUrl string `json:"cover_url"`
	Content  string `json:"content"`
	Title    string `json:"title"`
}

type VipStyleDelBinder struct {
	Id int `json:"id"`
}
