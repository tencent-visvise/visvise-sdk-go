package tests

import (
	"os"
	"testing"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

var (
	appID     = os.Getenv("VISVISE_APP_ID")
	secretKey = os.Getenv("VISVISE_SECRET_KEY")
	rtx       = os.Getenv("VISVISE_RTX")
	assetsDir = "assets"
)

// TestAtomicAPI_GetUserQuota tests the get_user_quota API
func TestAtomicAPI_GetUserQuota(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	client := visvise.NewClient(appID, secretKey,
		visvise.NewClientOptions().SetEnv(visvise.EnvDev).SetDebug(true))
	api := client.GetAPI()

	quota, err := api.GetUserQuota(rtx)
	if err != nil {
		t.Fatalf("GetUserQuota failed: %v", err)
	}

	if quota.Quota < 0 {
		t.Errorf("Expected quota >= 0, got %d", quota.Quota)
	}

	t.Logf("PASS: get_user_quota - quota=%d server_ts=%d", quota.Quota, quota.ServerTS)
}

// TestAtomicAPI_ListAlgorithmModel tests the list_algorithm_model API
func TestAtomicAPI_ListAlgorithmModel(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	client := visvise.NewClient(appID, secretKey, nil)
	api := client.GetAPI()

	testCases := []struct {
		nodeType int
		subType  *int
		label    string
	}{
		{7, nil, "图生360"},
		{3, nil, "图生高模"},
		{4, intPtr(1), "视频生动画"},
		{4, intPtr(2), "文生动画"},
		{5, nil, "骨骼架设"},
		{2, nil, "LOD"},
	}

	for _, tc := range testCases {
		models, err := api.ListAlgorithmModel(tc.nodeType, tc.subType, rtx)
		if err != nil {
			t.Errorf("ListAlgorithmModel %s failed: %v", tc.label, err)
			continue
		}
		if len(models) == 0 {
			t.Errorf("ListAlgorithmModel %s returned empty list", tc.label)
			continue
		}
		t.Logf("PASS: list_algorithm_model %s - first=%s count=%d", tc.label, models[0], len(models))
	}
}

// TestAtomicAPI_GetText2MotionPromptList tests the get_text2motion_prompt_list API
func TestAtomicAPI_GetText2MotionPromptList(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	client := visvise.NewClient(appID, secretKey, nil)
	api := client.GetAPI()

	for _, lang := range []string{"zh", "en"} {
		prompts, err := api.GetText2MotionPromptList(lang, rtx)
		if err != nil {
			t.Errorf("GetText2MotionPromptList lang=%s failed: %v", lang, err)
			continue
		}
		if len(prompts) == 0 {
			t.Errorf("GetText2MotionPromptList lang=%s returned empty list", lang)
			continue
		}
		t.Logf("PASS: get_text2motion_prompt_list lang=%s - count=%d first=%s", lang, len(prompts), prompts[0])
	}
}

// TestAtomicAPI_RemoveBackground tests the remove_bg API
func TestAtomicAPI_RemoveBackground(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(mainViewPath); os.IsNotExist(err) {
		t.Skip("Skipping test: main_view.png not found in tests/assets")
	}

	client := visvise.NewClient(appID, secretKey, nil)
	api := client.GetAPI()

	cosURL, err := client.UploadFile(mainViewPath, "", false, rtx)
	if err != nil {
		t.Fatalf("Upload file failed: %v", err)
	}

	resultURL, err := api.RemoveBackground(cosURL, rtx)
	if err != nil {
		t.Fatalf("RemoveBackground failed: %v", err)
	}

	if len(resultURL) < 4 || resultURL[:4] != "http" {
		t.Errorf("URL format abnormal: %s", resultURL)
	}

	t.Logf("PASS: remove_bg - output_url=%s...", resultURL[:80])
}

// TestAtomicAPI_DownloadModel tests the download_model API
func TestAtomicAPI_DownloadModel(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	client := visvise.NewClient(appID, secretKey, nil)
	api := client.GetAPI()

	knownModelID := "Model2026042300226056"
	url, err := api.DownloadModel(knownModelID, rtx)
	if err != nil {
		t.Fatalf("DownloadModel failed: %v", err)
	}

	if len(url) < 4 || url[:4] != "http" {
		t.Errorf("URL format abnormal: %s", url)
	}

	t.Logf("PASS: download_model - url=%s...", url[:80])
}

// TestAtomicAPI_DeleteModel tests the delete_model API
func TestAtomicAPI_DeleteModel(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	animModelPath := assetsDir + "/animation_model.fbx"
	poseRefPath := assetsDir + "/pose_ref.png"
	mainViewPath := assetsDir + "/main_view.png"
	if _, err := os.Stat(animModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: animation_model.fbx not found")
	}

	client := visvise.NewClient(appID, secretKey, nil)
	api := client.GetAPI()

	animModelURL, err := client.UploadFile(animModelPath, "", false, rtx)
	if err != nil {
		t.Fatalf("Upload animation model failed: %v", err)
	}
	poseRefURL, err := client.UploadFile(poseRefPath, "", false, rtx)
	if err != nil {
		t.Fatalf("Upload pose ref failed: %v", err)
	}
	mainViewURL, err := client.UploadFile(mainViewPath, "", false, rtx)
	if err != nil {
		t.Fatalf("Upload main view failed: %v", err)
	}

	params := map[string]interface{}{
		"algorithm_model":     "VISVISE-PosingAI-V1.0.0",
		"output_model_format": "fbx",
	}
	ids, err := api.BatchGenPose("delete_test", animModelURL, []string{poseRefURL, mainViewURL}, params, rtx)
	if err != nil {
		t.Fatalf("BatchGenPose failed: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 ids, got %d: %v", len(ids), ids)
	}
	t.Logf("Created test model ids: %v", ids)

	err = api.DeleteModel(ids[0], rtx)
	if err != nil {
		t.Fatalf("DeleteModel failed: %v", err)
	}
	t.Logf("PASS: delete_model - deleted %s", ids[0])

	err = api.BatchDeleteModel([]string{ids[1]}, rtx)
	if err != nil {
		t.Fatalf("BatchDeleteModel failed: %v", err)
	}
	t.Logf("PASS: batch_delete_model - batch deleted %s", ids[1])

	models, _, err := api.GetModelList(ids, nil, nil, "", 10, 1, rtx)
	if err != nil {
		t.Fatalf("GetModelList failed: %v", err)
	}
	remaining := make(map[string]bool)
	for _, m := range models {
		remaining[m.ModelID] = true
	}
	for _, mid := range ids {
		if remaining[mid] {
			t.Errorf("Model %s still queryable after deletion", mid)
		}
	}
	t.Logf("PASS: delete verification - all models deleted")
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
