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

	modelID, err := client.Gen360(mainViewPath, "", "test_auto_alg_360", "main_view.png", nil, "", nil, "", nil, "", nil, "")
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

	modelID, err := client.GenHighModel(mainViewPath, "", string(visvise.OutputModelFormatFBX), 1, "test_auto_alg_high", "main_view.png", nil, nil, "", nil, "", nil, "")
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

	modelID, err := client.GenMidModel(mainViewPath, backViewPath, leftViewPath, rightViewPath, "", string(visvise.OutputModelFormatFBX), 1, "test_auto_alg_mid", "main_view.png", "back_view.png", "left_view.png", "right_view.png", "")
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

	modelID, err := client.GenLowModel(mainViewPath, "", string(visvise.OutputModelFormatFBX), 1, "test_auto_alg_low", "main_view.png", nil, "", nil, "", nil, "")
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

	modelID, err := client.GenMeshRefine(modelPath, "", "test_auto_alg_mesh_refine", "", "", nil, nil, "")
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

	modelID, err := client.GenRetopology(modelPath, "", string(visvise.OutputModelFormatFBX), 2, "test_auto_alg_retopology", "", intPtr(2), nil)
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

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: 2}}
	modelIDs, err := client.GenLOD(modelPath, reduceFaces, "", string(visvise.OutputModelFormatFBX), "test_auto_alg_lod", "", 1)
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

	modelID, err := client.GenUV(modelPath, "", "test_auto_alg_uv", "", boolPtr(false))
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
	modelID, err := client.GenTexture(modelPath, "", "test_auto_alg_texture", "", view, intPtr(0), nil, "")
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

	modelID, err := client.GenRigging(modelPath, "", "", "test_auto_alg_rigging", "", nil, "")
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

	modelID, err := client.GenSkinning(modelPath, []string{"Body_Mesh"}, []string{"Bip001", "Bip001 Pelvis"}, "", "test_auto_alg_skinning", "")
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

	modelID, err := client.GenVideoMotion(modelPath, videoPath, "", string(visvise.OutputModelFormatFBX), "test_auto_alg_video_motion", "", "", boolPtr(false), boolPtr(false), nil)
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

	modelIDs, err := client.GenTextMotion(modelPath, "一个人在原地踏步", "", "", "test_auto_alg_text_motion", "")
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

	modelIDs, err := client.GenPose(modelPath, []visvise.FileInput{mainViewPath}, "", "", "test_auto_alg_pose", "", []string{"main_view.png"})
	if err != nil {
		t.Fatalf("GenPose failed: %v", err)
	}

	t.Logf("PASS: gen_pose (no algorithm_model) - model_ids=%v", modelIDs)
}
