package main

import (
	"fmt"
	"github.com/visvise/visvise-sdk-go/visvise"
	"log"
	"os"
)

// Example: gen_360 —— 图生360（生成多视图）
//
// 从单张图片生成 360° 多视图，输出结果可作为图生高模/中模/低模的输入。
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

	// 素材路径（请替换为实际路径）
	assetsDir := "./tests/assets"
	mainView := assetsDir + "/main_view.png"

	opts := visvise.NewGen360Options().
		SetAlgorithmModel("VISVISE-MultiView-V1.0.0").
		SetName("example_gen_360")

	fmt.Printf("[gen_360] 开始生成多视图，输入图片：%s\n", mainView)

	modelID, err := client.Gen360(mainView, rtx, opts)
	if err != nil {
		log.Fatalf("[gen_360] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_360] 任务已创建，model_id=%s\n", modelID)

	fmt.Println("[gen_360] 等待完成...")
	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 3.0,
		Timeout:  300,
	})
	if err != nil {
		log.Fatalf("[gen_360] 等待失败: %v", err)
	}

	fmt.Printf("[gen_360] 生成成功！耗时 %ds\n", model.TimeCost)
	if model.ImageGen360Output != nil && model.ImageGen360Output.OutputView != nil {
		v := model.ImageGen360Output.OutputView
		fmt.Printf("  main_view  : %s\n", v.MainView)
		fmt.Printf("  back_view  : %s\n", v.BackView)
		fmt.Printf("  left_view  : %s\n", v.LeftView)
		fmt.Printf("  right_view : %s\n", v.RightView)
	}
	if model.ImageGen360Output != nil && model.ImageGen360Output.HorizontalViewVideo != "" {
		fmt.Printf("  旋转视频   : %s\n", model.ImageGen360Output.HorizontalViewVideo)
	}
}
