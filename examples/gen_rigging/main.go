package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_rigging —— 骨骼架设（node_type=5）
//
// 为输入的 3D 模型自动生成骨骼结构。
// SDK 内部自动构建 JSON 参数文件并打包成 zip 上传。
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
	modelPath := assetsDir + "/rigging_model.fbx"

	fmt.Println("[gen_rigging] 开始骨骼架设...")

	// SDK 自动将 model.fbx + 参数 JSON 打包成 zip 上传
	// 无需手动准备 zip 包
	modelID, err := client.GenRigging(modelPath,
		visvise.NewGenRiggingOptions().
			SetAlgorithmModel("VISVISE-GoRigging-V1.0.0").
			SetMeshCategory(visvise.MeshCategoryHumanoid). // 人形（默认）或 visvise.MeshCategoryTetrapod（四足）
			SetName("example_gen_rigging"))
	if err != nil {
		log.Fatalf("[gen_rigging] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_rigging] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  600,
	})
	if err != nil {
		log.Fatalf("[gen_rigging] 等待失败: %v", err)
	}

	fmt.Printf("[gen_rigging] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
