package tests

import (
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// TestAutoAlgorithmModel_Gen360 tests gen_360 without algorithm_model parameter
func TestAutoAlgorithmModel_Gen360(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found in tests/assets")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGen360Options()

	modelID, err := client.Gen360(mainViewPath, opts)
	if err != nil {
		t.Fatalf("Gen360 failed: %v", err)
	}

	t.Logf("PASS: gen_360 (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenHighModel tests gen_high_model without algorithm_model parameter
func TestAutoAlgorithmModel_GenHighModel(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found in tests/assets")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenHighModelOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle)

	modelID, err := client.GenHighModel(mainViewPath, opts)
	if err != nil {
		t.Fatalf("GenHighModel failed: %v", err)
	}

	t.Logf("PASS: gen_high_model (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenMidModel tests gen_mid_model without algorithm_model parameter
func TestAutoAlgorithmModel_GenMidModel(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	backViewPath := assetsDir + "/back_view.png"
	leftViewPath := assetsDir + "/left_view.png"
	rightViewPath := assetsDir + "/right_view.png"

	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMidModelOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle)

	modelID, err := client.GenMidModel(mainViewPath, backViewPath, leftViewPath, rightViewPath, opts)
	if err != nil {
		t.Fatalf("GenMidModel failed: %v", err)
	}

	t.Logf("PASS: gen_mid_model (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenLowModel tests gen_low_model without algorithm_model parameter
func TestAutoAlgorithmModel_GenLowModel(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenLowModelOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle)

	modelID, err := client.GenLowModel(mainViewPath, opts)
	if err != nil {
		t.Fatalf("GenLowModel failed: %v", err)
	}

	t.Logf("PASS: gen_low_model (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenMeshRefine tests gen_mesh_refine without algorithm_model parameter
func TestAutoAlgorithmModel_GenMeshRefine(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMeshRefineOptions()

	modelID, err := client.GenMeshRefine(modelPath, opts)
	if err != nil {
		t.Fatalf("GenMeshRefine failed: %v", err)
	}

	t.Logf("PASS: gen_mesh_refine (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenRetopology tests gen_retopology without algorithm_model parameter
func TestAutoAlgorithmModel_GenRetopology(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenRetopologyOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad).
		SetDetailLevel(visvise.DetailLevelMedium)

	modelID, err := client.GenRetopology(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRetopology failed: %v", err)
	}

	t.Logf("PASS: gen_retopology (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenLOD tests gen_lod without algorithm_model parameter
func TestAutoAlgorithmModel_GenLOD(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad}}
	opts := visvise.NewGenLODOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetGenTimes(1)

	modelIDs, err := client.GenLOD(modelPath, reduceFaces, opts)
	if err != nil {
		t.Fatalf("GenLOD failed: %v", err)
	}

	t.Logf("PASS: gen_lod (no algorithm_model) - model_ids=%v", modelIDs)
}

// TestAutoAlgorithmModel_GenUV tests gen_uv without algorithm_model parameter
func TestAutoAlgorithmModel_GenUV(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenUVOptions().
		SetEnableAutoSmoothing(false)

	modelID, err := client.GenUV(modelPath, opts)
	if err != nil {
		t.Fatalf("GenUV failed: %v", err)
	}

	t.Logf("PASS: gen_uv (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenTexture tests gen_texture without algorithm_model parameter
func TestAutoAlgorithmModel_GenTexture(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	view := &visvise.View{MainView: mainViewPath}
	opts := visvise.NewGenTextureOptions().
		SetInputView(view)

	modelID, err := client.GenTexture(modelPath, opts)
	if err != nil {
		t.Fatalf("GenTexture failed: %v", err)
	}

	t.Logf("PASS: gen_texture (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenRigging tests gen_rigging without algorithm_model parameter
func TestAutoAlgorithmModel_GenRigging(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenRiggingOptions()

	modelID, err := client.GenRigging(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRigging failed: %v", err)
	}

	t.Logf("PASS: gen_rigging (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenSkinning tests gen_skinning without algorithm_model parameter
func TestAutoAlgorithmModel_GenSkinning(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/skinning_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: skinning_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenSkinningOptions([]string{"Body_Mesh"}, []string{"Bip001", "Bip001 Pelvis"})

	modelID, err := client.GenSkinning(modelPath, opts)
	if err != nil {
		t.Fatalf("GenSkinning failed: %v", err)
	}

	t.Logf("PASS: gen_skinning (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenVideoMotion tests gen_video_motion without algorithm_model parameter
func TestAutoAlgorithmModel_GenVideoMotion(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/animation_model.fbx"
	videoPath := assetsDir + "/animation_video.mp4"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenVideoMotionOptions().
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetWithHand(false).
		SetMultipleTrack(false)

	modelID, err := client.GenVideoMotion(modelPath, videoPath, opts)
	if err != nil {
		t.Fatalf("GenVideoMotion failed: %v", err)
	}

	t.Logf("PASS: gen_video_motion (no algorithm_model) - model_id=%s", modelID)
}

// TestAutoAlgorithmModel_GenTextMotion tests gen_text_motion without algorithm_model parameter
func TestAutoAlgorithmModel_GenTextMotion(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/animation_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenTextMotionOptions()

	modelIDs, err := client.GenTextMotion(modelPath, "一个人在原地踏步", opts)
	if err != nil {
		t.Fatalf("GenTextMotion failed: %v", err)
	}

	t.Logf("PASS: gen_text_motion (no algorithm_model) - model_ids=%v", modelIDs)
}

// TestAutoAlgorithmModel_GenPose tests gen_pose without algorithm_model parameter
func TestAutoAlgorithmModel_GenPose(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/animation_model.fbx"
	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenPoseOptions().
		SetImageFilenames([]string{"main_view.png"})

	modelIDs, err := client.GenPose(modelPath, []visvise.FileInput{mainViewPath}, opts)
	if err != nil {
		t.Fatalf("GenPose failed: %v", err)
	}

	t.Logf("PASS: gen_pose (no algorithm_model) - model_ids=%v", modelIDs)
}
