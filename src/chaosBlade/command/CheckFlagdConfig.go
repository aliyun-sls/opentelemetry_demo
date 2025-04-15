package command

import (
	"context"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
)

func initFeatureFlag() {
	once.Do(func() {
		provider := flagd.NewProvider(
			flagd.WithHost("flagd"),
			flagd.WithPort(8013),
		)
		if err := openfeature.SetProvider(provider); err != nil {
			log.Printf("设置provider失败: %v", err)
			return
		}
		flagClient = openfeature.NewClient("product")
	})
}

func CheckRDSFlagdConfig() string {
	// 初始化flag客户端(只执行一次)
	initFeatureFlag()

	// 获取feature flag值
	istrue := flagClient.String(
		context.Background(),
		"RDSLoss",
		"off",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取feature : %v", istrue)
	if istrue == "on" {
		return "app=product"
	} else {
		return ""
	}

}
