package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_retopology —— 重拓扑（node_type=1）
//
// 对高面数模型进行拓扑优化。
// 混元模型传 detail_level，VISVISE 自研模型传 face_num，二选一。
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
	modelPath := assetsDir + "/high_model.fbx"

	fmt.Println("[gen_retopology] 开始重拓扑...")

	opts := visvise.NewGenRetopologyOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad).
		SetDetailLevel(visvise.DetailLevelHigh). // 混元模型用 detail_level
		SetName("example_gen_retopology")

	modelID, err := client.GenRetopology(modelPath, rtx, opts)
	if err != nil {
		log.Fatalf("[gen_retopology] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_retopology] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_retopology] 等待失败: %v", err)
	}

	fmt.Printf("[gen_retopology] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}

// 另一种使用 VISVISE 自研模型的写法（传 face_num）：
// opts := visvise.NewGenRetopologyOptions().
//     SetAlgorithmModel("VISVISE-RTP-V1.0.0").
//     SetOutputModelFormat(visvise.OutputModelFormatFBX).
//     SetFaceType(visvise.FaceTypeQuad).
//     SetFaceNum(10000). // VISVISE 自研模型用 face_num
//     SetName("example_gen_retopology")
