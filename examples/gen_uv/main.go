package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_uv —— UV 展开（node_type=9）
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
	modelPath := assetsDir + "/tex_model.obj"

	fmt.Println("[gen_uv] 开始 UV 展开...")

	modelID, err := client.GenUV(modelPath, rtx,
		visvise.NewGenUVOptions().
			SetAlgorithmModel("hunyuan3D-UV-v2.0").
			SetEnableAutoSmoothing(true).
			SetName("example_gen_uv"))
	if err != nil {
		log.Fatalf("[gen_uv] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_uv] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  600,
	})
	if err != nil {
		log.Fatalf("[gen_uv] 等待失败: %v", err)
	}

	fmt.Printf("[gen_uv] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
