package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_low_model —— 图生低模（node_type=13）
//
// 仅需主视图，其余视图可选。
//
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_UID=xxx VISVISE_ENV=prod go run main.go

func main() {
	appID := os.Getenv("VISVISE_APP_ID")
	secretKey := os.Getenv("VISVISE_SECRET_KEY")
	uid := os.Getenv("VISVISE_UID")
	envStr := os.Getenv("VISVISE_ENV")

	if appID == "" || secretKey == "" || uid == "" {
		log.Fatal("请设置环境变量: VISVISE_APP_ID, VISVISE_SECRET_KEY, VISVISE_UID")
	}

	env := visvise.EnvProd
	switch envStr {
	case "dev":
		env = visvise.EnvDev
	case "test":
		env = visvise.EnvTest
	}

	client := visvise.NewClient(appID, secretKey, uid,
		visvise.NewClientOptions().SetEnv(env))

	assetsDir := "./tests/assets"
	mainView := assetsDir + "/main_view.png"

	fmt.Println("[gen_low_model] 开始生成低模...")

	modelID, err := client.GenLowModel(mainView,
		visvise.NewGenLowModelOptions().
			SetAlgorithmModel("Tripo-v1.0-快速生成").
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetFaceType(visvise.FaceTypeTriangle).
			SetName("example_gen_low_model"))
	if err != nil {
		log.Fatalf("[gen_low_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_low_model] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, &visvise.WaitOptions{
		Interval: 3.0,
		Timeout:  600,
	})
	if err != nil {
		log.Fatalf("[gen_low_model] 等待失败: %v", err)
	}

	fmt.Printf("[gen_low_model] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
