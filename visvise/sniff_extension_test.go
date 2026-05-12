package visvise

import (
	"os"
	"path/filepath"
	"testing"
)

// Expected extensions for real asset files
// Note: main_view.png contains JPEG data despite .png extension (same MD5 as main_view.jpg)
var expectedExtensions = map[string]string{
	"main_view.png":       ".png", // contains JPEG data
	"back_view.png":       ".png",
	"left_view.png":       ".png",
	"right_view.png":      ".png",
	"pose_ref.png":        ".png",
	"tex_ref_front.jpg":   ".jpg",
	"tex_ref_back.jpg":    ".jpg",
	"tex_ref_left.jpg":    ".jpg",
	"high_model.fbx":      ".fbx",
	"animation_model.fbx": ".fbx",
	"rigging_model.fbx":   ".fbx",
	"skinning_model.fbx":  ".fbx",
	"tex_model.obj":       ".obj",
	"animation_video.mp4": ".mp4",
}

// TestSniffExtension_RealAssets tests sniffExtension and genRandomFilename with real asset files
func TestSniffExtension_RealAssets(t *testing.T) {
	assetsDir := "../tests/assets"
	failures := []string{}

	for fname, expectExt := range expectedExtensions {
		path := filepath.Join(assetsDir, fname)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Logf("SKIP: %s not found, skipping", fname)
			continue
		}

		// Read only first 2KB to identify file type (memory efficient)
		file, err := os.Open(path)
		if err != nil {
			failures = append(failures, fname+": failed to open file: "+err.Error())
			continue
		}

		head := make([]byte, 2048)
		n, err := file.Read(head)
		file.Close()
		if err != nil && n == 0 {
			failures = append(failures, fname+": failed to read file: "+err.Error())
			continue
		}
		head = head[:n]

		gotExt := sniffExtension(head, ".bin")
		genName := genRandomFilename(".bin")

		sizeKB := float64(0)
		if info, err := os.Stat(path); err == nil {
			sizeKB = float64(info.Size()) / 1024
		}

		status := "PASS"
		if gotExt != expectExt {
			status = "FAIL"
			failures = append(failures, fname+": got="+gotExt+", want="+expectExt)
		}
		t.Logf("  %s %s (%.1f KB) -> %s (expected: %s) uuid: %s", status, fname, sizeKB, gotExt, expectExt, genName)
	}

	if len(failures) > 0 {
		t.Errorf("%d test cases failed:", len(failures))
		for _, f := range failures {
			t.Errorf("  - %s", f)
		}
	}
}
