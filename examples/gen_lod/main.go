package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_lod —— LOD 减面（node_type=2）
//
// gen_times=1 表示不抽卡，生成单个版本。
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
	modelPath := assetsDir + "/tex_model.obj"

	fmt.Println("[gen_lod] 开始 LOD 减面...")

	reduceFaces := []visvise.ReduceFace{
		{ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad},
		{ReduceLevel: 2, ReducePercent: 25, FaceType: visvise.FaceTypeQuad},
		{ReduceLevel: 3, ReducePercent: 13, FaceType: visvise.FaceTypeQuad},
	}

	modelIDs, err := client.GenLOD(modelPath, reduceFaces,
		visvise.NewGenLODOptions().
			SetAlgorithmModel("VISVISE-LOD-V1.0.0").
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetGenTimes(1).
			SetName("example_gen_lod"))
	if err != nil {
		log.Fatalf("[gen_lod] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_lod] 任务已创建，model_ids=%v\n", modelIDs)

	for _, mid := range modelIDs {
		model, err := client.WaitModel(mid, &visvise.WaitOptions{
			Interval: 5.0,
			Timeout:  600,
		})
		if err != nil {
			log.Printf("[gen_lod] %s 等待失败: %v", mid, err)
			continue
		}
		fmt.Printf("[gen_lod] %s 生成成功！耗时 %ds\n", mid, model.TimeCost)
		if model.LODOutput != nil {
			for _, lf := range model.LODOutput.LODFiles {
				fmt.Printf("  LOD%d: %s\n", lf.ReduceLevel, lf.DownloadURL)
			}
		} else {
			fmt.Printf("  output_model: %s\n", model.OutputModel)
		}
	}
}
