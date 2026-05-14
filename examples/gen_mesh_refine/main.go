package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_mesh_refine —— 重布线/布线优化（node_type=10）
//
// 对模型网格进行布线重建优化。
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
	modelPath := assetsDir + "/high_model.fbx"

	fmt.Println("[gen_mesh_refine] 开始重布线...")

	modelID, err := client.GenMeshRefine(modelPath,
		visvise.NewGenMeshRefineOptions().
			SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").
			SetInputModelFormat("fbx").
			SetName("example_gen_mesh_refine"))
	if err != nil {
		log.Fatalf("[gen_mesh_refine] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_mesh_refine] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_mesh_refine] 等待失败: %v", err)
	}

	fmt.Printf("[gen_mesh_refine] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
