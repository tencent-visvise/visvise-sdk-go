package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_video_motion —— 视频生动画（node_type=4）
//
// 从视频中提取动作数据，驱动 3D 模型生成动画。
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
	videoPath := assetsDir + "/animation_video.mp4"

	fmt.Println("[gen_video_motion] 开始视频生动画...")

	modelID, err := client.GenVideoMotion(modelPath, videoPath, rtx,
		visvise.NewGenVideoMotionOptions().
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetWithHand(true).
			SetName("example_gen_video_motion"))
	if err != nil {
		log.Fatalf("[gen_video_motion] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_video_motion] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_video_motion] 等待失败: %v", err)
	}

	fmt.Printf("[gen_video_motion] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
