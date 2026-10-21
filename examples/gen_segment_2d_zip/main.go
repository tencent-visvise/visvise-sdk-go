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
// zip_path 是算法侧产出的拆分结果压缩包（COS 路径），后端只落库不解析，里面装什么由算法定义。
// 后端对外返回时会自动加临时签名 ⇒ **浏览器打开 / curl / 代码 GET 都能直接下载，无需鉴权头**。
//
// 方式一（推荐，最简单）：GenSegment2D + WaitModel
//   GenSegment2D  → modelID
//   WaitModel     → ModelInfo.SegmentOutput.ZipPath
//   顺带还能拿到各视图的部件表 Components / npz DataPath / 彩色预览图 ColorImgs。
//   链接在「查询那一刻」现签，24h 内有效，过期重新查一次详情即可拿新链接。
//
// 方式二（生成瞬间就拿到）：底层 api.InitSegment 自己消费 SSE 帧
//   高层 GenSegment2D 只把 model_id 返回、reply 帧被丢弃；
//   想在同一条 SSE 连接里直接取 zip_path，就自己遍历事件帧拿 reply。
//   同一个文件，只是签名时刻更早（从生成那一刻起算 24h），省掉一次查询。
//
// Usage:
//   VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod \
//   MV_360_MODEL_ID=Model2026xxxx go run main.go

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
		log.Fatal("[gen_segment_2d_zip] 需要先运行 examples/gen_360 并设置环境变量 MV_360_MODEL_ID")
	}

	outDir := os.Getenv("SEG_OUT_DIR")
	if outDir == "" {
		outDir = "."
	}

	way1HighLevel(client, mvModelID, rtx, outDir)
	fmt.Println()
	way2SSEReply(client, mvModelID, rtx, outDir)

	fmt.Println("\n[提示] zip_path 是带签名的 URL，24h 内直接 GET 即可下载；" +
		"过期后重新调用 WaitModel / GetModelList 就能拿到新链接。")
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

// way1HighLevel 方式一：高层方法 + 资产详情（WaitModel），最省事的一种。
func way1HighLevel(client *visvise.Client, mvModelID, rtx, outDir string) {
	fmt.Println("[方式一] GenSegment2D → WaitModel → SegmentOutput.ZipPath")

	segModelID, err := client.GenSegment2D(mvModelID, rtx,
		visvise.NewGenSegment2DOptions().
			SetName("example_gen_segment_2d_zip_way1").
			SetSplitType(visvise.SegmentSplitFourView)) // 2 多视图拆分（按实际传入的视图，1~6 个）
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

	so := modelInfo.SegmentOutput
	if so == nil {
		fmt.Println("  [warn] SegmentOutput 为空：只有 2D 拆分（node_type=14）资产才有拆分产物")
		return
	}
	fmt.Printf("  zip_path = %s\n", so.ZipPath)

	// 顺带看下拆分结果里的其它内容（这才是 SDK 能拿到的最全信息）
	if so.MainView != nil {
		fmt.Printf("  正视图 shape=%v\n", so.MainView.Shape)
		fmt.Printf("  正视图 npz=%s\n", so.MainView.DataPath)
		fmt.Printf("  正视图 mask PNG=%s\n", so.MainView.MaskImagePath)
		for _, comp := range so.MainView.Components {
			fmt.Printf("    部件 label=%d name=%s color=%s  ← label 可作为 ComponentLabel 送图生中模\n",
				comp.Label, comp.Name, comp.Color)
		}
	}
	if so.ColorImgs != nil && so.ColorImgs.MainView != "" {
		fmt.Printf("  正视图彩色预览图=%s\n", so.ColorImgs.MainView)
	}

	download(so.ZipPath, filepath.Join(outDir, "segment_output_way1.zip"))
}

// way2SSEReply 方式二：自己消费 SSE 流，在 reply 帧里直接拿 zip_path（不用再查一次）。
func way2SSEReply(client *visvise.Client, mvModelID, rtx, outDir string) {
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

	splitType := visvise.SegmentSplitFourView
	iter, err := api.InitSegment(
		"example_gen_segment_2d_zip_way2", // name
		algorithmModel,                    // algorithm_model
		mvModelID,                         // model_id（图生360 的 model_id）
		nil,                               // input_view（传 model_id 时留空）
		&splitType,                        // split_type
		nil,                               // params（skip_360 / skip_remove_background / plan）
		0,                                 // readTimeout，0 表示用默认 20 分钟
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
