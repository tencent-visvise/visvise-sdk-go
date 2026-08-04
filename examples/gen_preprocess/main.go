package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_style_transfer / gen_patter_auto_remove —— 2D 预处理。
//
// 输入本地图片或 VISVISE 平台 COS URL；本地图片会自动上传。
// 未设置 VISVISE_PREPROCESS_INPUT 时，使用 tests/assets/preprocess.png 示例图片。
// 每次仅执行一个预处理工作流并同步返回保存资产的 model_id。
// Usage: VISVISE_APP_ID=... VISVISE_SECRET_KEY=... VISVISE_RTX=... go run main.go
func main() {
	appID := os.Getenv("VISVISE_APP_ID")
	secretKey := os.Getenv("VISVISE_SECRET_KEY")
	rtx := os.Getenv("VISVISE_RTX")
	if appID == "" || secretKey == "" || rtx == "" {
		log.Fatal("请设置环境变量: VISVISE_APP_ID, VISVISE_SECRET_KEY, VISVISE_RTX")
	}
	assetsDir := "./tests/assets"
	inputView := envOrDefault("VISVISE_PREPROCESS_INPUT", assetsDir+"/preprocess.png")

	env := visvise.EnvProd
	switch os.Getenv("VISVISE_ENV") {
	case "dev":
		env = visvise.EnvDev
	case "test":
		env = visvise.EnvTest
	}
	client := visvise.NewClient(appID, secretKey, visvise.NewClientOptions().SetEnv(env).SetTimeout(180))

	name := envOrDefault("VISVISE_PREPROCESS_NAME", "example_gen_preprocess")
	algorithmModel := os.Getenv("VISVISE_PREPROCESS_ALGORITHM_MODEL")

	var modelID string
	var err error
	switch mode := envOrDefault("VISVISE_PREPROCESS_MODE", "stylized"); mode {
	case "stylized":
		opts := visvise.NewGenStyleTransferOptions().
			SetName(name)
		if algorithmModel != "" {
			opts.SetAlgorithmModel(algorithmModel)
		}
		modelID, err = client.GenStyleTransfer(inputView, styleType(), rtx, opts)
	case "auto-remove":
		opts := visvise.NewGenPatterAutoRemoveOptions().SetName(name)
		if algorithmModel != "" {
			opts.SetAlgorithmModel(algorithmModel)
		}
		modelID, err = client.GenPatterAutoRemove(inputView, rtx, opts)
	default:
		log.Fatalf("不支持的 VISVISE_PREPROCESS_MODE: %s", mode)
	}
	if err != nil {
		log.Fatalf("[gen_preprocess] 创建资产失败: %v", err)
	}
	fmt.Printf("[gen_preprocess] 预处理资产已创建，model_id=%s\n", modelID)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func styleType() visvise.StyleType {
	switch envOrDefault("VISVISE_PREPROCESS_STYLE", "grayscale") {
	case "grayscale":
		return visvise.StyleTypeGrayscale
	case "pixel":
		return visvise.StyleTypePixel
	case "realistic":
		return visvise.StyleTypeRealistic
	case "cartoon":
		return visvise.StyleTypeCartoon
	default:
		log.Fatal("VISVISE_PREPROCESS_STYLE 仅支持 grayscale、pixel、realistic 或 cartoon")
		return 0
	}
}
