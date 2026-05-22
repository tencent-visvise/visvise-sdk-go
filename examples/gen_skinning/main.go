package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_skinning —— 蒙皮生成（node_type=6）
//
// 为带骨骼的 3D 模型自动生成蒙皮权重。
// SDK 内部自动构建 JSON 参数文件（含 mesh_names / joint_names）并打包成 zip 上传。
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
	modelPath := assetsDir + "/skinning_model.fbx"

	// 来自 skinning_model.json
	meshNames := []string{
		"CH_M_ZB_Human_01_Body",
		"CH_M_ZB_Human_01_Hair",
		"CH_M_ZB_Human_01_Weapon",
	}

	jointNames := []string{
		"Bip002",
		"Bip002 Pelvis",
		"Bip002 Spine",
		"Bip002 Spine1",
		"Bip002 Spine2",
		"Bip002 Neck",
		"Bip002 L Clavicle",
		"Bip002 L UpperArm",
		"Bip002 L Forearm",
		"Bip002 L Hand",
		"Bip002 R Clavicle",
		"Bip002 R UpperArm",
		"Bip002 R Forearm",
		"Bip002 R Hand",
		"Bip002 L Thigh",
		"Bip002 L Calf",
		"Bip002 L Foot",
		"Bip002 R Thigh",
		"Bip002 R Calf",
		"Bip002 R Foot",
	}

	fmt.Println("[gen_skinning] 开始蒙皮生成...")

	modelID, err := client.GenSkinning(modelPath, rtx,
		visvise.NewGenSkinningOptions(meshNames, jointNames).
			SetName("example_gen_skinning"))
	if err != nil {
		log.Fatalf("[gen_skinning] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_skinning] 任务已创建，model_id=%s\n", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  600,
	})
	if err != nil {
		log.Fatalf("[gen_skinning] 等待失败: %v", err)
	}

	fmt.Printf("[gen_skinning] 生成成功！耗时 %ds\n", model.TimeCost)
	fmt.Printf("  output_model : %s\n", model.OutputModel)
}
