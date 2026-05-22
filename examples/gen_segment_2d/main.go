package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_segment_2d —— 2D 拆分（node_type=14）
//
// 对图生 360 输出的多视图进行组件分割，生成的分割资产可作为后续图生中模 / 低模的
// segment_model_id 输入，用于基于分割结果的精细化生成。
//
// 使用 SSE 协议，过程中会推送 thinking 思考态。
//
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod MV_360_MODEL_ID=xxx go run main.go

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

	// 需要先有 gen_360 输出的 model_id
	mvModelID := os.Getenv("MV_360_MODEL_ID")
	if mvModelID == "" {
		fmt.Println("[gen_segment_2d] 需要先运行 gen_360 并设置环境变量 MV_360_MODEL_ID")
		fmt.Println("Usage: MV_360_MODEL_ID=xxx go run main.go")
		os.Exit(1)
	}

	fmt.Printf("[gen_segment_2d] 基于 360 模型 %s 进行 2D 拆分...\n", mvModelID)

	segModelID, err := client.GenSegment2D(mvModelID, rtx,
		visvise.NewGenSegment2DOptions().
			SetSplitType(visvise.SegmentSplitFrontView).      // 1 正视图（默认） / 2 四视图
			SetGranularity(visvise.SegmentGranularityMedium). // 1 粗 / 2 中（默认） / 3 细
			SetPrompt("").                                    // 可选：自然语言描述拆分规则
			SetName("example_gen_segment_2d").
			SetOnThinking(func(thought string) {
				fmt.Printf("  [thinking] %s\n", thought)
			}))
	if err != nil {
		log.Fatalf("[gen_segment_2d] 分割失败: %v", err)
	}

	fmt.Printf("[gen_segment_2d] 分割完成，model_id=%s\n", segModelID)
	fmt.Printf("  → 可作为图生中模/低模的 segment_model_id 参数使用\n")
}
