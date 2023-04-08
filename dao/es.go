package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"sanHeRecruitment/config"
	"strings"
)

var ESClient *elastic.Client
var eSServerURL = []string{config.ESServerURL}

func init() {
	var err error
	ESClient, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(eSServerURL...))
	if err != nil {
		log.Printf("Failed to build elasticsearch connection: %s %s", strings.Join(eSServerURL, ","), err.Error())
		panic(any("Failed to build elasticsearch connection"))
	}
	info, code, err := ESClient.Ping(strings.Join(eSServerURL, ",")).Do(context.Background())
	if err != nil {
		log.Printf("ping es failed,%s", err.Error())
		panic(any("es ping es failed"))
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}
