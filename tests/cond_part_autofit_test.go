package tests

import (
	"os"
	"testing"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// TestAtomicAPI_CondPartAutofit tests the cond_part_autofit (裸模拆分) API
func TestAtomicAPI_CondPartAutofit(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}
	originalPath := assetsDir + "/original_model.obj"
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		t.Skip("Skipping test: original_model.obj not found in tests/assets")
	}

	// 同步算法耗时较长，超时至少 300s
	client := visvise.NewClient(appID, secretKey, visvise.NewClientOptions().SetEnv(visvise.EnvTest).SetTimeout(360))

	result, err := client.CondPartAutofit(originalPath, originalPath, nil, rtx)
	if err != nil {
		t.Fatalf("CondPartAutofit failed: %v", err)
	}
	if result.OutputObjPath == "" {
		t.Fatal("CondPartAutofit returned empty output_obj_path")
	}
	t.Logf("PASS: cond_part_autofit - background_color=%s url=%s", result.BackgroundColor, result.OutputObjPath)
}
