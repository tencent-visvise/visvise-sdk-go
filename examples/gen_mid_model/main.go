package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_mid_model —— 图生中模（node_type=3）
//
// 中模要求四视图全部必传。
// 优先从环境变量 MV_360_MODEL_ID 读取 gen_360 的输出，自动提取四视图。
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

	var mainView, backView, leftView, rightView string

	if mvModelID != "" {
		fmt.Printf("[gen_mid_model] 从 gen_360 输出提取四视图 (model_id=%s)\n", mvModelID)
		models, _, err := client.GetAPI().GetModelList([]string{mvModelID}, nil, nil, "", 10, 1, rtx)
		if err != nil || len(models) == 0 {
			log.Fatalf("[gen_mid_model] 获取模型失败: %v", err)
		}
		out := models[0].ImageGen360Output.OutputView
		mainView = stripSign(out.MainView)
		backView = stripSign(out.BackView)
		leftView = stripSign(out.LeftView)
		rightView = stripSign(out.RightView)
	} else {
		mainView = getEnvOrDefault("MV_MAIN", assetsDir+"/main_view.png")
		backView = getEnvOrDefault("MV_BACK", assetsDir+"/back_view.png")
		leftView = getEnvOrDefault("MV_LEFT", assetsDir+"/left_view.png")
		rightView = getEnvOrDefault("MV_RIGHT", assetsDir+"/right_view.png")
	}

	fmt.Println("[gen_mid_model] 开始生成中模...")

	modelID, err := client.GenMidModel(mainView, backView, leftView, rightView, rtx,
		visvise.NewGenMidModelOptions().
			SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetFaceType(visvise.FaceTypeTriangle).
			SetName("example_gen_mid_model"))
	if err != nil {
		log.Fatalf("[gen_mid_model] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_mid_model] 任务已创建，model_id=%s\n", modelID)

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
