package main

import (
	"fmt"
	"log"
	"os"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// Example: gen_pose —— 批量图生 Pose（node_type=12）
//
// 从参考图片中提取姿态，驱动 3D 模型生成对应 Pose。
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
	modelPath := assetsDir + "/animation_model.fbx"
	poseRef := assetsDir + "/pose_ref.png"

	fmt.Println("[gen_pose] 开始图生 Pose...")

	// []visvise.FileInput 类型转换
	inputImages := []visvise.FileInput{poseRef}

	modelIDs, err := client.GenPose(modelPath, inputImages,
		visvise.NewGenPoseOptions().
			SetAlgorithmModel("VISVISE-PosingAI-V1.0.0").
			SetOutputModelFormat(visvise.OutputModelFormatFBX).
			SetName("example_gen_pose"))
	if err != nil {
		log.Fatalf("[gen_pose] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_pose] 任务已创建，model_ids=%v\n", modelIDs)

	for _, mid := range modelIDs {
		model, err := client.WaitModel(mid, &visvise.WaitOptions{
			Interval: 3.0,
			Timeout:  600,
		})
		if err != nil {
			log.Printf("[gen_pose] %s 等待失败: %v", mid, err)
			continue
		}
		fmt.Printf("[gen_pose] %s 生成成功！耗时 %ds\n", mid, model.TimeCost)
		fmt.Printf("  output_model : %s\n", model.OutputModel)
	}
}
