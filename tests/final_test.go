package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// TestFinal_QueryYesterdayModels tests querying yesterday's batch2 model_id results
func TestFinal_QueryYesterdayModels(t *testing.T) {
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
		{"mr preserve=True", "Model2026042300225347"},
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
			t.Logf("PASS: [%s] model_id=%s status=%d time_cost=%d", m.name, models[0].ModelID, models[0].Status, models[0].TimeCost)
		} else {
			t.Logf("Model %s not found (already deleted)", m.id)
		}
	}
}

// TestFinal_CompleteBatch2 tests completing batch2 remaining test cases
func TestFinal_CompleteBatch2(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	modelID, err := client.GenMidModel(mv["main"], mv["back"], mv["left"], mv["right"],
		"VISVISE-MeshGen-V1.0.0", string(visvise.OutputModelFormatFBX), 2, "opt_mid_b_final", "", "", "", "", "")
	if err != nil {
		t.Fatalf("GenMidModel face_type=2 fbx failed: %v", err)
	}
	fmt.Printf("Submitted [mid face_type=2 fbx] -> %s\n", modelID)

	modelPath := assetsDir + "/tex_model.obj"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: tex_model.obj not found")
	}

	modelID, err = client.GenRetopology(modelPath, "hunyuan3D-RTP-v1.5", string(visvise.OutputModelFormatFBX), 1, "opt_rtp_b_final", "obj", intPtr(3), nil)
	if err != nil {
		t.Fatalf("GenRetopology detail_level=3 face_type=1 failed: %v", err)
	}
	fmt.Printf("Submitted [rtp detail=3 face=1] -> %s\n", modelID)

	modelID, err = client.GenUV(modelPath, "hunyuan3D-UV-v2.0", "opt_uv_b_final", "", boolPtr(false))
	if err != nil {
		t.Fatalf("GenUV smooth=False failed: %v", err)
	}
	fmt.Printf("Submitted [uv smooth=False] -> %s\n", modelID)

	refFrontPath := assetsDir + "/tex_ref_front.jpg"
	view := &visvise.View{MainView: refFrontPath}
	modelID, err = client.GenTexture(modelPath, "hunyuan3D-TEX-v2.0", "opt_tex_b_final", "", view, intPtr(2048), boolPtr(true), "")
	if err != nil {
		t.Fatalf("GenTexture res=2048 unwarp_uv=True failed: %v", err)
	}
	fmt.Printf("Submitted [tex res=2048 uv=True] -> %s\n", modelID)

	reduceFaces := []visvise.ReduceFace{{ReduceLevel: 1, ReducePercent: 50, FaceType: 2}}
	modelIDs, err := client.GenLOD(modelPath, reduceFaces, "VISVISE-LOD-V1.0.0", string(visvise.OutputModelFormatFBX), "opt_lod_a_final", "", 1)
	if err != nil {
		t.Fatalf("GenLOD gen_times=1 failed: %v", err)
	}
	fmt.Printf("Submitted [lod gen_times=1] -> %v\n", modelIDs)
}

// TestFinal_AnimationTests tests batch3 animation tests
func TestFinal_AnimationTests(t *testing.T) {
	if appID == "" || secretKey == "" || uid == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_UID not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	videoPath := assetsDir + "/animation_video.mp4"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, uid, visvise.EnvProd, 30)

	modelID, err := client.GenVideoMotion(animModelPath, videoPath, "VISVISE-FramingAI-Base-V1.5.0", string(visvise.OutputModelFormatFBX), "opt_vm_a_final", "", "", boolPtr(true), boolPtr(false), nil)
	if err != nil {
		t.Fatalf("GenVideoMotion with_hand=True failed: %v", err)
	}
	fmt.Printf("Submitted [vm with_hand=True] -> %s\n", modelID)

	modelID, err = client.GenVideoMotion(animModelPath, videoPath, "VISVISE-FramingAI-Base-V1.5.0", string(visvise.OutputModelFormatFBX), "opt_vm_b_final", "", "", boolPtr(false), boolPtr(false), nil)
	if err != nil {
		t.Fatalf("GenVideoMotion hand=False multi=False failed: %v", err)
	}
	fmt.Printf("Submitted [vm hand=False multi=False] -> %s\n", modelID)

	modelIDs, err := client.GenTextMotion(animModelPath, "一个人在挥手打招呼", "VISVISE-TextMotion-V1.1.0", string(visvise.OutputModelFormatFBX), "opt_tm_a_final", "")
	if err != nil {
		t.Fatalf("GenTextMotion prompt=挥手 failed: %v", err)
	}
	fmt.Printf("Submitted [tm prompt=挥手] -> %v\n", modelIDs)

	modelIDs, err = client.GenTextMotion(animModelPath, "一个人在原地踏步", "VISVISE-TextMotion-V1.1.0", string(visvise.OutputModelFormatGLB), "opt_tm_b_final", "")
	if err != nil {
		t.Fatalf("GenTextMotion prompt=踏步 glb failed: %v", err)
	}
	fmt.Printf("Submitted [tm prompt=踏步 glb] -> %v\n", modelIDs)
}
