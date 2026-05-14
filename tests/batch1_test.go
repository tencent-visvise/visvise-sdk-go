package tests

import (
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

const mvBase = "https://visvise-weaver-bj-rel-1311802504.cos.accelerate.myqcloud.com/weaver/user-p_5sxfmuvwtfj58ssbt97q2k2kc83697p/Model2026042300225069"

var mv = map[string]string{
	"main":  mvBase + "/example_gen_360_MultiView(2)_MainView.png",
	"back":  mvBase + "/example_gen_360_MultiView(2)_BackView.png",
	"left":  mvBase + "/example_gen_360_MultiView(2)_LeftView.png",
	"right": mvBase + "/example_gen_360_MultiView(2)_RightView.png",
}

// TestBatch1_Gen360APose tests gen_360 with enable_a_pose=True
func TestBatch1_Gen360APose(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGen360Options().
		SetAlgorithmModel("VISVISE-MultiView-V1.0.0").
		SetEnableAPose(true)

	modelID, err := client.Gen360(mainViewPath, opts)
	if err != nil {
		t.Fatalf("Gen360 a_pose=True failed: %v", err)
	}

	t.Logf("PASS: gen_360 a_pose=True - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 3, Timeout: 600}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_Gen360NoAPose tests gen_360 with enable_a_pose=False
func TestBatch1_Gen360NoAPose(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGen360Options().
		SetAlgorithmModel("VISVISE-MultiView-V1.0.0").
		SetEnableAPose(false)

	modelID, err := client.Gen360(mainViewPath, opts)
	if err != nil {
		t.Fatalf("Gen360 a_pose=False failed: %v", err)
	}

	t.Logf("PASS: gen_360 a_pose=False - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 3, Timeout: 600}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_HighModelFaceNum tests gen_high_model with face_num=100000
func TestBatch1_HighModelFaceNum(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenHighModelOptions().
		SetAlgorithmModel("Tripo-v3.1-ultra").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle).
		SetFaceNum(100000)

	modelID, err := client.GenHighModel(mv["main"], opts)
	if err != nil {
		t.Fatalf("GenHighModel face_num=100000 failed: %v", err)
	}

	t.Logf("PASS: high face_num=100000 - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_HighModelFaceType2 tests gen_high_model with face_type=2 and glb format
func TestBatch1_HighModelFaceType2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenHighModelOptions().
		SetAlgorithmModel("Tripo-v3.1-ultra").
		SetOutputModelFormat(visvise.OutputModelFormatGLB).
		SetFaceType(visvise.FaceTypeQuad)

	modelID, err := client.GenHighModel(mv["main"], opts)
	if err != nil {
		t.Fatalf("GenHighModel face_type=2 glb failed: %v", err)
	}

	t.Logf("PASS: high face_type=2 glb - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_LowModelFaceType1Back tests gen_low_model with face_type=1 and back_view
func TestBatch1_LowModelFaceType1Back(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenLowModelOptions().
		SetAlgorithmModel("Tripo-v1.0-快速生成").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle).
		SetBackView(mv["back"], "")

	modelID, err := client.GenLowModel(mv["main"], opts)
	if err != nil {
		t.Fatalf("GenLowModel face_type=1 back failed: %v", err)
	}

	t.Logf("PASS: low face_type=1 back - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_LowModelFaceType2 tests gen_low_model with face_type=2 and fbx format
func TestBatch1_LowModelFaceType2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenLowModelOptions().
		SetAlgorithmModel("Tripo-v1.0-快速生成").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad)

	modelID, err := client.GenLowModel(mv["main"], opts)
	if err != nil {
		t.Fatalf("GenLowModel face_type=2 fbx failed: %v", err)
	}

	t.Logf("PASS: low face_type=2 fbx - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch1_LODGenTimes tests gen_lod with gen_times=1
func TestBatch1_LODGenTimes(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad}}
	opts := visvise.NewGenLODOptions().
		SetAlgorithmModel("VISVISE-LOD-V1.0.0").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetGenTimes(1)

	modelIDs, err := client.GenLOD(modelPath, reduceFaces, opts)
	if err != nil {
		t.Fatalf("GenLOD gen_times=1 failed: %v", err)
	}

	t.Logf("PASS: lod gen_times=1 - model_ids=%v", modelIDs)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelIDs[0], waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}
