package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_segment_2d（跳过360 + 跳过去背景 + 上传六视图）—— 2D 拆分（node_type=14）
//
// 与 examples/gen_segment_2d 的区别：不依赖图生360产出的 model_id，直接上传本地原画，
// 并同时跳过 360 前置（skip_360_preprocess）与去背景预处理（skip_remove_background），
// 原始图原样送拆分算法，不再经过任何图像预处理。
//
// 视图按需传入：split_type=2（多视图）支持 1~6 个视图，素材缺失的视图留空即可，
// 后端按实际传入的视图投递；只有正视图也能跑（等价于正视图拆分）。
// 本示例会把 assets 目录下存在的视图全部带上，缺哪个就跳过哪个。
//
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod \
//   SEG_ASSETS_DIR=./tests/assets go run main.go

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

	assetsDir := os.Getenv("SEG_ASSETS_DIR")
	if assetsDir == "" {
		assetsDir = "./tests/assets"
	}

	// 六视图逐个探测，本地路径由 SDK 自动上传 COS；不存在则跳过该视图
	view := &visvise.View{}
	pick := func(name string) string {
		path := assetsDir + "/" + name
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("  [skip] 素材不存在，跳过该视图: %s\n", path)
			return ""
		}
		fmt.Printf("  [view] %s\n", path)
		return path
	}
	view.MainView = pick("main_view.png")
	if view.MainView == "" {
		log.Fatal("[gen_segment_2d_skip360] 正视图 main_view.png 必传")
	}
	view.BackView = pick("back_view.png")
	view.LeftView = pick("left_view.png")
	view.RightView = pick("right_view.png")
	view.TopView = pick("top_view.png")
	view.BottomView = pick("bottom_view.png")

	fmt.Printf("[gen_segment_2d_skip360] 上传原画并跳过 360 前置 / 去背景，直接 2D 拆分...\n")

	segModelID, err := client.GenSegment2D("", rtx,
		visvise.NewGenSegment2DOptions().
			SetName("example_gen_segment_2d_skip360").
			SetSplitType(visvise.SegmentSplitFourView). // 2 多视图拆分（按实际传入的视图，1~6 个）
			SetSkip360Preprocess(true).                 // 跳过图生360前置，直接用 input_view 拆分
			SetSkipRemoveBackground(true).              // 跳过去背景，原始图直接送算法
			SetInputView(view).
			// 可选：拆解 plan（json 文件），本地路径自动上传；不传时算法自行推理
			SetDecompositionPlan(os.Getenv("SEG_DECOMPOSITION_PLAN")).
			SetOnThinking(func(thought string) {
				fmt.Printf("  [thinking] %s\n", thought)
			}))
	if err != nil {
		log.Fatalf("[gen_segment_2d_skip360] 分割失败: %v", err)
	}

	fmt.Printf("[gen_segment_2d_skip360] 分割完成，model_id=%s\n", segModelID)
	fmt.Printf("  → 可作为图生中模/低模的 segment_model_id 参数使用\n")
}
