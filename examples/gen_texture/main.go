package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_texture —— 贴图纹理生成（node_type=8）
//
// input_view 中的本地图片路径会被 SDK 自动上传到 COS。
// input_view.main_view 和 prompt 必须传其中一个。
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

	fmt.Println("[gen_texture] 开始贴图纹理生成...")

	// 本地图片路径，SDK 内部自动上传到 COS
	inputView := &visvise.View{
		MainView: assetsDir + "/tex_ref_front.jpg",
		BackView: assetsDir + "/tex_ref_back.jpg",
		LeftView: assetsDir + "/tex_ref_left.jpg",
	}

	modelID, err := client.GenTexture(assetsDir+"/tex_model.obj", rtx,
		visvise.NewGenTextureOptions().
			SetAlgorithmModel("hunyuan3D-TEX-v2.0").
			SetInputView(inputView).
			SetResolution(2048).
			SetUnwarpUV(false).
			SetName("example_gen_texture"))
	if err != nil {
		log.Fatalf("[gen_texture] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_texture] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_texture] 等待失败: %v", err)
	}

	fmt.Printf("[gen_texture] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
