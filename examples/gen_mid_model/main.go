package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_mid_model —— 图生中模（node_type=3）
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod go run main.go

func main() {
	appID := os.Getenv("VISVISE_APP_ID")
	secretKey := os.Getenv("VISVISE_SECRET_KEY")
	rtx := os.Getenv("VISVISE_RTX")
	envStr := os.Getenv("VISVISE_ENV")
	mvModelID := os.Getenv("MV_360_MODEL_ID")

	if appID == "" || secretKey == "" || rtx == "" || mvModelID == "" {
		log.Fatal("请设置环境变量: VISVISE_APP_ID, VISVISE_SECRET_KEY, VISVISE_RTX, MV_360_MODEL_ID")
	}

	env := visvise.EnvProd
	switch envStr {
	case "dev":
		env = visvise.EnvDev
	case "test":
		env = visvise.EnvTest
	}

	client := visvise.NewClient(appID, secretKey, visvise.NewClientOptions().SetEnv(env))

	// modelID := example1(client, rtx)
	// modelID := example2(client, mvModelID, rtx)
	modelID := example3(client, mvModelID, rtx)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_mid_model] 等待失败: %v", err)
	}

	fmt.Printf("[gen_mid_model] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}

// =================  使用场景一：用户上传原画视图  =====================
func example1(client *visvise.Client, rtx string) string {
	assetsDir := "./tests/assets"
	mainView := assetsDir + "/main_view.png"
	fmt.Println("[gen_mid_model] 开始生成中模...")
	modelID, err := client.GenMidModel(stripSign(mainView), nil, nil, nil, rtx,
		visvise.NewGenMidModelOptions().
			SetName("example_gen_mid_model"))
	if err != nil {
		log.Fatalf("[gen_mid_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_mid_model] 任务已创建，model_id=%s\n", modelID)
	return modelID
}

// =================  使用场景二：基于【图生360】的生成结果生成模型资产  =====================
func example2(client *visvise.Client, mvModelID, rtx string) string {
	if mvModelID == "" {
		models, _, err := client.GetAPI().GetModelList(nil, []int{7}, []int{3}, "", 10, 1, rtx)
		if err != nil || len(models) == 0 {
			log.Fatalf("[gen_mid_model] 获取模型失败: %v", err)
		}
		mvModelID = models[0].ModelID
	}

	fmt.Printf("[gen_mid_model] gen_360 (model_id=%s)\n", mvModelID)

	fmt.Println("[gen_mid_model] 开始生成中模...")
	modelID, err := client.GenMidModel(nil, nil, nil, nil, rtx,
		visvise.NewGenMidModelOptions().
			SetName("example_gen_mid_model").
			SetModelID360(mvModelID))
	if err != nil {
		log.Fatalf("[gen_mid_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_mid_model] 任务已创建，model_id=%s\n", modelID)
	return modelID
}

// =================  使用场景三：基于【2D拆分】的生成结果生成模型资产  =====================
func example3(client *visvise.Client, mvModelID, rtx string) string {
	if mvModelID == "" {
		models, _, err := client.GetAPI().GetModelList(nil, []int{14}, []int{3}, "", 10, 1, rtx)
		if err != nil || len(models) == 0 {
			log.Fatalf("[gen_mid_model] 获取模型失败: %v", err)
		}
		mvModelID = models[0].ModelID
	}

	fmt.Printf("[gen_mid_model] 2D拆分 (model_id=%s)\n", mvModelID)

	fmt.Println("[gen_mid_model] 开始生成中模...")
	modelID, err := client.GenMidModel(nil, nil, nil, nil, rtx,
		visvise.NewGenMidModelOptions().
			SetName("example_gen_mid_model").
			SetSegmentModelID(mvModelID))
	if err != nil {
		log.Fatalf("[gen_mid_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_mid_model] 任务已创建，model_id=%s\n", modelID)
	return modelID
}

func stripSign(url string) string {
	for i := 0; i < len(url); i++ {
		if url[i] == '?' {
			return url[:i]
		}
	}
	return url
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
