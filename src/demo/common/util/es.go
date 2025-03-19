package util

import (
	"sls-mall-go/common/config"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func InitES() {
	addresses := strings.Split(config.EsAddr, ",")
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  config.EsUser,
		Password:  config.EsPass,
	}
	//fmt.Println(cfg)
	transport := NewOTELHTTPTransport(cfg.Transport)
	cfg.Transport = transport

	es, err := elasticsearch.NewClient(cfg)
	Chk(err)
	ESClient = es
}
