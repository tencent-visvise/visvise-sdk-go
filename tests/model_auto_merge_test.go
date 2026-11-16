package tests

import (
	"os"
	"testing"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// TestAtomicAPI_ModelAutoMerge tests the model_auto_merge (模型自动合并) API
func TestAtomicAPI_ModelAutoMerge(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}
	partsZip := assetsDir + "/parts.zip"
	if _, err := os.Stat(partsZip); os.IsNotExist(err) {
		t.Skip("Skipping test: parts.zip not found in tests/assets")
	}

	// 同步算法耗时较长，超时至少 300s
	client := visvise.NewClient(appID, secretKey, visvise.NewClientOptions().SetEnv(visvise.EnvTest).SetTimeout(360))

	result, err := client.ModelAutoMerge(partsZip, "", false, rtx)
	if err != nil {
		t.Fatalf("ModelAutoMerge failed: %v", err)
	}
	if result.OutputObjPath == "" {
		t.Fatal("ModelAutoMerge returned empty output_obj_path")
	}
	t.Logf("PASS: model_auto_merge - elapsed=%.2fs mean_iou=%.4f url=%s",
		result.ElapsedTime, result.MeanBboxIou, result.OutputObjPath)
}
