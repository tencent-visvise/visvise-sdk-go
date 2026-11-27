package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// Example: gen_rtp_connected —— 重拓扑-连通体（node_type=19）
//
// 两段式自研（VV）重拓扑的第一段：对高面数模型做连通体分组，
// 产出连通体网格与分组数据（产物在资产的 output_model zip 中）。
//
// 完整链路：
//  1. GenRTPConnected：创建连通体分组任务（node_type=19）
//  2. 等待完成，下载 output_model（连通体分组产物 zip）并在前端编辑连通体分组
//  3. 将编辑结果重新打包成 zip 并上传到当前账号的 COS 目录
//  4. SaveRTPConnected：保存该产物，得到新的连通体资产 ID（node_type=19，status=3）
//  5. GenRetopology：以保存得到的 model_id 作为 parent_model_id 发起第二段重拓扑（node_type=1）
//
// Usage:
//
//	VISVISE_APP_ID=xxx VISVISE_SECRET_KEY=xxx VISVISE_RTX=xxx VISVISE_ENV=prod go run main.go
//
// 可选环境变量：
//
//	VISVISE_RTP_CONNECTED_ZIP  前端编辑后重新打包的连通体产物 zip 路径（本地文件或 COS URL）
//	                             传入后本示例会继续执行「保存产物 + 第二段重拓扑」
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

	// ── 第一段：重拓扑-连通体 ────────────────────────────────────────────
	fmt.Println("[gen_rtp_connected] 创建重拓扑-连通体任务（node_type=19）...")

	connectedOpts := visvise.NewGenRTPConnectedOptions().
		SetName("example_gen_rtp_connected")

	connectedID, err := client.GenRTPConnected(modelPath, rtx, connectedOpts)
	if err != nil {
		log.Fatalf("[gen_rtp_connected] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_rtp_connected] 任务已创建，model_id=%s\n", connectedID)

	fmt.Println("[gen_rtp_connected] 等待完成...")
	connected, err := client.WaitModel(connectedID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_rtp_connected] 等待失败: %v", err)
	}
	fmt.Printf("[gen_rtp_connected] 连通体分组产物: %s\n", connected.OutputModel)
	fmt.Println("  请下载该 zip，在前端编辑连通体分组后重新打包上传到当前账号的 COS 目录，")
	fmt.Println("  然后调用 SaveRTPConnected 保存产物，并以返回的 model_id 作为重拓扑的 parent_model_id。")

	// ── 保存产物 + 第二段重拓扑 ─────────────────────────────────────────
	editedZip := os.Getenv("VISVISE_RTP_CONNECTED_ZIP")
	if editedZip == "" {
		fmt.Println("\n未设置 VISVISE_RTP_CONNECTED_ZIP，跳过「保存产物 + 第二段重拓扑」。")
		return
	}

	fmt.Println("\n[save_rtp_connected] 保存编辑后的连通体产物...")
	savedID, err := client.SaveRTPConnected(connectedID, editedZip, "example_rtp_connected_saved", rtx)
	if err != nil {
		log.Fatalf("[save_rtp_connected] 保存失败: %v", err)
	}
	fmt.Printf("[save_rtp_connected] 保存成功，新的连通体资产 ID=%s\n", savedID)

	fmt.Println("[gen_retopology] 发起第二段重拓扑（node_type=1）...")
	rtpOpts := visvise.NewGenRetopologyOptions().
		SetName("example_gen_retopology_vv").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad).
		SetFaceNum(10000).        // VISVISE 自研模型用 face_num
		SetSmoothStrength(1.0).   // 平滑强度
		SetParentModelID(savedID) // 第一段保存后的连通体资产 ID

	rtpID, err := client.GenRetopology(editedZip, rtx, rtpOpts)
	if err != nil {
		log.Fatalf("[gen_retopology] 创建任务失败: %v", err)
	}
	fmt.Printf("[gen_retopology] 任务已创建，model_id=%s\n", rtpID)

	rtpModel, err := client.WaitModel(rtpID, rtx, &visvise.WaitOptions{
		Interval: 5.0,
		Timeout:  900,
	})
	if err != nil {
		log.Fatalf("[gen_retopology] 等待失败: %v", err)
	}
	fmt.Printf("[gen_retopology] 生成成功！耗时 %ds\n", rtpModel.TimeCost)
	fmt.Printf("  output_model : %s\n", rtpModel.OutputModel)
}
