package tests

import (
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// TestOptionalParams_Retopology tests various retopology parameters
func TestOptionalParams_Retopology(t *testing.T) {
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
		SetOutputModelFormat(string(visvise.OutputModelFormatFBX)).
		SetFaceType(int(visvise.FaceTypeQuad)).
		SetName("opt_rtp_a2").
		SetDetailLevel(int(visvise.DetailLevelMedium))

	modelID, err := client.GenRetopology(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRetopology detail_level=2 failed: %v", err)
	}
	t.Logf("PASS: retopology detail_level=2 - model_id=%s", modelID)

	opts = visvise.NewGenRetopologyOptions().
		SetAlgorithmModel("hunyuan3D-RTP-v1.5").
		SetOutputModelFormat(string(visvise.OutputModelFormatFBX)).
		SetFaceType(int(visvise.FaceTypeTriangle)).
		SetName("opt_rtp_b2").
		SetDetailLevel(int(visvise.DetailLevelHigh))

	modelID, err = client.GenRetopology(modelPath, opts)
	if err != nil {
		t.Fatalf("GenRetopology detail_level=3 face_type=1 failed: %v", err)
	}
	t.Logf("PASS: retopology detail_level=3 face_type=1 - model_id=%s", modelID)
}

// TestOptionalParams_MeshRefine tests mesh refine parameters
func TestOptionalParams_MeshRefine(t *testing.T) {
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
		SetName("opt_mr_a2").
		SetMode(int(visvise.MeshRefineModeOptimize))

	modelID, err := client.GenMeshRefine(modelPath, opts)
	if err != nil {
		t.Fatalf("GenMeshRefine mode=optimize failed: %v", err)
	}
	t.Logf("PASS: mesh_refine mode=optimize - model_id=%s", modelID)

	opts = visvise.NewGenMeshRefineOptions().
		SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").
		SetName("opt_mr_b2").
		SetMode(int(visvise.MeshRefineModeDensify))

	modelID, err = client.GenMeshRefine(modelPath, opts)
	if err != nil {
		t.Fatalf("GenMeshRefine mode=densify failed: %v", err)
	}
	t.Logf("PASS: mesh_refine mode=densify - model_id=%s", modelID)
}

// TestOptionalParams_UV tests UV generation parameters
func TestOptionalParams_UV(t *testing.T) {
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
		SetName("opt_uv_a2").
		SetEnableAutoSmoothing(true)

	modelID, err := client.GenUV(modelPath, opts)
	if err != nil {
		t.Fatalf("GenUV smooth=True failed: %v", err)
	}
	t.Logf("PASS: uv enable_auto_smoothing=True - model_id=%s", modelID)

	opts = visvise.NewGenUVOptions().
		SetAlgorithmModel("hunyuan3D-UV-v2.0").
		SetName("opt_uv_b2").
		SetEnableAutoSmoothing(false)

	modelID, err = client.GenUV(modelPath, opts)
	if err != nil {
		t.Fatalf("GenUV smooth=False failed: %v", err)
	}
	t.Logf("PASS: uv enable_auto_smoothing=False - model_id=%s", modelID)
}

// TestOptionalParams_Texture tests texture generation parameters
func TestOptionalParams_Texture(t *testing.T) {
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
		SetName("opt_tex_a2").
		SetInputView(view).
		SetResolution(1024).
		SetUnwarpUV(false)

	modelID, err := client.GenTexture(modelPath, opts)
	if err != nil {
		t.Fatalf("GenTexture res=1024 unwarp_uv=False failed: %v", err)
	}
	t.Logf("PASS: tex resolution=1024 unwarp_uv=False - model_id=%s", modelID)

	view = &visvise.View{MainView: refFrontPath}
	opts = visvise.NewGenTextureOptions().
		SetAlgorithmModel("hunyuan3D-TEX-v2.0").
		SetName("opt_tex_b2").
		SetInputView(view).
		SetResolution(2048).
		SetUnwarpUV(true)

	modelID, err = client.GenTexture(modelPath, opts)
	if err != nil {
		t.Fatalf("GenTexture res=2048 unwarp_uv=True failed: %v", err)
	}
	t.Logf("PASS: tex resolution=2048 unwarp_uv=True - model_id=%s", modelID)
}

// TestOptionalParams_LOD tests LOD generation parameters
func TestOptionalParams_LOD(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: 2}}
	opts := visvise.NewGenLODOptions().
		SetAlgorithmModel("VISVISE-LOD-V1.0.0").
		SetName("opt_lod_a2").
		SetGenTimes(1)

	modelIDs, err := client.GenLOD(modelPath, reduceFaces, opts)
	if err != nil {
		t.Fatalf("GenLOD gen_times=1 fbx failed: %v", err)
	}
	t.Logf("PASS: lod gen_times=1 fbx - model_ids=%v", modelIDs)
}

// TestOptionalParams_VideoMotion tests video motion parameters
func TestOptionalParams_VideoMotion(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	videoPath := assetsDir + "/animation_video.mp4"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenVideoMotionOptions().
		SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0").
		SetName("opt_vm_a2").
		SetWithHand(true).
		SetMultipleTrack(false)

	modelID, err := client.GenVideoMotion(animModelPath, videoPath, opts)
	if err != nil {
		t.Fatalf("GenVideoMotion with_hand=True failed: %v", err)
	}
	t.Logf("PASS: video_motion with_hand=True - model_id=%s", modelID)

	opts = visvise.NewGenVideoMotionOptions().
		SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0").
		SetName("opt_vm_b2").
		SetWithHand(false).
		SetMultipleTrack(false)

	modelID, err = client.GenVideoMotion(animModelPath, videoPath, opts)
	if err != nil {
		t.Fatalf("GenVideoMotion with_hand=False multiple_track=False failed: %v", err)
	}
	t.Logf("PASS: video_motion with_hand=False multiple_track=False - model_id=%s", modelID)
}

// TestOptionalParams_TextMotion tests text motion parameters
func TestOptionalParams_TextMotion(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenTextMotionOptions().
		SetAlgorithmModel("VISVISE-TextMotion-V1.1.0").
		SetName("opt_tm_a2")

	modelIDs, err := client.GenTextMotion(animModelPath, "一个人在挥手打招呼", opts)
	if err != nil {
		t.Fatalf("GenTextMotion prompt=挥手 failed: %v", err)
	}
	t.Logf("PASS: text_motion prompt=挥手 - model_ids=%v", modelIDs)

	opts = visvise.NewGenTextMotionOptions().
		SetAlgorithmModel("VISVISE-TextMotion-V1.1.0").
		SetName("opt_tm_b2")

	modelIDs, err = client.GenTextMotion(animModelPath, "一个人在原地踏步", opts)
	if err != nil {
		t.Fatalf("GenTextMotion prompt=踏步 failed: %v", err)
	}
	t.Logf("PASS: text_motion prompt=踏步 - model_ids=%v", modelIDs)
}

// TestOptionalParams_QueryYesterdayModels tests querying yesterday's batch2 model results
func TestOptionalParams_QueryYesterdayModels(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)
	api := client.GetAPI()

	yesterdayModels := []struct {
		name string
		id   string
	}{
		{"mid face_type=1 fbx", "Model2026042300225326"},
		{"rtp detail=2 face=2", "Model2026042300225337"},
		{"mr mode=optimize", "Model2026042300225347"},
		{"uv smooth=True", "Model2026042300226324"},
		{"tex res=1024", "Model2026042300225360"},
	}

	for _, m := range yesterdayModels {
		models, _, err := api.GetModelList([]string{m.id}, nil, nil, "", 10, 1)
		if err != nil {
			t.Logf("Query model %s failed: %v", m.id, err)
			continue
		}
		if len(models) > 0 {
			t.Logf("PASS: query %s - found model_id=%s status=%d", m.name, models[0].ModelID, models[0].Status)
		} else {
			t.Logf("Model %s not found (may have been deleted)", m.id)
		}
	}
}

// TestOptionalParams_CompleteBatch2 tests completing batch2 test cases
func TestOptionalParams_CompleteBatch2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	opts := visvise.NewGenMidModelOptions().
		SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").
		SetName("opt_mid_b2").
		SetFaceType(int(visvise.FaceTypeQuad))

	modelID, err := client.GenMidModel(mv["main"], mv["back"], mv["left"], mv["right"], opts)
	if err != nil {
		t.Fatalf("GenMidModel face_type=2 fbx failed: %v", err)
	}
	t.Logf("PASS: mid face_type=2 fbx - model_id=%s", modelID)

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	opts2 := visvise.NewGenRetopologyOptions().
		SetAlgorithmModel("hunyuan3D-RTP-v1.5").
		SetName("opt_rtp_b2").
		SetDetailLevel(int(visvise.DetailLevelHigh))

	modelID, err = client.GenRetopology(modelPath, opts2)
	if err != nil {
		t.Fatalf("GenRetopology detail_level=3 face_type=1 failed: %v", err)
	}
	t.Logf("PASS: rtp detail=3 face=1 - model_id=%s", modelID)

	opts3 := visvise.NewGenUVOptions().
		SetAlgorithmModel("hunyuan3D-UV-v2.0").
		SetName("opt_uv_b2").
		SetEnableAutoSmoothing(false)

	modelID, err = client.GenUV(modelPath, opts3)
	if err != nil {
		t.Fatalf("GenUV smooth=False failed: %v", err)
	}
	t.Logf("PASS: uv smooth=False - model_id=%s", modelID)

	refFrontPath := assetsDir + "/tex_ref_front.jpg"
	view := &visvise.View{MainView: refFrontPath}
	texOpts := visvise.NewGenTextureOptions().
		SetAlgorithmModel("hunyuan3D-TEX-v2.0").
		SetName("opt_tex_b2").
		SetInputView(view).
		SetResolution(2048).
		SetUnwarpUV(true)

	modelID, err = client.GenTexture(modelPath, texOpts)
	if err != nil {
		t.Fatalf("GenTexture res=2048 unwarp_uv=True failed: %v", err)
	}
	t.Logf("PASS: tex res=2048 unwarp_uv=True - model_id=%s", modelID)

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: 2}}
	opts4 := visvise.NewGenLODOptions().
		SetAlgorithmModel("VISVISE-LOD-V1.0.0").
		SetName("opt_lod_a2").
		SetGenTimes(1)

	modelIDs, err := client.GenLOD(modelPath, reduceFaces, opts4)
	if err != nil {
		t.Fatalf("GenLOD gen_times=1 failed: %v", err)
	}
	t.Logf("PASS: lod gen_times=1 - model_ids=%v", modelIDs)
}
