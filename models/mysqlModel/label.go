package mysqlModel

import (
	"encoding/json"
	"fmt"
	"sanHeRecruitment/dao"
)

type Label struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Level     int    `json:"level"`
	ParentId  int    `json:"-"`
	Label     string `json:"text"`
	Type      string `json:"type"`
	Recommend int    `json:"recommend"`
	SortNum   int    `json:"sort_num"`
}

type MaxLabelCount struct {
	MaxNum int `json:"max_num"`
}

type LabelOut struct {
	Label
	Value    int           `json:"value"`
	Children []interface{} `json:"children"`
}

type RecommendLabel struct {
	ID    int    `json:"id" gorm:"primary_key"`
	Label string `json:"text"`
	Type  string `json:"type"`
}

// AddLabel 将切片存入，测试函数
func AddLabel() []string {
	var lab Label
	a := []string{"北京", "唐山"}
	b, _ := json.Marshal(a)
	fmt.Println(b)
	c := string(b)
	fmt.Println(c)
	lab.Label = c
	err := dao.DB.Save(&lab).Error
	if err != nil {
		fmt.Println(err)
	}
	var x []string
	_ = json.Unmarshal([]byte(c), &x)
	fmt.Println("x:", x)
	return x
}
