package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_text_motion —— 文本生动画（node_type=4）
//
// 通过提示词描述动作自动生成 3D 动画，一次返回 4 个版本供抽卡选择。
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
	modelPath := assetsDir + "/animation_model.fbx"

	fmt.Println("[gen_text_motion] 开始文本生动画...")

	// 一次生成 4 个版本供抽卡
	modelIDs, err := client.GenTextMotion(modelPath, "一个人在跳街舞", rtx,
		visvise.NewGenTextMotionOptions().
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetName("example_gen_text_motion"))
	if err != nil {
		log.Fatalf("[gen_text_motion] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_text_motion] 任务已创建，共 %d 个版本：%v\n", len(modelIDs), modelIDs)

	fmt.Println("[gen_text_motion] 等待第一个版本完成（可按需等待全部）...")
	model, err := client.WaitModel(modelIDs[0], rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_text_motion] 等待失败: %v", err)
	}

	fmt.Printf("[gen_text_motion] model_ids[0] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
