package tests

import (
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// TestBatch3_VideoMotionWithHand tests gen_video_motion with with_hand=True
func TestBatch3_VideoMotionWithHand(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	videoPath := assetsDir + "/animation_video.mp4"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenVideoMotionOptions().
		SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0").
		SetWithHand(true).
		SetMultipleTrack(false)

	modelID, err := client.GenVideoMotion(animModelPath, videoPath, rtx, opts)
	if err != nil {
		t.Fatalf("GenVideoMotion with_hand=True failed: %v", err)
	}

	t.Logf("PASS: vm with_hand=True - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, rtx, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch3_VideoMotionNoHandNoMulti tests gen_video_motion with with_hand=False and multiple_track=False
func TestBatch3_VideoMotionNoHandNoMulti(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	videoPath := assetsDir + "/animation_video.mp4"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenVideoMotionOptions().
		SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0").
		SetWithHand(false).
		SetMultipleTrack(false)

	modelID, err := client.GenVideoMotion(animModelPath, videoPath, rtx, opts)
	if err != nil {
		t.Fatalf("GenVideoMotion hand=False multi=False failed: %v", err)
	}

	t.Logf("PASS: vm hand=False multi=False - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, rtx, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch3_TextMotionWave tests gen_text_motion with prompt="挥手"
func TestBatch3_TextMotionWave(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenTextMotionOptions().
		SetAlgorithmModel("VISVISE-TextMotion-V1.1.0")

	modelIDs, err := client.GenTextMotion(animModelPath, "一个人在挥手打招呼", rtx, opts)
	if err != nil {
		t.Fatalf("GenTextMotion prompt=挥手 failed: %v", err)
	}

	t.Logf("PASS: tm prompt=挥手 - model_ids=%v", modelIDs)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelIDs[0], rtx, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch3_TextMotionStep tests gen_text_motion with prompt="踏步" and glb format
func TestBatch3_TextMotionStep(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenTextMotionOptions().
		SetAlgorithmModel("VISVISE-TextMotion-V1.1.0").
		SetOutputModelFormat(visvise.OutputModelFormatGLB)

	modelIDs, err := client.GenTextMotion(animModelPath, "一个人在原地踏步", rtx, opts)
	if err != nil {
		t.Fatalf("GenTextMotion prompt=踏步 glb failed: %v", err)
	}

	t.Logf("PASS: tm prompt=踏步 glb - model_ids=%v", modelIDs)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelIDs[0], rtx, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}
