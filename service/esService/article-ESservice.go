package esService

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
)

type ArticleESservice struct {
}

func (aes *ArticleESservice) FuzzyArticlesQuery(pageNum int, fuzVal string) {
	offseter := sqlUtil.PageNumToSqlPage(pageNum, config.PageSize)
	//匹配查询
	matchPhraseQuery := elastic.NewMultiMatchQuery(fuzVal, "content", "title") //
	//matchPhraseQuery2 := elastic.NewMatchPhraseQuery("art_type", "job")
	//query := elastic.NewBoolQuery().
	//		Must(elastic.NewWildcardQuery("content", fuzVal))

	sortQuery1 := elastic.NewFieldSort("create_time").Desc()

	searchByPhrase, err := dao.ESClient.Search().Index(config.ArticleESIndex).
		SortBy(sortQuery1).
		From(offseter).Size(config.PageSize).
		//Query(matchPhraseQuery2).
		Query(matchPhraseQuery).
		Do(context.Background())
	for _, item := range searchByPhrase.Each(reflect.TypeOf(mysqlModel.Article{})) {
		language := item.(mysqlModel.Article)
		fmt.Println(language.CreateTime)
		fmt.Printf("search by phrase: %#v \n", language)
	}
	fmt.Println(err)
}
