package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_high_model —— 图生高模（node_type=3）
//
// 高模用第二个算法模型 Tripo-v3.1-ultra（支持单图输入，无需四视图）。
// 也可传入四视图以提升质量，优先从 MV_360_MODEL_ID 获取。
//
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod go run main.go

func main() {
	appID := os.Getenv("VISVISE_APP_ID")
	secretKey := os.Getenv("VISVISE_SECRET_KEY")
	rtx := os.Getenv("VISVISE_RTX")
	envStr := os.Getenv("VISVISE_ENV")

	if appID == "" || secretKey == "" || rtx == "" {
		log.Fatal("请设置环境变量: VISVISE_APP_ID, VISVISE_SECRET_KEY, VISVISE_RTX")
	}

	env := visvise.EnvProd
	switch envStr {
	case "dev":
		env = visvise.EnvDev
	case "test":
		env = visvise.EnvTest
	}

	client := visvise.NewClient(appID, secretKey,
		visvise.NewClientOptions().SetEnv(env))

	assetsDir := "./tests/assets"
	mvModelID := os.Getenv("MV_360_MODEL_ID")

	var mainView string

	if mvModelID != "" {
		fmt.Printf("[gen_high_model] 从 gen_360 输出提取四视图 (model_id=%s)\n", mvModelID)
		models, _, err := client.GetAPI().GetModelList([]string{mvModelID}, nil, nil, "", 10, 1, rtx)
		if err != nil || len(models) == 0 {
			log.Fatalf("[gen_high_model] 获取模型失败: %v", err)
		}
		out := models[0].ImageGen360Output.OutputView
		mainView = stripSign(out.MainView)
	} else {
		// 使用本地主视图（Tripo-v3.1-ultra 支持单图）
		mainView = assetsDir + "/main_view.png"
	}

	opts := visvise.NewGenHighModelOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle).
		SetName("example_gen_high_model")

	modelID, err := client.GenHighModel(mainView, rtx, opts)
	if err != nil {
		log.Fatalf("[gen_high_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_high_model] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 10.0,
		Timeout:  1200,
	})
	if err != nil {
		log.Fatalf("[gen_high_model] 等待失败: %v", err)
	}

	fmt.Printf("[gen_high_model] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}

func stripSign(url string) string {
	for i := 0; i < len(url); i++ {
		if url[i] == '?' {
			return url[:i]
		}
	}
	return url
}
