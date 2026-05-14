package tests

import (
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// TestBatch2_MidModelFaceType1 tests gen_mid_model with face_type=1 and fbx format
func TestBatch2_MidModelFaceType1(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMidModelOptions().
		SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle)

	modelID, err := client.GenMidModel(mv["main"], mv["back"], mv["left"], mv["right"], opts)
	if err != nil {
		t.Fatalf("GenMidModel face_type=1 fbx failed: %v", err)
	}

	t.Logf("PASS: mid face_type=1 fbx - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_MidModelFaceType2 tests gen_mid_model with face_type=2 and fbx format
func TestBatch2_MidModelFaceType2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMidModelOptions().
		SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad)

	modelID, err := client.GenMidModel(mv["main"], mv["back"], mv["left"], mv["right"], opts)
	if err != nil {
		t.Fatalf("GenMidModel face_type=2 fbx failed: %v", err)
	}

	t.Logf("PASS: mid face_type=2 fbx - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_RetopologyDetailLevel2Face2 tests gen_retopology with detail_level=2 and face_type=2
func TestBatch2_RetopologyDetailLevel2Face2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenRetopologyOptions().
		SetAlgorithmModel("hunyuan3D-RTP-v1.5").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad).
		SetDetailLevel(visvise.DetailLevelMedium)

	modelID, err := client.GenRetopology(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRetopology detail=2 face=2 failed: %v", err)
	}

	t.Logf("PASS: rtp detail=2 face=2 - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_RetopologyDetailLevel3Face1 tests gen_retopology with detail_level=3 and face_type=1
func TestBatch2_RetopologyDetailLevel3Face1(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenRetopologyOptions().
		SetAlgorithmModel("hunyuan3D-RTP-v1.5").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeTriangle).
		SetDetailLevel(visvise.DetailLevelHigh)

	modelID, err := client.GenRetopology(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRetopology detail=3 face=1 failed: %v", err)
	}

	t.Logf("PASS: rtp detail=3 face=1 - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_MeshRefinePreserveTrue tests gen_mesh_refine with mode=optimize
func TestBatch2_MeshRefinePreserveTrue(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMeshRefineOptions().
		SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").
		SetMode(visvise.MeshRefineModeOptimize)

	modelID, err := client.GenMeshRefine(modelPath, opts)
	if err != nil {
		t.Fatalf("GenMeshRefine mode=optimize failed: %v", err)
	}

	t.Logf("PASS: mr mode=optimize - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_MeshRefinePreserveFalse tests gen_mesh_refine with mode=densify
func TestBatch2_MeshRefinePreserveFalse(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMeshRefineOptions().
		SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").
		SetMode(visvise.MeshRefineModeDensify)

	modelID, err := client.GenMeshRefine(modelPath, opts)
	if err != nil {
		t.Fatalf("GenMeshRefine mode=densify failed: %v", err)
	}

	t.Logf("PASS: mr mode=densify - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_UVSmoothTrue tests gen_uv with enable_auto_smoothing=True
func TestBatch2_UVSmoothTrue(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenUVOptions().
		SetAlgorithmModel("hunyuan3D-UV-v2.0").
		SetEnableAutoSmoothing(true)

	modelID, err := client.GenUV(modelPath, opts)
	if err != nil {
		t.Fatalf("GenUV smooth=True failed: %v", err)
	}

	t.Logf("PASS: uv smooth=True - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_UVSmoothFalse tests gen_uv with enable_auto_smoothing=False
func TestBatch2_UVSmoothFalse(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenUVOptions().
		SetAlgorithmModel("hunyuan3D-UV-v2.0").
		SetEnableAutoSmoothing(false)

	modelID, err := client.GenUV(modelPath, opts)
	if err != nil {
		t.Fatalf("GenUV smooth=False failed: %v", err)
	}

	t.Logf("PASS: uv smooth=False - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_TextureRes1024 tests gen_texture with resolution=1024
func TestBatch2_TextureRes1024(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	refFrontPath := assetsDir + "/tex_ref_front.jpg"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	view := &visvise.View{MainView: refFrontPath}
	opts := visvise.NewGenTextureOptions().
		SetAlgorithmModel("hunyuan3D-TEX-v2.0").
		SetInputView(view).
		SetResolution(1024).
		SetUnwarpUV(false)

	modelID, err := client.GenTexture(modelPath, opts)
	if err != nil {
		t.Fatalf("GenTexture res=1024 failed: %v", err)
	}

	t.Logf("PASS: tex res=1024 - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}

// TestBatch2_TextureRes2048 tests gen_texture with resolution=2048
func TestBatch2_TextureRes2048(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	refFrontPath := assetsDir + "/tex_ref_front.jpg"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	view := &visvise.View{MainView: refFrontPath}
	opts := visvise.NewGenTextureOptions().
		SetAlgorithmModel("hunyuan3D-TEX-v2.0").
		SetInputView(view).
		SetResolution(2048).
		SetUnwarpUV(true)

	modelID, err := client.GenTexture(modelPath, opts)
	if err != nil {
		t.Fatalf("GenTexture res=2048 uv=True failed: %v", err)
	}

	t.Logf("PASS: tex res=2048 uv=True - model_id=%s", modelID)

	waitOpts := &visvise.WaitOptions{Interval: 5, Timeout: 900}
	model, err := client.WaitModel(modelID, waitOpts)
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
	} else {
		t.Logf("Model completed: status=%d time_cost=%d", model.Status, model.TimeCost)
	}
}
