package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: 2D 拆分产物 zip_path 的两种获取方式（node_type=14）
//
// 输入方式：直接上传本地六视图原画，不走图生360
//   InputView              六视图（正视图必传，其余按需，缺哪个就跳过哪个）
//   Skip360Preprocess      true —— 跳过图生360前置，直接用 input_view 拆分
//   SkipRemoveBackground   true —— 跳过去背景，原始图原样送拆分算法
//   SplitType              2（多视图拆分），支持 1~6 个视图，按实际传入的视图投递
//
// zip_path 是算法侧产出的拆分结果压缩包（COS 路径），后端只落库不解析，里面装什么由算法定义。
// 后端对外返回时会自动加临时签名 ⇒ **浏览器打开 / curl / 代码 GET 都能直接下载，无需鉴权头**。
//
// 方式一（推荐，最简单）：GenSegment2D + WaitModel
//   GenSegment2D  → modelID（本地文件由 SDK 自动上传 COS）
//   WaitModel     → ModelInfo.SegmentOutput.ZipPath
//   顺带还能拿到各视图的部件表 Components / npz DataPath / 彩色预览图 ColorImgs。
//   链接在「查询那一刻」现签，24h 内有效，过期重新查一次详情即可拿新链接。
//
// 方式二（生成瞬间就拿到）：底层 api.InitSegment 自己消费 SSE 帧
//   高层 GenSegment2D 只把 model_id 返回、reply 帧被丢弃；
//   想在同一条 SSE 连接里直接取 zip_path，就自己遍历事件帧拿 reply。
//   注意底层接口只接受「已上传到平台的 COS URL」，本地文件需要先自己上传一次。
//
// 两种方式拿到的是同一个 COS 文件，区别只在「什么时候签的名」和「要不要多查一次接口」。
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
	inputView := buildInputView(assetsDir)

	outDir := os.Getenv("SEG_OUT_DIR")
	if outDir == "" {
		outDir = "."
	}

	way1HighLevel(client, inputView, rtx, outDir)
	fmt.Println()
	way2SSEReply(client, inputView, rtx, outDir)

	fmt.Println("\n[提示] zip_path 是带签名的 URL，24h 内直接 GET 即可下载；" +
		"过期后重新调用 WaitModel / GetModelList 就能拿到新链接。")
}

// buildInputView 逐个探测六视图素材，存在则带上，不存在则跳过该视图。
func buildInputView(assetsDir string) *visvise.View {
	view := &visvise.View{}
	pick := func(name string) string {
		path := filepath.Join(assetsDir, name)
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("  [skip] 素材不存在，跳过该视图: %s\n", path)
			return ""
		}
		fmt.Printf("  [view] %s\n", path)
		return path
	}
	view.MainView = pick("main_view.png")
	if view.MainView == "" {
		log.Fatal("[gen_segment_2d_zip] 正视图 main_view.png 必传")
	}
	view.BackView = pick("back_view.png")
	view.LeftView = pick("left_view.png")
	view.RightView = pick("right_view.png")
	view.TopView = pick("top_view.png")
	view.BottomView = pick("bottom_view.png")
	return view
}

// download zip_path 是带 COS 签名的 URL，直接 GET 即可，无需任何鉴权头。
func download(url, saveAs string) {
	if url == "" {
		fmt.Println("  [warn] zip_path 为空（该资产没有算法产物，或后端未落库），跳过下载")
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  [warn] 下载失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		fmt.Printf("  [warn] 下载失败: HTTP %d, %s\n", resp.StatusCode, string(body))
		return
	}

	f, err := os.Create(saveAs)
	if err != nil {
		fmt.Printf("  [warn] 创建文件失败: %v\n", err)
		return
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		fmt.Printf("  [warn] 写入文件失败: %v\n", err)
		return
	}
	fmt.Printf("  [download] %s（%d bytes）\n", saveAs, n)
}

// printSegmentOutput 打印各视图的分割数据（没参与拆分的视图为 nil）。
func printSegmentOutput(so *visvise.SegmentOutput) {
	if so == nil {
		fmt.Println("  [warn] SegmentOutput 为空：只有 2D 拆分（node_type=14）资产才有拆分产物")
		return
	}
	fmt.Printf("  zip_path = %s\n", so.ZipPath)

	views := []struct {
		name string
		v    *visvise.SegmentViewData
	}{
		{"main_view", so.MainView},
		{"back_view", so.BackView},
		{"left_view", so.LeftView},
		{"right_view", so.RightView},
		{"top_view", so.TopView},
		{"bottom_view", so.BottomView},
	}
	for _, item := range views {
		if item.v == nil {
			continue // 该视图没参与拆分
		}
		fmt.Printf("  %s: shape=%v npz=%s mask PNG=%s\n",
			item.name, item.v.Shape, item.v.DataPath, item.v.MaskImagePath)
		for _, comp := range item.v.Components {
			fmt.Printf("    部件 label=%d name=%s color=%s  ← label 可作为 ComponentLabel 送图生中模\n",
				comp.Label, comp.Name, comp.Color)
		}
	}
	if so.ColorImgs != nil && so.ColorImgs.MainView != "" {
		fmt.Printf("  正视图彩色预览图=%s\n", so.ColorImgs.MainView)
	}
}

// way1HighLevel 方式一：高层方法 + 资产详情（WaitModel），最省事的一种。
func way1HighLevel(client *visvise.Client, inputView *visvise.View, rtx, outDir string) {
	fmt.Println("[方式一] GenSegment2D → WaitModel → SegmentOutput.ZipPath")

	segModelID, err := client.GenSegment2D("", rtx, // 不用图生360 的 model_id
		visvise.NewGenSegment2DOptions().
			SetName("example_gen_segment_2d_zip_way1").
			SetSplitType(visvise.SegmentSplitFourView). // 2 多视图拆分（按实际传入的视图，1~6 个）
			SetSkip360Preprocess(true).                 // 跳过图生360前置，直接用 input_view 拆分
			SetSkipRemoveBackground(true).              // 跳过去背景，原始图直接送算法
			SetInputView(inputView))                    // 本地文件由 SDK 自动上传 COS
	if err != nil {
		log.Fatalf("[方式一] 分割失败: %v", err)
	}
	fmt.Printf("  model_id=%s\n", segModelID)

	// WaitModel 轮询到 status=3 才返回，此时 ZipPath 一定已落库
	modelInfo, err := client.WaitModel(segModelID, rtx,
		&visvise.WaitOptions{Interval: 2, Timeout: 600})
	if err != nil {
		log.Fatalf("[方式一] 等待完成失败: %v", err)
	}

	printSegmentOutput(modelInfo.SegmentOutput)
	if modelInfo.SegmentOutput != nil {
		download(modelInfo.SegmentOutput.ZipPath, filepath.Join(outDir, "segment_output_way1.zip"))
	}
}

// way2SSEReply 方式二：自己消费 SSE 流，在 reply 帧里直接拿 zip_path（不用再查一次）。
func way2SSEReply(client *visvise.Client, inputView *visvise.View, rtx, outDir string) {
	fmt.Println("[方式二] GetAPI().InitSegment 自己遍历 SSE 帧 → reply 帧里的 zip_path")

	api := client.GetAPI()

	// 底层接口需要显式传算法模型名，高层 GenSegment2D 会自动选，这里手动取一个
	models, err := api.ListAlgorithmModel(int(visvise.NodeTypeSegment2D), nil, rtx)
	if err != nil || len(models) == 0 {
		fmt.Printf("  [warn] 取算法模型失败或当前账号没有可用的 2D 拆分模型: %v\n", err)
		return
	}
	algorithmModel := models[0]
	fmt.Printf("  algorithm_model=%s\n", algorithmModel)

	// 底层接口只接受已上传到平台的 COS URL，本地文件先上传一次
	cosView, err := uploadView(client, inputView, rtx)
	if err != nil {
		fmt.Printf("  [warn] 上传视图失败: %v\n", err)
		return
	}

	splitType := visvise.SegmentSplitFourView
	skip360, skipRembg := true, true
	iter, err := api.InitSegment(
		"example_gen_segment_2d_zip_way2", // name
		algorithmModel,                    // algorithm_model
		"",                                // model_id：不走图生360，留空
		cosView,                           // input_view（六视图的 COS URL）
		&splitType,                        // split_type
		&visvise.SegmentParams{            // params（skip_360 / skip_remove_background / plan）
			Skip360Preprocess:    &skip360,
			SkipRemoveBackground: &skipRembg,
		},
		0,   // readTimeout，0 表示用默认 20 分钟
		rtx,
	)
	if err != nil {
		fmt.Printf("  [warn] 发起拆分失败: %v\n", err)
		return
	}
	defer iter.Close()

	var zipPath, clientID string
	for {
		event, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("  [warn] 读取 SSE 失败: %v\n", err)
			return
		}

		switch event.Event {
		case "pre_create":
			if m, ok := event.Data.(map[string]interface{}); ok {
				if id, ok := m["model_id"].(string); ok {
					fmt.Printf("  [pre_create] model_id=%s\n", id)
				}
			}
		case "thinking":
			if s, ok := event.Data.(string); ok {
				fmt.Printf("  [thinking] %s\n", s)
			}
		case "reply":
			if m, ok := event.Data.(map[string]interface{}); ok {
				if s, ok := m["zip_path"].(string); ok {
					zipPath = s
				}
				if s, ok := m["client_id"].(string); ok {
					clientID = s
				}
			}
			fmt.Printf("  [reply] zip_path=%s\n", zipPath)
		case "error":
			fmt.Printf("  [error] %v\n", event.Data)
			return
		}
		if zipPath != "" {
			break
		}
	}

	fmt.Printf("  client_id=%s（可传给 OpenSegment 做二次编辑）\n", clientID)
	download(zipPath, filepath.Join(outDir, "segment_output_way2.zip"))
}

// uploadView 把本地六视图逐个上传到平台 COS，返回只含 COS URL 的 View。
func uploadView(client *visvise.Client, view *visvise.View, rtx string) (*visvise.View, error) {
	uploaded := &visvise.View{}
	items := []struct {
		src *string
		dst *string
	}{
		{&view.MainView, &uploaded.MainView},
		{&view.BackView, &uploaded.BackView},
		{&view.LeftView, &uploaded.LeftView},
		{&view.RightView, &uploaded.RightView},
		{&view.TopView, &uploaded.TopView},
		{&view.BottomView, &uploaded.BottomView},
	}
	for _, item := range items {
		if *item.src == "" {
			continue
		}
		url, err := client.UploadFile(*item.src, filepath.Base(*item.src), false, rtx)
		if err != nil {
			return nil, err
		}
		*item.dst = url
	}
	return uploaded, nil
}
