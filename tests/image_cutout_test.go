package tests

import (
	"os"
	"testing"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// TestAtomicAPI_ImageCutout tests the image_cutout API
func TestAtomicAPI_ImageCutout(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}
	inputPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found in tests/assets")
	}

	client := visvise.NewClient(appID, secretKey, visvise.NewClientOptions().SetEnv(visvise.EnvTest))

	result, err := client.ImageCutout(
		[]visvise.FileInput{inputPath},
		visvise.ImageCutoutModeDirect,
		"white",
		rtx,
	)
	if err != nil {
		t.Fatalf("ImageCutout failed: %v", err)
	}
	if len(result.ImageList) != 1 {
		t.Fatalf("Expected 1 result image, got %d", len(result.ImageList))
	}
	t.Logf("PASS: image_cutout - background_color=%s url=%s", result.BackgroundColor, result.ImageList[0])
}
